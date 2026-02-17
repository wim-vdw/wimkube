package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	client kubernetes.Interface
}

func NewClientFromKubeconfig(kubeconfigPath string, contextName string) (kubernetes.Interface, error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}
	configOverrides.CurrentContext = contextName
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to load kubeconfig from %s: %w", kubeconfigPath, err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create a client: %w", err)
	}

	return client, nil
}

func NewClient(kubeconfigFilename, contextName string) (*Client, error) {
	c := &Client{}
	client, err := NewClientFromKubeconfig(kubeconfigFilename, contextName)
	if err != nil {
		return nil, err
	}
	c.client = client

	return c, nil
}

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
