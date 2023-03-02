package services

import (
	"fmt"
	"start-feishubot/initialization"
	"testing"
)

func TestCompletions(t *testing.T) {
	initialization.LoadConfig("../config.yaml")
	msg := []Messages{
		{Role: "user", Content: "你好"},
	}
	resp, err := Completions(msg)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(resp.Content, resp.Role)
}
