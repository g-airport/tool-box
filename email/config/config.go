package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type VerifierConfig struct {
	IP          []string
	SourceEmail []string
}

var Config = &VerifierConfig{}

var configPath = "/Users/tqll/work/go/src/github.com/g-airport/tool-box/email/config/config.json"

func InitConfig() {
	vConfig := viper.New()
	vConfig.SetConfigFile(configPath)
	err := vConfig.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("read config panic: %v", err))
	}
	Config.IP = vConfig.GetStringSlice("ip_source_list")
	Config.SourceEmail = vConfig.GetStringSlice("email_list")
}
