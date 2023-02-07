package initialization

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/spf13/viper"
)

var larkClient *lark.Client

func LoadLarkClient() {
	larkClient = lark.NewClient(viper.GetString("APP_ID"),
		viper.GetString("APP_SECRET"))
}

func GetLarkClient() *lark.Client {
	return larkClient
}
