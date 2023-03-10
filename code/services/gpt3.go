package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"start-feishubot/initialization"
	"start-feishubot/services/loadbalancer"
	"strings"
	"time"
)

const (
	maxTokens   = 2000
	temperature = 0.7
	engine      = "gpt-3.5-turbo"
)

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatGPTChoiceItem    `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}
type ChatGPTChoiceItem struct {
	Message      Messages `json:"message"`
	Index        int      `json:"index"`
	FinishReason string   `json:"finish_reason"`
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string     `json:"model"`
	Messages         []Messages `json:"messages"`
	MaxTokens        int        `json:"max_tokens"`
	Temperature      float32    `json:"temperature"`
	TopP             int        `json:"top_p"`
	FrequencyPenalty int        `json:"frequency_penalty"`
	PresencePenalty  int        `json:"presence_penalty"`
}
type ChatGPT struct {
	Lb        *loadbalancer.LoadBalancer
	ApiKey    []string
	ApiUrl    string
	HttpProxy string
}

type ImageGenerationRequestBody struct {
	Prompt         string `json:"prompt"`
	N              int    `json:"n"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
}

type ImageGenerationResponseBody struct {
	Created int64 `json:"created"`
	Data    []struct {
		Base64Json string `json:"b64_json"`
	} `json:"data"`
}

func (gpt ChatGPT) doRequest(url, method string,
	requestBody interface{}, responseBody interface{},
	client *http.Client) error {
	api := gpt.Lb.GetAPI()
	if api == nil {
		return errors.New("no available API")
	}

	requestData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api.Key)

	response, err := client.Do(req)
	if err != nil {
		gpt.Lb.SetAvailability(api.Key, false)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode/2 != 100 {
		gpt.Lb.SetAvailability(api.Key, false)
		return fmt.Errorf("%s api %s", strings.ToUpper(method), response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, responseBody)
	if err != nil {
		return err
	}

	gpt.Lb.SetAvailability(api.Key, true)
	return nil
}

func (gpt ChatGPT) sendRequest(link, method string,
	requestBody interface{}, responseBody interface{}) error {
	var err error
	client := &http.Client{Timeout: 110 * time.Second}
	if gpt.HttpProxy == "" {
		err = gpt.doRequest(link, method, requestBody, responseBody, client)
	} else {
		//fmt.Println("using proxy: " + gpt.HttpProxy)
		proxyUrl, err := url.Parse(gpt.HttpProxy)
		if err != nil {
			return err
		}

		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		proxyClient := &http.Client{
			Transport: transport,
			Timeout:   110 * time.Second,
		}

		err = gpt.doRequest(link, method, requestBody, responseBody, proxyClient)
	}

	return err
}

func (gpt ChatGPT) Completions(msg []Messages) (resp Messages, err error) {
	requestBody := ChatGPTRequestBody{
		Model:            engine,
		Messages:         msg,
		MaxTokens:        maxTokens,
		Temperature:      temperature,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	gptResponseBody := &ChatGPTResponseBody{}
	err = gpt.sendRequest(gpt.ApiUrl+"/v1/chat/completions", "POST",
		requestBody, gptResponseBody)

	if err == nil && len(gptResponseBody.Choices) > 0 {
		resp = gptResponseBody.Choices[0].Message
	} else {
		resp = Messages{}
		err = errors.New("openai 请求失败")
	}
	return resp, err
}

func (gpt ChatGPT) GenerateImage(prompt string, size string, n int) ([]string, error) {
	requestBody := ImageGenerationRequestBody{
		Prompt:         prompt,
		N:              n,
		Size:           size,
		ResponseFormat: "b64_json",
	}

	imageResponseBody := &ImageGenerationResponseBody{}
	err := gpt.sendRequest(gpt.ApiUrl+"/v1/images/generations",
		"POST", requestBody, imageResponseBody)

	if err != nil {
		return nil, err
	}

	var b64Pool []string
	for _, data := range imageResponseBody.Data {
		b64Pool = append(b64Pool, data.Base64Json)
	}
	return b64Pool, nil
}

func (gpt ChatGPT) GenerateOneImage(prompt string, size string) (string, error) {
	b64s, err := gpt.GenerateImage(prompt, size, 1)
	if err != nil {
		return "", err
	}
	return b64s[0], nil
}

func NewChatGPT(config initialization.Config) *ChatGPT {
	apiKeys := config.OpenaiApiKeys
	apiUrl := config.OpenaiApiUrl
	httpProxy := config.HttpProxy
	lb := loadbalancer.NewLoadBalancer(apiKeys)
	return &ChatGPT{
		Lb:        lb,
		ApiKey:    apiKeys,
		ApiUrl:    apiUrl,
		HttpProxy: httpProxy,
	}
}
