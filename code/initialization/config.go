package initialization

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	FeishuAppId                string
	FeishuAppSecret            string
	FeishuAppEncryptKey        string
	FeishuAppVerificationToken string
	FeishuBotName              string
	OpenaiApiKeys              []string
	HttpPort                   int
	HttpsPort                  int
	UseHttps                   bool
	CertFile                   string
	KeyFile                    string
}

func LoadConfig(cfg string) *Config {
	viper.SetConfigFile(cfg)
	if err := viper.ReadInConfig(); err != nil {
		return nil
	}
	viper.AutomaticEnv()

	config := &Config{
		FeishuAppId:                getViperStringValue("APP_ID", ""),
		FeishuAppSecret:            getViperStringValue("APP_SECRET", ""),
		FeishuAppEncryptKey:        getViperStringValue("APP_ENCRYPT_KEY", ""),
		FeishuAppVerificationToken: getViperStringValue("APP_VERIFICATION_TOKEN", ""),
		FeishuBotName:              getViperStringValue("BOT_NAME", ""),
		OpenaiApiKeys:              getViperStringValueTable("OPENAI_KEY", nil),
		HttpPort:                   getViperIntValue("HTTP_PORT", 9000),
		HttpsPort:                  getViperIntValue("HTTPS_PORT", 9001),
		UseHttps:                   getViperBoolValue("USE_HTTPS", false),
		CertFile:                   getViperStringValue("CERT_FILE", "cert.pem"),
		KeyFile:                    getViperStringValue("KEY_FILE", "key.pem"),
	}

	return config
}

func getViperStringValue(key string, defaultValue string) string {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getViperStringValueTable(key string, defaultValue []string) []string {
	value := viper.GetStringSlice(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func getViperIntValue(key string, defaultValue int) int {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		fmt.Printf("Invalid value for %s, using default value %d\n", key, defaultValue)
		return defaultValue
	}
	return intValue
}

func getViperBoolValue(key string, defaultValue bool) bool {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		fmt.Printf("Invalid value for %s, using default value %v\n", key, defaultValue)
		return defaultValue
	}
	return boolValue
}

func (config *Config) GetCertFile() string {
	if config.CertFile == "" {
		return "cert.pem"
	}
	if _, err := os.Stat(config.CertFile); err != nil {
		fmt.Printf("Certificate file %s does not exist, using default file cert.pem\n", config.CertFile)
		return "cert.pem"
	}
	return config.CertFile
}

func (config *Config) GetKeyFile() string {
	if config.KeyFile == "" {
		return "key.pem"
	}
	if _, err := os.Stat(config.KeyFile); err != nil {
		fmt.Printf("Key file %s does not exist, using default file key.pem\n", config.KeyFile)
		return "key.pem"
	}
	return config.KeyFile
}
