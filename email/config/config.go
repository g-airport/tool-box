package config

import (
	"fmt"

	"github.com/spf13/viper"
)

//This is used to config

type VerifierConfig struct {
	SourceIP    []string
	SourceEmail []string
}

var Config = &VerifierConfig{}

var configPath = "$HOME/work/go/src/github.com/g-airport/tool-box/email/config/"

func InitConfig() {
	vCfg := viper.New()
	vCfg.AddConfigPath(configPath)
	err := vCfg.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("read config panic: %v", err))
	}
	Config.SourceIP = vCfg.GetStringSlice("ip_source_list")
	Config.SourceEmail = vCfg.GetStringSlice("email_list")
}
