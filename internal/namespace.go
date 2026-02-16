package internal

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Namespace struct {
	client kubernetes.Interface
}

func NewNamespace(kubeconfigFilename, contextName string) (*Namespace, error) {
	n := &Namespace{}
	client, err := NewClientFromKubeconfig(kubeconfigFilename, contextName)
	if err != nil {
		return nil, err
	}
	n.client = client

	return n, nil
}

func (n *Namespace) GetCurrenNamespace() (string, error) {
	return "default", nil
}

func (n *Namespace) GetNamespaces() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()

	namespaces, err := n.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get namespaces: %w", err)
	}
	out := make([]string, 0, len(namespaces.Items))
	for _, ns := range namespaces.Items {
		out = append(out, ns.Name)
	}

	return out, nil
}
