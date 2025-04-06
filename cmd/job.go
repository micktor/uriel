package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// jobCmd represents the migrate command
var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "run a job",
	Long:  `run a job`,
	Run:   runJob,
}

func init() {
	rootCmd.AddCommand(jobCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jobCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jobCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runJob(cmd *cobra.Command, args []string) {
	fmt.Println("Error: must also specify name of the job")
}
