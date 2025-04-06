package cmd

import (
	"github.com/spf13/cobra"
	"github.com/your_org/uriel/internal/job"
)

// singCmd represents the migrate command
var singCmd = &cobra.Command{
	Use:   "sing",
	Short: "sing praise",
	Long:  `shout to the lord`,
	Run:   singJob,
}

func init() {
	jobCmd.AddCommand(singCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// singJob.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// singJob.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func singJob(cmd *cobra.Command, args []string) {
	job.Sing()
}
