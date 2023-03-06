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
	BASEURL     = "/v1/"
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
	ApiUrl string
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
	req, err := http.NewRequest("POST", gpt.ApiUrl+BASEURL+"chat/completions", bytes.NewBuffer(requestData))
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
