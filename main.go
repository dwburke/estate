package main

import (
	"strings"

	"github.com/spf13/viper"

	"github.com/dwburke/prefs/cmd"
)

func init() {
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetDefault("prefs.storage.type", "memory")
	viper.SetDefault("prefs.storage.table", "prefs")

	viper.AutomaticEnv()
}

func main() {
	cmd.Execute()
}
