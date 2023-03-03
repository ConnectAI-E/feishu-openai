package main

import (
	"fmt"
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

	handler := dispatcher.NewEventDispatcher(viper.GetString(
		"APP_VERIFICATION_TOKEN"), viper.GetString("APP_ENCRYPT_KEY")).
		OnP2MessageReceiveV1(handlers.Handler)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 在已有 Gin 实例上注册消息处理路由
	r.POST("/webhook/event", sdkginext.NewEventHandlerFunc(handler))

	fmt.Println("http server started",
		"http://localhost:9000/webhook/event")

	r.Run(":9000")

}
