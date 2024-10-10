package initialization

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

var larkClient *lark.Client

func LoadLarkClient(config Config) {
	options := []lark.ClientOptionFunc{
		lark.WithLogLevel(larkcore.LogLevelDebug),
	}
	if config.FeishuBaseUrl != "" {
		options = append(options, lark.WithOpenBaseUrl(config.FeishuBaseUrl))
	}

	larkClient = lark.NewClient(config.FeishuAppId, config.FeishuAppSecret, options...)

}

func GetLarkClient() *lark.Client {
	return larkClient
}
