package cmd

import (
	"fmt"
	"log"

	"github.com/kabukky/httpscerts"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(generateCertCmd)
}

var generateCertCmd = &cobra.Command{
	Use:   "generatecert",
	Short: "Generate self-signed certificate",
	Run: func(cmd *cobra.Command, args []string) {
		err := httpscerts.Check("cert.pem", "key.pem")
		// If they are not available, generate new ones.
		if err != nil {
			err = httpscerts.Generate("cert.pem", "key.pem", "127.0.0.1:8081")
			if err != nil {
				log.Fatal("Error: Couldn't create https certs.")
			}
		} else {
			fmt.Println("Cert already exists.")
		}
	},
}
