package services

import (
	"context"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/viper"
	"log"
)

const (
	maxTokens   = 2000
	temperature = 0.7
	engine      = gpt3.TextDavinci003Engine
)

func GetAnswer(question string) (reply string, ok bool) {
	client := gpt3.NewClient(viper.GetString("OPENAI_KEY"))

	ok = false
	reply = ""
	ctx := context.Background()
	resp, err := client.CompletionWithEngine(ctx, engine, gpt3.CompletionRequest{
		Prompt: []string{
			question,
		},
		MaxTokens:   gpt3.IntPtr(maxTokens),
		Temperature: gpt3.Float32Ptr(temperature),
	})
	if err != nil {
		log.Fatalln(err)
	}
	reply = resp.Choices[0].Text
	if reply != "" {
		ok = true
	}
	return reply, ok
}

func FormatQuestion(question string) string {
	return "Answer:" + question
}
