package cmd

import (
	"GoBookstoreAPI/routes"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start the server with the specified port",
	Long:  `start the server with the specified port`,

	Run: func(cmd *cobra.Command, args []string) {
		host, exists := os.LookupEnv("BOOKSTORE_LISTEN")
		if !exists || host == "" {
			host = "localhost"
			fmt.Println("BOOKSTORE_LISTEN environment variable not found")
		}

		port, _ := cmd.Flags().GetString("port")
		if port == "" {
			log.Fatal("port flag not found")
		}
		fmt.Println("Bookstore API server started on " + host + ":" + port)
		routes.StartAPI(host, port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("port", "3000", "Port of the server")
}
