package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	BASEURL     = "https://api.openai.com/v1/"
	maxTokens   = 2000
	temperature = 0.7
	engine      = "gpt-3.5-turbo"
)

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChoiceItem           `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}
type ChoiceItem struct {
	Message      Messages `json:"message"`
	Index        int      `json:"index"`
	FinishReason string   `json:"finish_reason"`
}

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
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
	ApiKey string
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
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return resp, err
	}
	log.Printf("request gtp json string : %v", string(requestData))
	req, err := http.NewRequest("POST", BASEURL+"chat/completions", bytes.NewBuffer(requestData))
	if err != nil {
		return resp, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+gpt.ApiKey)
	client := &http.Client{Timeout: 110 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	defer response.Body.Close()
	if response.StatusCode/2 != 100 {
		return resp, fmt.Errorf("gtp api %s", response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return resp, err
	}

	gptResponseBody := &ChatGPTResponseBody{}
	// log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		return resp, err
	}

	resp = gptResponseBody.Choices[0].Message
	return resp, nil
}

func FormatQuestion(question string) string {
	return "Answer:" + question
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
		URL string `json:"url"`
	} `json:"data"`
}

func (gpt ChatGPT) GenerateImage(prompt string, size string,
	n int) ([]string, error) {
	requestBody := ImageGenerationRequestBody{
		Prompt:         prompt,
		N:              n,
		Size:           size,
		ResponseFormat: "url",
	}
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", BASEURL+"images/generations", bytes.NewBuffer(requestData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+gpt.ApiKey)
	client := &http.Client{Timeout: 110 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode/2 != 100 {
		return nil, fmt.Errorf("image generation api %s",
			response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("image generation response body: %s", string(body))

	imageResponseBody := &ImageGenerationResponseBody{}
	err = json.Unmarshal(body, imageResponseBody)
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, data := range imageResponseBody.Data {
		urls = append(urls, data.URL)
	}
	return urls, nil

}

func (gpt ChatGPT) GenerateOneImage(prompt string, size string) (string, error) {
	urls, err := gpt.GenerateImage(prompt, size, 1)
	if err != nil {
		return "", err
	}
	return urls[0], nil
}
