package main

import (
	"strings"

	"github.com/spf13/viper"

	"github.com/dwburke/estate/cmd"
)

func init() {
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetDefault("estate.storage.type", "memory")
	viper.SetDefault("estate.storage.table", "estate")

	viper.AutomaticEnv()
}

func main() {
	cmd.Execute()
}
