package CMD

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "bookstore",
	Short: "bookstore REST API written in Go",
	Long:  `bookstore REST API written in Golang`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().String("port", "8000", "Port of the server")
}
