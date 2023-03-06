package main

import (
	"fmt"
	"start-feishubot/handlers"
	"start-feishubot/initialization"
	"start-feishubot/services"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"

	sdkginext "github.com/larksuite/oapi-sdk-gin"

	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
)

var (
	cfg = pflag.StringP("config", "c", "./config.yaml", "apiserver config file path.")
)

func main() {
	pflag.Parse()
	config := initialization.LoadConfig(*cfg)
	initialization.LoadLarkClient(*config)

	gpt := &services.ChatGPT{ApiKey: config.OpenaiApiKey, ApiUrl: config.OpenaiApiUrl}
	handlers.InitHanders(*gpt, *config)

	eventHandler := dispatcher.NewEventDispatcher(
		config.FeishuAppVerificationToken, config.FeishuAppEncryptKey).
		OnP2MessageReceiveV1(handlers.Handler)

	cardHandler := larkcard.NewCardActionHandler(
		config.FeishuAppVerificationToken, config.FeishuAppEncryptKey,
		handlers.CardHandler())

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/webhook/event",
		sdkginext.NewEventHandlerFunc(eventHandler))
	r.POST("/webhook/card",
		sdkginext.NewCardActionHandlerFunc(
			cardHandler))

	fmt.Println("http server started",
		"http://localhost:9000/webhook/event")
	r.Run(":9000")

}
