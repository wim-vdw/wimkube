package internal

import (
	"fmt"
	"sort"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type Kubeconfig struct {
	KubeconfigFilename string
}

func NewKubeconfig(kubeconfigFilename string) (*Kubeconfig, error) {
	k := &Kubeconfig{}
	if err := k.init(kubeconfigFilename); err != nil {
		return nil, err
	}

	return k, nil
}

func (k *Kubeconfig) init(kubeconfigFilename string) error {
	k.KubeconfigFilename = kubeconfigFilename

	return nil
}

func (k *Kubeconfig) LoadContexts() (*api.Config, error) {
	config, err := clientcmd.LoadFromFile(k.KubeconfigFilename)
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

	if err := clientcmd.WriteToFile(*config, k.KubeconfigFilename); err != nil {
		return fmt.Errorf("could not write kubeconfig: %w", err)
	}

	return nil
}
