package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wim-vdw/wimkube/internal"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage contexts.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return showContextMenu()
	},
}

var contextListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all contexts.",
	RunE:  execContextList,
}

var contextGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get current context.",
	RunE:  execContextGet,
}

var contextSetCmd = &cobra.Command{
	Use:   "set [context]",
	Short: "Set current context.",
	Args:  cobra.ExactArgs(1),
	RunE:  execContextSet,
}

func showContextMenu() error {
	var option, contextName string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select an option").
				Options(
					huh.NewOption("Get current context", "1"),
					huh.NewOption("List contexts", "2"),
					huh.NewOption("Set context", "3"),
				).
				Value(&option),
		),
	)
	err := form.Run()
	if err != nil {
		return err
	}
	switch option {
	case "1":
		return execContextGet(nil, nil)
	case "2":
		return execContextList(nil, nil)
	case "3":
		kc, err := internal.NewKubeconfig(viper.GetString("kubeconfig"))
		if err != nil {
			return err
		}
		contextNames, err := kc.GetContextNames()
		if err != nil {
			return err
		}
		if len(contextNames) == 0 {
			return fmt.Errorf("no contexts found in kubeconfig")
		}
		currentContext, _ := kc.GetCurrentContext()
		contextName = currentContext
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select a context").
					Options(huh.NewOptions(contextNames...)...).
					Value(&contextName),
			),
		)
		err = form.Run()
		if err != nil {
			return err
		}
		return execContextSet(nil, []string{contextName})
	}

	return nil
}

func execContextList(cmd *cobra.Command, args []string) error {
	kc, err := internal.NewKubeconfig(viper.GetString("kubeconfig"))
	if err != nil {
		return err
	}
	contextNames, err := kc.GetContextNames()
	if err != nil {
		return err
	}
	if len(contextNames) == 0 {
		fmt.Println("No contexts found in kubeconfig.")
		return nil
	}
	for _, contextName := range contextNames {
		fmt.Println(contextName)
	}

	return nil
}

func execContextGet(cmd *cobra.Command, args []string) error {
	kc, err := internal.NewKubeconfig(viper.GetString("kubeconfig"))
	if err != nil {
		return err
	}
	currentContext, err := kc.GetCurrentContext()
	if err != nil {
		return err
	}
	fmt.Println(currentContext)

	return nil
}

func execContextSet(cmd *cobra.Command, args []string) error {
	contextName := args[0]
	kc, err := internal.NewKubeconfig(viper.GetString("kubeconfig"))
	if err != nil {
		return err
	}
	err = kc.SetContext(contextName)
	if err != nil {
		return err
	}
	fmt.Printf("Current context set to: %s\n", contextName)

	return nil
}

func init() {
	rootCmd.AddCommand(contextCmd)
	contextCmd.AddCommand(contextListCmd)
	contextCmd.AddCommand(contextSetCmd)
	contextCmd.AddCommand(contextGetCmd)
}
