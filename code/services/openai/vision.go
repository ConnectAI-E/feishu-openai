package openai

import (
	"errors"
	"start-feishubot/logger"
)

type ImageURL struct {
	URL    string `json:"url,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type ContentType struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}
type VisionMessages struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type VisionRequestBody struct {
	Model     string           `json:"model"`
	Messages  []VisionMessages `json:"messages"`
	MaxTokens int              `json:"max_tokens"`
}

func (gpt *ChatGPT) GetVisionInfo(msg []VisionMessages) (
	resp Messages, err error) {
	requestBody := VisionRequestBody{
		Model:     "gpt-4-vision-preview",
		Messages:  msg,
		MaxTokens: gpt.MaxTokens,
	}
	gptResponseBody := &ChatGPTResponseBody{}
	url := gpt.FullUrl("chat/completions")
	logger.Debug("request body ", requestBody)
	if url == "" {
		return resp, errors.New("无法获取openai请求地址")
	}
	//gpt.ChangeMode("gpt-4-vision-preview")
	//fmt.Println("model", gpt.Model)
	err = gpt.sendRequestWithBodyType(url, "POST", jsonBody, requestBody, gptResponseBody)
	if err == nil && len(gptResponseBody.Choices) > 0 {
		resp = gptResponseBody.Choices[0].Message
	} else {
		logger.Errorf("ERROR %v", err)
		resp = Messages{}
		err = errors.New("openai 请求失败")
	}
	return resp, err
}
