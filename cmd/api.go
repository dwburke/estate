package cmd

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dwburke/lode/api"
)

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.Flags().Int("port", 4444, "Port to run server on")
	viper.BindPFlag("port", apiCmd.Flags().Lookup("port"))

	apiCmd.Flags().String("ssl-key", "key.pem", "SSL certificate key")
	viper.BindPFlag("ssl-key", apiCmd.Flags().Lookup("ssl-key"))

	apiCmd.Flags().String("ssl-cert", "cert.pem", "SSL certificate")
	viper.BindPFlag("ssl-cert", apiCmd.Flags().Lookup("ssl-cert"))

	apiCmd.Flags().Bool("https", false, "Use HTTPS")
	viper.BindPFlag("https", apiCmd.Flags().Lookup("https"))
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the REST api",
	Long:  `Start the REST api`,
	Run: func(cmd *cobra.Command, args []string) {

		r := gin.Default()
		r.Use(cors.Default())

		api.SetupRoutes(r)

		listen := fmt.Sprintf(":%d", viper.GetInt("port"))

		if viper.GetBool("https") == true {
			r.RunTLS(listen, viper.GetString("ssl-cert"), viper.GetString("ssl-key"))
		} else {
			r.Run(listen)
		}

	},
}
