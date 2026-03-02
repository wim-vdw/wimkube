package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/term"

	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
)

type Client struct {
	client kubernetes.Interface
	config *rest.Config
}

// NewClient creates a new Kubernetes client using the specified kubeconfig file and context name.
// It returns an error if the kubeconfig file cannot be loaded or if the client cannot be created.
func NewClient(kubeconfigFilename, contextName string) (*Client, error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigFilename}
	configOverrides := &clientcmd.ConfigOverrides{}
	configOverrides.CurrentContext = contextName
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to load kubeconfig from %s: %w", kubeconfigFilename, err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create a client: %w", err)
	}

	return &Client{
		client: client,
		config: config,
	}, nil
}

// GetNamespaces retrieves the list of namespaces in the Kubernetes cluster.
// It returns a slice of namespace names and an error if the namespaces cannot be retrieved.
func (c *Client) GetNamespaces() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("request-timeout"))*time.Second)
	defer cancel()

	namespaces, err := c.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get namespaces: %w", err)
	}
	out := make([]string, 0, len(namespaces.Items))
	for _, ns := range namespaces.Items {
		out = append(out, ns.Name)
	}

	return out, nil
}

// GetPods retrieves the list of pods in the specified namespace.
// It returns a slice of pod names and an error if the pods cannot be retrieved.
func (c *Client) GetPods(namespace string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("request-timeout"))*time.Second)
	defer cancel()

	pods, err := c.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get pods: %w", err)
	}
	out := make([]string, 0, len(pods.Items))
	for _, pod := range pods.Items {
		out = append(out, pod.Name)
	}

	return out, nil
}

// GetContainers retrieves the list of container names in the specified pod and namespace.
// It returns a slice of container names and an error if the pod cannot be retrieved.
func (c *Client) GetContainers(namespace, podName string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("request-timeout"))*time.Second)
	defer cancel()

	pod, err := c.client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get pod %s in namespace %s: %w", podName, namespace, err)
	}
	out := make([]string, 0, len(pod.Spec.Containers))
	for _, container := range pod.Spec.Containers {
		out = append(out, container.Name)
	}

	return out, nil
}

// ExecInContainer executes a shell command in the specified container, pod, and namespace.
// It returns an error if the command cannot be executed or if there is an issue with the terminal setup.
func (c *Client) ExecInContainer(namespace, podName, containerName string) error {
	ctx := context.Background()

	req := c.client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   []string{"/bin/sh", "-c", "command -v bash >/dev/null 2>&1 && bash || sh"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(c.config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("unable to create executor: %w", err)
	}

	// Save original terminal state
	oldState, err := setupTerminal(os.Stdin)
	if err != nil {
		return fmt.Errorf("unable to setup terminal: %w", err)
	}
	defer func(f *os.File, state *term.State) {
		err := restoreTerminal(f, state)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to restore terminal: %v\n", err)
		}
	}(os.Stdin, oldState)

	err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		return fmt.Errorf("unable to execute command: %w", err)
	}

	return nil
}

// GetPodLogs retrieves the logs from a specific container in a pod.
// It returns the logs as a string and an error if the logs cannot be retrieved.
func (c *Client) GetPodLogs(namespace, podName, containerName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("request-timeout"))*time.Second)
	defer cancel()

	logOptions := &corev1.PodLogOptions{
		Container: containerName,
	}

	req := c.client.CoreV1().Pods(namespace).GetLogs(podName, logOptions)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to get pod logs for %s in namespace %s: %w", podName, namespace, err)
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("unable to read pod logs: %w", err)
	}

	return buf.String(), nil
}

// setupTerminal puts the terminal into raw mode and returns the original terminal state.
// It returns an error if the terminal cannot be put into raw mode.
func setupTerminal(f *os.File) (*term.State, error) {
	state, err := term.MakeRaw(int(f.Fd()))
	if err != nil {
		return nil, err
	}
	return state, nil
}

// restoreTerminal restores the terminal to its original state.
// It returns an error if the terminal cannot be restored.
func restoreTerminal(f *os.File, state *term.State) error {
	return term.Restore(int(f.Fd()), state)
}
