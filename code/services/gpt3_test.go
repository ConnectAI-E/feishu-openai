package services

import (
	"fmt"
	"start-feishubot/initialization"
	"testing"
)

func TestCompletions(t *testing.T) {
	initialization.LoadConfig("../config.yaml")
	msg := []Messages{
		{Role: "system", Content: "你是一个专业的翻译官，负责中英文翻译。"},
		{Role: "user", Content: "翻译这段话: The assistant messages help store prior responses. They can also be written by a developer to help give examples of desired behavior."},
	}
	resp, err := Completions(msg)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(resp.Content, resp.Role)
}
