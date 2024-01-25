package CMD

import (
	"GoBookstoreAPI/Routes"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start the server with the specified port",
	Long:  `start the server with the specified port`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		if port == "" {
			log.Fatal("port flag not found")
		}
		fmt.Println("Bookstore API server started on port: " + port)
		Routes.StartAPI(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("port", "8000", "Port of the server")
}
