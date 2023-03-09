package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"start-feishubot/handlers"
	"start-feishubot/initialization"
	"start-feishubot/services"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"

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

	gpt := services.NewChatGPT(config.OpenaiApiKeys)
	gpt.StartApiKeyAvailabilityCheck()
	handlers.InitHandlers(gpt, *config)

	eventHandler := dispatcher.NewEventDispatcher(
		config.FeishuAppVerificationToken, config.FeishuAppEncryptKey).
		OnP2MessageReceiveV1(handlers.Handler).
		OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
			return handlers.ReadHandler(ctx, event)
		})

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

	if config.UseHttps {
		certFile := config.GetCertFile()
		keyFile := config.GetKeyFile()

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			panic(err)
		}

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", config.HttpsPort),
			Handler: r,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		}

		fmt.Printf("https server started: https://localhost:%d/webhook/event\n", config.HttpsPort)
		err = server.ListenAndServeTLS("", "")
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("http server started: http://localhost:%d/webhook/event\n", config.HttpPort)
		err := r.Run(fmt.Sprintf(":%d", config.HttpPort))
		if err != nil {
			panic(err)
		}
	}
}
