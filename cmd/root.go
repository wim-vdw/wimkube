package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wim-vdw/wimkube/internal"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeConfig *internal.KubeConfig

var rootCmd = &cobra.Command{
	Use:   "wimkube",
	Short: "Interactive Kubernetes CLI.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.HasParent() {
			return nil
		}
		if viper.GetString("kubeconfig") == "" {
			if env := os.Getenv("KUBECONFIG"); env != "" {
				viper.Set("kubeconfig", filepath.SplitList(env)[0])
			} else {
				viper.Set("kubeconfig", clientcmd.RecommendedHomeFile)
			}
		}
		k, err := internal.NewKubeConfig(viper.GetString("kubeconfig"))
		if err != nil {
			return err
		}
		kubeConfig = k
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func SetVersion(version string) {
	rootCmd.Version = version
}

var commit string

func SetCommit(commitHash string) {
	commit = commitHash
}

var buildTime string

func SetBuildTime(bt string) {
	buildTime = bt
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("help", "h", false, "Display this help message.")
	rootCmd.Flags().BoolP("version", "v", false, "Display version information.")
	rootCmd.PersistentFlags().StringP("kubeconfig", "", "", "Path to the kubeconfig file to use. If not specified, $KUBECONFIG or ~/.kube/config will be used.")
	rootCmd.PersistentFlags().IntP("request-timeout", "t", 30, "Timeout in seconds for Kubernetes API requests.")
	rootCmd.SetVersionTemplate("wimkube version: {{ .Version }}\n")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.SilenceUsage = true
	_ = viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
	_ = viper.BindPFlag("request-timeout", rootCmd.PersistentFlags().Lookup("request-timeout"))
}
