package internal

import (
	"fmt"
	"os"
	"sort"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type KubeConfig struct {
	FilePath string
	Config   *api.Config
}

type KubeConfigManager interface {
	GetCurrentContext() (string, error)
	SetContext(contextName string) error
	GetContextNames() ([]string, error)
	GetCurrentNamespace() (string, error)
	SetNamespace(namespace string) error
}

// NewKubeConfig creates a new KubeConfig instance by loading the kubeconfig file from the specified path.
// It returns an error if the file cannot be loaded or if there are no contexts in the kubeconfig.
func NewKubeConfig(filePath string) (*KubeConfig, error) {
	k := &KubeConfig{}
	if err := k.init(filePath); err != nil {
		return nil, err
	}

	return k, nil
}

// init loads the kubeconfig file from the specified path and initializes the KubeConfig struct.
// It checks if the file is accessible and if it contains any contexts.
// If the file cannot be loaded or if there are no contexts, it returns an error.
func (k *KubeConfig) init(filePath string) error {
	k.FilePath = filePath
	if _, err := os.Stat(k.FilePath); err != nil {
		return fmt.Errorf("kubeconfig file not accessible: %w", err)
	}
	config, err := clientcmd.LoadFromFile(k.FilePath)
	if err != nil {
		return fmt.Errorf("could not load kubeconfig from %s: %w", k.FilePath, err)
	}
	if len(config.Contexts) == 0 {
		return fmt.Errorf("no contexts found in kubeconfig: %s", k.FilePath)
	}
	k.Config = config

	return nil
}

// GetCurrentContext returns the name of the current context set in the kubeconfig file.
// It returns an error if there is no current context set.
func (k *KubeConfig) GetCurrentContext() (string, error) {
	if k.Config.CurrentContext == "" {
		return "", fmt.Errorf("no current context set in kubeconfig")
	}

	return k.Config.CurrentContext, nil
}

// SetContext sets the current context in the kubeconfig file.
// It returns an error if the specified context does not exist or if there is an issue writing the updated kubeconfig back to the file.
func (k *KubeConfig) SetContext(contextName string) error {
	if _, exists := k.Config.Contexts[contextName]; !exists {
		return fmt.Errorf("context '%s' does not exist", contextName)
	}
	if k.Config.CurrentContext == contextName {
		return nil
	}
	k.Config.CurrentContext = contextName
	if err := clientcmd.WriteToFile(*k.Config, k.FilePath); err != nil {
		return fmt.Errorf("could not write kubeconfig: %w", err)
	}

	return nil
}

// GetContextNames returns a sorted list of all context names available in the kubeconfig file.
func (k *KubeConfig) GetContextNames() ([]string, error) {
	contextNames := make([]string, 0, len(k.Config.Contexts))
	for context := range k.Config.Contexts {
		contextNames = append(contextNames, context)
	}
	sort.Strings(contextNames)

	return contextNames, nil
}

// GetCurrentNamespace returns the namespace for the current context in the kubeconfig file.
// If the current context does not have a namespace set, it returns "default".
// It returns an error if there is no current context or if the current context does not exist.
func (k *KubeConfig) GetCurrentNamespace() (string, error) {
	if k.Config.CurrentContext == "" {
		return "", fmt.Errorf("no current context set in kubeconfig")
	}
	ctx, exists := k.Config.Contexts[k.Config.CurrentContext]
	if !exists {
		return "", fmt.Errorf("current context '%s' does not exist", k.Config.CurrentContext)
	}
	if ctx.Namespace == "" {
		return "default", nil
	}

	return ctx.Namespace, nil
}

// SetNamespace sets the namespace for the current context in the kubeconfig file.
// It returns an error if there is no current context or if the current context does not exist.
func (k *KubeConfig) SetNamespace(namespace string) error {
	if k.Config.CurrentContext == "" {
		return fmt.Errorf("no current context set in kubeconfig")
	}
	context, exists := k.Config.Contexts[k.Config.CurrentContext]
	if !exists {
		return fmt.Errorf("current context '%s' does not exist", k.Config.CurrentContext)
	}
	if context.Namespace == namespace {
		return nil
	}
	context.Namespace = namespace
	if err := clientcmd.WriteToFile(*k.Config, k.FilePath); err != nil {
		return fmt.Errorf("could not write kubeconfig: %w", err)
	}

	return nil
}
