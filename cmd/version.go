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
		fmt.Println("Version:      ", rootCmd.Version)
		fmt.Println("Go version:   ", runtime.Version())
		fmt.Println("Git commit:   ", commit)
		fmt.Println("Build time:   ", buildTime)
		fmt.Println("OS/Arch:      ", runtime.GOOS+"/"+runtime.GOARCH)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	// Hide flags only while rendering this command's help; the flags are shared with rootCmd.
	versionCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		for _, name := range []string{"kubeconfig", "request-timeout"} {
			if f := rootCmd.PersistentFlags().Lookup(name); f != nil {
				f.Hidden = true
				defer func() { f.Hidden = false }()
			}
		}
		rootCmd.HelpFunc()(cmd, args)
	})
}
