package cmd

import (
	"github.com/christfirst/uriel/internal/httpd"
	"github.com/christfirst/uriel/internal/repository"
	"github.com/christfirst/uriel/internal/service"
	"github.com/spf13/cobra"
)

// httpdCmd represents the httpd command
var httpdCmd = &cobra.Command{
	Use:   "httpd",
	Short: "http service",
	Long:  `handles all http requests and routing`,
	Run:   runHTTPServer,
}

func init() {
	rootCmd.AddCommand(httpdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// httpdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// httpdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runHTTPServer(cmd *cobra.Command, args []string) {
	dbClient := repository.Connect(&conf)
	defer repository.Close(dbClient)

	repo := repository.NewRepository(dbClient)
	svc := service.NewService(&conf, repo)
	handler := httpd.NewHandler(&conf, svc)

	handler.Run()
	return
}
