package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "wimkube",
	Short: "Interactive Kubernetes CLI.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetString("kubeconfig") == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("could not determine home directory: %w", err)
			}
			viper.Set("kubeconfig", filepath.Join(homeDir, ".kube", "config"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func SetVersion(version string) {
	rootCmd.Version = version
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Display this help message.")
	rootCmd.Flags().BoolP("version", "v", false, "Display version info.")
	rootCmd.PersistentFlags().StringP("kubeconfig", "", "", "Path to the kubeconfig file to use. If not specified, the default kubeconfig will be used.")
	rootCmd.PersistentFlags().IntP("request-timeout", "t", 30, "Timeout in seconds for Kubernetes API requests.")
	rootCmd.SetVersionTemplate("wimkube version: {{ .Version }}\n")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.SilenceUsage = true
	_ = viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
	_ = viper.BindPFlag("request-timeout", rootCmd.PersistentFlags().Lookup("request-timeout"))
}
