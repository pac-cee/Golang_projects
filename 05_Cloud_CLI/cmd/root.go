package cmd

import (
	"fmt"
	"os"

	"cloud-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "cloud",
	Short: "A CLI tool for managing cloud resources",
	Long: `A command-line interface tool for managing cloud resources across different providers.
Supports AWS, GCP, and Azure with a unified interface for common operations.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cloud-cli.yaml)")
	rootCmd.PersistentFlags().StringP("provider", "p", "", "cloud provider (aws, gcp, azure)")
	rootCmd.PersistentFlags().StringP("region", "r", "", "cloud region")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "output format (table, json, yaml)")
}

func initConfig() {
	var err error
	cfg, err = config.LoadConfig(cfgFile)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
}
