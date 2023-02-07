package initialization

import (
	"fmt"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigFile("./config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
