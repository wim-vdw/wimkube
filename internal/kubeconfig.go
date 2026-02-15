package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type Kubeconfig struct {
	KubeconfigPath string
}

func NewKubeconfig() (*Kubeconfig, error) {
	k := &Kubeconfig{}
	if err := k.init(); err != nil {
		return nil, err
	}

	return k, nil
}

func (k *Kubeconfig) init() error {
	// Try to get kubeconfig path from KUBECONFIG environment variable
	kubeconfigPath := os.Getenv("KUBECONFIG")

	// If not set, use default location
	if kubeconfigPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("could not determine home directory: %w", err)
		}
		kubeconfigPath = filepath.Join(homeDir, ".kube", "config")
	}

	k.KubeconfigPath = kubeconfigPath

	return nil
}

func (k *Kubeconfig) LoadContexts() (*api.Config, error) {
	config, err := clientcmd.LoadFromFile(k.KubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("could not load kubeconfig: %w", err)
	}

	return config, nil
}

func (k *Kubeconfig) GetCurrentContext() (string, error) {
	config, err := k.LoadContexts()
	if err != nil {
		return "", err
	}
	if config.CurrentContext == "" {
		return "", fmt.Errorf("no current context set in kubeconfig")
	}

	return config.CurrentContext, nil
}

func (k *Kubeconfig) GetContextNames() ([]string, error) {
	config, err := k.LoadContexts()
	if err != nil {
		return nil, err
	}

	contextNames := make([]string, 0, len(config.Contexts))
	for context := range config.Contexts {
		contextNames = append(contextNames, context)
	}
	sort.Strings(contextNames)

	return contextNames, nil
}

func (k *Kubeconfig) SetContext(contextName string) error {
	config, err := k.LoadContexts()
	if err != nil {
		return err
	}

	if _, exists := config.Contexts[contextName]; !exists {
		return fmt.Errorf("context '%s' does not exist", contextName)
	}

	config.CurrentContext = contextName

	if err := clientcmd.WriteToFile(*config, k.KubeconfigPath); err != nil {
		return fmt.Errorf("could not write kubeconfig: %w", err)
	}

	return nil
}
