package CMD

import (
	"GoBookstoreAPI/Routes"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// addCmd represents the add command
var rootCmd = &cobra.Command{
	Use:   "bookstore",
	Short: "bookstore REST API written in Go",
	Long:  `bookstore REST API written in Golang`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		fmt.Println(port)
		Routes.StartAPI(port)
	},
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
