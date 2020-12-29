package config

import (
	"strings"

	"github.com/spf13/viper"
)

func NewConfig(fileName, prefix string, cfg interface{}) error {
	v := viper.New()

	v.SetConfigName(fileName)
	v.AddConfigPath("configs")
	v.AddConfigPath(".")
	v.SetConfigType("yaml")
	v.SetEnvPrefix(prefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(&cfg)
	if err != nil {
		return err
	}
	return nil
}
