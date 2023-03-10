package services

import (
	"fmt"
	"start-feishubot/initialization"
	"testing"
)

func TestCompletions(t *testing.T) {
	config := initialization.LoadConfig("../config.yaml")

	msgs := []Messages{
		{Role: "system", Content: "你是一个专业的翻译官，负责中英文翻译。"},
		{Role: "user", Content: "翻译这段话: The assistant messages help store prior responses. They can also be written by a developer to help give examples of desired behavior."},
	}

	gpt := NewChatGPT(*config)

	resp, err := gpt.Completions(msgs)
	if err != nil {
		t.Errorf("TestCompletions failed with error: %v", err)
	}

	fmt.Println(resp.Content, resp.Role)
}

func TestGenerateOneImage(t *testing.T) {
	config := initialization.LoadConfig("../config.yaml")

	gpt := NewChatGPT(*config)
	prompt := "a red apple"
	size := "256x256"

	imageURL, err := gpt.GenerateOneImage(prompt, size)
	if err != nil {
		t.Errorf("TestGenerateOneImage failed with error: %v", err)
	}

	if imageURL == "" {
		t.Errorf("TestGenerateOneImage returned empty imageURL")
	}
}
