package openai

import (
	"errors"
)

const (
	defaultMaxTokens = 2000
	temperature      = 0.7
)

type Model string

const (
	Gpt4       Model = "gpt-4"
	Gpt432k    Model = "gpt-4-32k"
	Gpt35Turbo Model = "gpt-3.5-turbo"
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
	Model            Model      `json:"model"`
	Messages         []Messages `json:"messages"`
	MaxTokens        int        `json:"max_tokens"`
	Temperature      float32    `json:"temperature"`
	TopP             int        `json:"top_p"`
	FrequencyPenalty int        `json:"frequency_penalty"`
	PresencePenalty  int        `json:"presence_penalty"`
}

func NewChatGPTRequestBody(msg []Messages) *ChatGPTRequestBody {
	return &ChatGPTRequestBody{
		Model:            Gpt35Turbo,
		Messages:         msg,
		MaxTokens:        defaultMaxTokens,
		Temperature:      temperature,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}

}

func (gpt ChatGPT) Completions(requestBody *ChatGPTRequestBody) (resp Messages, err error) {
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
