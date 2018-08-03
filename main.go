package main

import (
	"strings"

	"github.com/spf13/viper"

	"github.com/dwburke/lode/cmd"
)

func init() {
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetDefault("lode.storage.type", "memory")
	viper.SetDefault("lode.storage.table", "lode")

	viper.AutomaticEnv()
}

func main() {
	cmd.Execute()
}
