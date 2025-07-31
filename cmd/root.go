package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bookstore",
	Short: "bookstore REST API written in Go",
	Long:  `bookstore REST API written in Golang`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("bookstore REST API written in Go")
		fmt.Println("To start the server, pass the arguments 'serve --port <port>'")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
