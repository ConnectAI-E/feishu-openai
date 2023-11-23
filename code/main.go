package main

import (
	"context"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"start-feishubot/handlers"
	"start-feishubot/initialization"
	"start-feishubot/logger"
	"start-feishubot/utils"

	"github.com/gin-gonic/gin"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/spf13/pflag"
	"start-feishubot/services/openai"
)

func main() {
	initialization.InitRoleList()
	pflag.Parse()

	// 第一次加载以后，不再加载。好处是更快了，坏处是没法动态修改配置
	globalConfig := initialization.GetConfig()
	// 打印一下实际读取到的配置(开发的时候用的，正式环境别用，因为会打印出来 Key)
	//globalConfigPrettyString, _ := json.MarshalIndent(globalConfig, "", "    ")
	//log.Println(string(globalConfigPrettyString))

	initialization.LoadLarkClient(*globalConfig)
	gpt := openai.NewChatGPT(*globalConfig)
	handlers.InitHandlers(gpt, *globalConfig)

	// 是否开启文件日志
	if globalConfig.EnableLog {
		logger2 := enableLog()
		defer utils.CloseLogger(logger2)
	}

	eventHandler := dispatcher.NewEventDispatcher(
		globalConfig.FeishuAppVerificationToken, globalConfig.FeishuAppEncryptKey).
		OnP2MessageReceiveV1(handlers.Handler).
		OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
			logger.Debugf("收到请求 %v", event.RequestURI)
			return handlers.ReadHandler(ctx, event)
		})

	cardHandler := larkcard.NewCardActionHandler(
		globalConfig.FeishuAppVerificationToken, globalConfig.FeishuAppEncryptKey,
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

	if err := initialization.StartServer(*globalConfig, r); err != nil {
		logger.Fatalf("failed to start server: %v", err)
	}
}

func enableLog() *lumberjack.Logger {
	// Set up the logger
	var logger *lumberjack.Logger

	logger = &lumberjack.Logger{
		Filename: "logs/app.log",
		MaxSize:  100,      // megabytes
		MaxAge:   365 * 10, // days
	}

	fmt.Printf("logger %T\n", logger)

	// Set up the logger to write to both file and console
	log.SetOutput(io.MultiWriter(logger, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)

	// Write some log messages
	log.Println("Starting application...")

	return logger
}
