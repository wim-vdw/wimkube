package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wim-vdw/wimkube/internal"
)

var podCmd = &cobra.Command{
	Use:   "pod",
	Short: "Manage pods.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return showPodMenu()
	},
}

var podListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all pods.",
	RunE:  execPodList,
}

var podContainerListCmd = &cobra.Command{
	Use:   "list-containers [pod-name]",
	Short: "List all containers of a pod.",
	Args:  cobra.ExactArgs(1),
	RunE:  execPodContainerList,
}

var podContainerExecCmd = &cobra.Command{
	Use:   "exec [pod-name] [container-name]",
	Short: "Execute an interactive shell in a container of a pod.",
	Args:  cobra.ExactArgs(2),
	RunE:  execPodContainerExec,
}

func showPodMenu() error {
	var option, podName, containerName string
	kc, err := internal.NewKubeconfig(viper.GetString("kubeconfig"))
	if err != nil {
		return err
	}
	currentContext, err := kc.GetCurrentContext()
	if err != nil {
		return err
	}
	c, err := internal.NewClient(viper.GetString("kubeconfig"), currentContext)
	if err != nil {
		return err
	}
	currentNamespace, err := kc.GetCurrentNamespace()
	if err != nil {
		return err
	}
	title := fmt.Sprintf("Select an option (namespace: %s)", currentNamespace)
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(title).
				Options(
					huh.NewOption("List all pods", "1"),
					huh.NewOption("List all containers of a pod", "2"),
					huh.NewOption("Execute an interactive shell in a container of a pod", "3"),
				).
				Value(&option),
		),
	)
	err = form.Run()
	if err != nil {
		return err
	}
	switch option {
	case "1":
		return execPodList(nil, nil)
	case "2":
		pods, err := c.GetPods(currentNamespace)
		if err != nil {
			return err
		}
		if len(pods) == 0 {
			fmt.Printf("No resources found in %s namespace.", currentNamespace)
			return nil
		}
		title := fmt.Sprintf("Select a pod (namespace: %s)", currentNamespace)
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title(title).
					Options(huh.NewOptions(pods...)...).
					Value(&podName),
			),
		)
		err = form.Run()
		if err != nil {
			return err
		}
		return execPodContainerList(nil, []string{podName})
	case "3":
		pods, err := c.GetPods(currentNamespace)
		if err != nil {
			return err
		}
		if len(pods) == 0 {
			fmt.Printf("No resources found in %s namespace.", currentNamespace)
			return nil
		}
		title := fmt.Sprintf("Select a pod (namespace: %s)", currentNamespace)
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title(title).
					Options(huh.NewOptions(pods...)...).
					Value(&podName),
			),
		)
		err = form.Run()
		if err != nil {
			return err
		}
		containers, err := c.GetContainers(currentNamespace, podName)
		if err != nil {
			return err
		}
		title = fmt.Sprintf("Select a container (namespace: %s, pod: %s)", currentNamespace, podName)
		form = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title(title).
					Options(huh.NewOptions(containers...)...).
					Value(&containerName),
			),
		)
		err = form.Run()
		if err != nil {
			return err
		}
		return execPodContainerExec(nil, []string{podName, containerName})
	}

	return nil
}

func execPodList(cmd *cobra.Command, args []string) error {
	kc, err := internal.NewKubeconfig(viper.GetString("kubeconfig"))
	if err != nil {
		return err
	}
	currentContext, err := kc.GetCurrentContext()
	if err != nil {
		return err
	}
	c, err := internal.NewClient(viper.GetString("kubeconfig"), currentContext)
	if err != nil {
		return err
	}
	currentNamespace, err := kc.GetCurrentNamespace()
	if err != nil {
		return err
	}
	pods, err := c.GetPods(currentNamespace)
	if err != nil {
		return err
	}
	if len(pods) == 0 {
		fmt.Printf("No resources found in %s namespace.", currentNamespace)
		return nil
	}
	for _, podName := range pods {
		fmt.Println(podName)
	}

	return nil
}

func execPodContainerList(cmd *cobra.Command, args []string) error {
	podName := args[0]
	kc, err := internal.NewKubeconfig(viper.GetString("kubeconfig"))
	if err != nil {
		return err
	}
	currentContext, err := kc.GetCurrentContext()
	if err != nil {
		return err
	}
	c, err := internal.NewClient(viper.GetString("kubeconfig"), currentContext)
	if err != nil {
		return err
	}
	currentNamespace, err := kc.GetCurrentNamespace()
	if err != nil {
		return err
	}
	containers, err := c.GetContainers(currentNamespace, podName)
	if err != nil {
		return err
	}
	for _, containerName := range containers {
		fmt.Println(containerName)
	}

	return nil
}

func execPodContainerExec(cmd *cobra.Command, args []string) error {
	podName := args[0]
	containerName := args[1]
	kc, err := internal.NewKubeconfig(viper.GetString("kubeconfig"))
	if err != nil {
		return err
	}
	currentContext, err := kc.GetCurrentContext()
	if err != nil {
		return err
	}
	c, err := internal.NewClient(viper.GetString("kubeconfig"), currentContext)
	if err != nil {
		return err
	}
	currentNamespace, err := kc.GetCurrentNamespace()
	if err != nil {
		return err
	}
	err = c.ExecInContainer(currentNamespace, podName, containerName)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(podCmd)
	podCmd.AddCommand(podListCmd)
	podCmd.AddCommand(podContainerListCmd)
	podCmd.AddCommand(podContainerExecCmd)
}
