package main

import (
	"fmt"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	"start-feishubot/handlers"
	"start-feishubot/initialization"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	sdkginext "github.com/larksuite/oapi-sdk-gin"

	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
)

var (
	cfg = pflag.StringP("config", "c", "./config.yaml", "apiserver config file path.")
)

func main() {
	pflag.Parse()
	initialization.LoadConfig(*cfg)
	initialization.LoadLarkClient()

	eventHandler := dispatcher.NewEventDispatcher(
		viper.GetString("APP_VERIFICATION_TOKEN"),
		viper.GetString("APP_ENCRYPT_KEY")).
		OnP2MessageReceiveV1(handlers.Handler)

	cardHandler := larkcard.NewCardActionHandler(
		viper.GetString("APP_VERIFICATION_TOKEN"),
		viper.GetString("APP_ENCRYPT_KEY"),
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
