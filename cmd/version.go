package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display detailed version information.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Version:   ", rootCmd.Version)
		fmt.Println("Go version:", runtime.Version())
		fmt.Println("Git commit:", commit)
		fmt.Println("Build time:", buildTime)
		fmt.Println("OS/Arch:   ", runtime.GOOS+"/"+runtime.GOARCH)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	_ = versionCmd.InheritedFlags().MarkHidden("kubeconfig")
	_ = versionCmd.InheritedFlags().MarkHidden("request-timeout")
}
