package cmd

import (
	"fmt"
	"github.com/iamolegga/enviper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/your_org/uriel/internal/config"
	"log/slog"
	"os"
)

var (
	cfgFile string
	conf    config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "uriel",
	Short: "Short description",
	Long:  `Long description`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Set Global logger to json format
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	e := enviper.New(viper.New())

	if cfgFile != "" {
		// Use config file from the flag.
		e.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.Getwd()
		if err != nil {
			slog.Error("getting home dir", err)
			os.Exit(1)
		}

		// Search config in home directory with name ".config" (without extension).
		e.AddConfigPath(home)
		e.SetConfigName(".config")
	}

	e.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := e.ReadInConfig(); err == nil {
		slog.Info("config file read", slog.String("config_file", e.ConfigFileUsed()))
	}

	if err := e.Unmarshal(&conf); err == nil {
		slog.Info("config file unmarshalled", slog.String("config_file", viper.ConfigFileUsed()))
	}
}
