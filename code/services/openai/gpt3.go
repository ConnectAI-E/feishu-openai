package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
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

type AudioToTextRequestBody struct {
	File           string `json:"file"`
	Model          string `json:"model"`
	ResponseFormat string `json:"response_format"`
}

type AudioToTextResponseBody struct {
	Text string `json:"text"`
}

type requestBodyType int

const (
	jsonBody requestBodyType = iota
	formDataBody
)

func (gpt ChatGPT) doAPIRequestWithRetry(url, method string, bodyType requestBodyType,
	requestBody interface{}, responseBody interface{}, client *http.Client, maxRetries int) error {
	var api *loadbalancer.API
	var requestBodyData []byte
	var err error
	var writer *multipart.Writer

	switch bodyType {
	case jsonBody:
		api = gpt.Lb.GetAPI()
		requestBodyData, err = json.Marshal(requestBody)
		if err != nil {
			return err
		}
	case formDataBody:
		api = gpt.Lb.GetAPI()
		formBody := &bytes.Buffer{}
		writer = multipart.NewWriter(formBody)
		err = audioMultipartForm(requestBody.(AudioToTextRequestBody), writer)
		if err != nil {
			return err
		}
		err = writer.Close()
		if err != nil {
			return err
		}
		requestBodyData = formBody.Bytes()
	default:
		return errors.New("unknown request body type")
	}

	if api == nil {
		return errors.New("no available API")
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(requestBodyData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if bodyType == formDataBody {
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}
	req.Header.Set("Authorization", "Bearer "+api.Key)

	var response *http.Response
	var retry int
	for retry = 0; retry <= maxRetries; retry++ {
		response, err = client.Do(req)
		//fmt.Println("req", req)
		//fmt.Println("response", response, "err", err)
		if err != nil || response.StatusCode < 200 || response.StatusCode >= 300 {
			gpt.Lb.SetAvailability(api.Key, false)
			if retry == maxRetries {
				break
			}
			time.Sleep(time.Duration(retry+1) * time.Second)
		} else {
			break
		}
	}
	if response != nil {
		defer response.Body.Close()
	}

	if response == nil || response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("%s api failed after %d retries", strings.ToUpper(method), retry)
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

func (gpt ChatGPT) sendRequestWithBodyType(link, method string, bodyType requestBodyType,
	requestBody interface{}, responseBody interface{}) error {
	var err error
	client := &http.Client{Timeout: 110 * time.Second}
	if gpt.HttpProxy == "" {
		err = gpt.doAPIRequestWithRetry(link, method, bodyType,
			requestBody, responseBody, client, 3)
	} else {
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
		err = gpt.doAPIRequestWithRetry(link, method, bodyType,
			requestBody, responseBody, proxyClient, 3)
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
	err = gpt.sendRequestWithBodyType(gpt.ApiUrl+"/v1/chat/completions", "POST",
		jsonBody,
		requestBody, gptResponseBody)

	if err == nil && len(gptResponseBody.Choices) > 0 {
		resp = gptResponseBody.Choices[0].Message
	} else {
		resp = Messages{}
		err = errors.New("openai 请求失败")
	}
	return resp, err
}

// audioMultipartForm creates a form with audio file contents and the name of the model to use for
// audio processing.
func audioMultipartForm(request AudioToTextRequestBody, w *multipart.Writer) error {
	f, err := os.Open(request.File)
	if err != nil {
		return fmt.Errorf("opening audio file: %w", err)
	}

	fw, err := w.CreateFormFile("file", f.Name())
	if err != nil {
		return fmt.Errorf("creating form file: %w", err)
	}

	if _, err = io.Copy(fw, f); err != nil {
		return fmt.Errorf("reading from opened audio file: %w", err)
	}

	fw, err = w.CreateFormField("model")
	if err != nil {
		return fmt.Errorf("creating form field: %w", err)
	}

	modelName := bytes.NewReader([]byte(request.Model))
	if _, err = io.Copy(fw, modelName); err != nil {
		return fmt.Errorf("writing model name: %w", err)
	}
	w.Close()

	return nil
}

func (gpt ChatGPT) GenerateImage(prompt string, size string, n int) ([]string, error) {
	requestBody := ImageGenerationRequestBody{
		Prompt:         prompt,
		N:              n,
		Size:           size,
		ResponseFormat: "b64_json",
	}

	imageResponseBody := &ImageGenerationResponseBody{}
	err := gpt.sendRequestWithBodyType(gpt.ApiUrl+"/v1/images/generations",
		"POST", jsonBody, requestBody, imageResponseBody)

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

func (gpt ChatGPT) AudioToText(audio string) (string, error) {
	requestBody := AudioToTextRequestBody{
		File:           audio,
		Model:          "whisper-1",
		ResponseFormat: "text",
	}
	audioToTextResponseBody := &AudioToTextResponseBody{}
	err := gpt.sendRequestWithBodyType(gpt.ApiUrl+"/v1/audio/transcriptions",
		"POST", formDataBody, requestBody, audioToTextResponseBody)
	//fmt.Println(audioToTextResponseBody)
	if err != nil {
		//fmt.Println(err)
		return "", err
	}

	return audioToTextResponseBody.Text, nil
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
