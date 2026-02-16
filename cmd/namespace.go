package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wim-vdw/wimkube/internal"
)

var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Manage namespaces.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return showNamespaceMenu()
	},
}

var namespaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all namespaces.",
	RunE:  execNamespaceList,
}

var namespaceGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get current namespace.",
	RunE:  execNamespaceGet,
}

var namespaceSetCmd = &cobra.Command{
	Use:   "set [namespace]",
	Short: "Set current namespace.",
	Args:  cobra.ExactArgs(1),
	RunE:  execNamespaceSet,
}

func showNamespaceMenu() error {
	var option, namespace string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select an option").
				Options(
					huh.NewOption("Get current namespace", "1"),
					huh.NewOption("List namespaces", "2"),
					huh.NewOption("Set namespace", "3"),
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
		return execNamespaceGet(nil, nil)
	case "2":
		return execNamespaceList(nil, nil)
	case "3":
		kc, err := internal.NewKubeconfig(viper.GetString("kubeconfig"))
		if err != nil {
			return err
		}
		currentContext, err := kc.GetCurrentContext()
		if err != nil {
			return err
		}
		n, err := internal.NewNamespace(viper.GetString("kubeconfig"), currentContext)
		if err != nil {
			return err
		}
		namespaces, err := n.GetNamespaces()
		if err != nil {
			return err
		}
		currentNamespace, _ := kc.GetCurrentNamespace()
		namespace = currentNamespace
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select a namespace").
					Options(huh.NewOptions(namespaces...)...).
					Value(&namespace),
			),
		)
		err = form.Run()
		if err != nil {
			return err
		}
		return execNamespaceSet(nil, []string{namespace})
	}

	return nil
}

func execNamespaceList(cmd *cobra.Command, args []string) error {
	kc, err := internal.NewKubeconfig(viper.GetString("kubeconfig"))
	if err != nil {
		return err
	}
	currentContext, err := kc.GetCurrentContext()
	if err != nil {
		return err
	}
	n, err := internal.NewNamespace(viper.GetString("kubeconfig"), currentContext)
	if err != nil {
		return err
	}
	namespaces, err := n.GetNamespaces()
	if err != nil {
		return err
	}
	for _, ns := range namespaces {
		fmt.Println(ns)
	}

	return nil
}

func execNamespaceGet(cmd *cobra.Command, args []string) error {
	kc, err := internal.NewKubeconfig(viper.GetString("kubeconfig"))
	if err != nil {
		return err
	}
	currentNamespace, err := kc.GetCurrentNamespace()
	if err != nil {
		return err
	}
	fmt.Println(currentNamespace)

	return nil
}

func execNamespaceSet(cmd *cobra.Command, args []string) error {
	namespace := args[0]
	kc, err := internal.NewKubeconfig(viper.GetString("kubeconfig"))
	if err != nil {
		return err
	}
	err = kc.SetNamespace(namespace)
	if err != nil {
		return err
	}
	fmt.Printf("Current namespace set to: %s\n", namespace)

	return nil
}

func init() {
	rootCmd.AddCommand(namespaceCmd)
	namespaceCmd.AddCommand(namespaceListCmd)
	namespaceCmd.AddCommand(namespaceGetCmd)
	namespaceCmd.AddCommand(namespaceSetCmd)
}
