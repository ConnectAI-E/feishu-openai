package openai

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

type AudioToTextRequestBody struct {
	File           string `json:"file"`
	Model          string `json:"model"`
	ResponseFormat string `json:"response_format"`
}

type AudioToTextResponseBody struct {
	Text string `json:"text"`
}

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

func (gpt *ChatGPT) AudioToText(audio string) (string, error) {
	requestBody := AudioToTextRequestBody{
		File:           audio,
		Model:          "whisper-1",
		ResponseFormat: "text",
	}
	audioToTextResponseBody := &AudioToTextResponseBody{}
	err := gpt.sendRequestWithBodyType(gpt.ApiUrl+"/v1/audio/transcriptions",
		"POST", formVoiceDataBody, requestBody, audioToTextResponseBody)
	//fmt.Println(audioToTextResponseBody)
	if err != nil {
		//fmt.Println(err)
		return "", err
	}

	return audioToTextResponseBody.Text, nil
}
