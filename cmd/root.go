package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wimkube",
	Short: "Interactive Kubernetes CLI.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func SetVersion(version string) {
	rootCmd.Version = version
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Display this help message.")
	rootCmd.Flags().BoolP("version", "v", false, "Display version info.")
	rootCmd.SetVersionTemplate("wimkube version: {{ .Version }}\n")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.SilenceUsage = true
}
