package initialization

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	FeishuAppId                string
	FeishuAppSecret            string
	FeishuAppEncryptKey        string
	FeishuAppVerificationToken string
	FeishuBotName              string
	OpenaiApiKey               string
}

func LoadConfig(cfg string) *Config {
	viper.SetConfigFile(cfg)
	viper.ReadInConfig()
	viper.AutomaticEnv()

	return &Config{
		FeishuAppId:                getViperStringValue("APP_ID"),
		FeishuAppSecret:            getViperStringValue("APP_SECRET"),
		FeishuAppEncryptKey:        getViperStringValue("APP_ENCRYPT_KEY"),
		FeishuAppVerificationToken: getViperStringValue("APP_VERIFICATION_TOKEN"),
		FeishuBotName:              getViperStringValue("BOT_NAME"),
		OpenaiApiKey:               getViperStringValue("OPENAI_KEY"),
	}

}

func getViperStringValue(key string) string {
	value := viper.GetString(key)
	if value == "" {
		panic(fmt.Errorf("%s MUST be provided in environment or config.yaml file", key))
	}
	return value
}
