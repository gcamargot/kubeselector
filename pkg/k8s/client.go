package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Client wraps the Kubernetes clientset
type Client struct {
	clientset *kubernetes.Clientset
	namespace string
}

// NewClient creates a new Kubernetes client using the default kubeconfig
func NewClient(namespace string) (*Client, error) {
	// Use the current context in kubeconfig
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	// If no namespace specified, get current namespace from context
	if namespace == "" {
		namespace, _, err = kubeConfig.Namespace()
		if err != nil || namespace == "" {
			namespace = "default"
		}
	}

	return &Client{
		clientset: clientset,
		namespace: namespace,
	}, nil
}

// GetPods retrieves all pods in the namespace
func (c *Client) GetPods(ctx context.Context) ([]corev1.Pod, error) {
	pods, err := c.clientset.CoreV1().Pods(c.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}
	return pods.Items, nil
}

// GetSecrets retrieves all secrets in the namespace
func (c *Client) GetSecrets(ctx context.Context) ([]corev1.Secret, error) {
	secrets, err := c.clientset.CoreV1().Secrets(c.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}
	return secrets.Items, nil
}

// GetConfigMaps retrieves all configmaps in the namespace
func (c *Client) GetConfigMaps(ctx context.Context) ([]corev1.ConfigMap, error) {
	configMaps, err := c.clientset.CoreV1().ConfigMaps(c.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list configmaps: %w", err)
	}
	return configMaps.Items, nil
}

// DeletePod deletes a specific pod
func (c *Client) DeletePod(ctx context.Context, name string) error {
	err := c.clientset.CoreV1().Pods(c.namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete pod %s: %w", name, err)
	}
	return nil
}

// DeleteSecret deletes a specific secret
func (c *Client) DeleteSecret(ctx context.Context, name string) error {
	err := c.clientset.CoreV1().Secrets(c.namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete secret %s: %w", name, err)
	}
	return nil
}

// DeleteConfigMap deletes a specific configmap
func (c *Client) DeleteConfigMap(ctx context.Context, name string) error {
	err := c.clientset.CoreV1().ConfigMaps(c.namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete configmap %s: %w", name, err)
	}
	return nil
}

// GetPodLogs retrieves logs for a specific pod
func (c *Client) GetPodLogs(ctx context.Context, podName string, containerName string) (string, error) {
	opts := &corev1.PodLogOptions{}
	if containerName != "" {
		opts.Container = containerName
	}

	req := c.clientset.CoreV1().Pods(c.namespace).GetLogs(podName, opts)
	logs, err := req.DoRaw(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get logs for pod %s: %w", podName, err)
	}
	return string(logs), nil
}

// DescribePod returns detailed information about a pod
func (c *Client) DescribePod(ctx context.Context, name string) (*corev1.Pod, error) {
	pod, err := c.clientset.CoreV1().Pods(c.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %s: %w", name, err)
	}
	return pod, nil
}

// DescribeSecret returns detailed information about a secret
func (c *Client) DescribeSecret(ctx context.Context, name string) (*corev1.Secret, error) {
	secret, err := c.clientset.CoreV1().Secrets(c.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s: %w", name, err)
	}
	return secret, nil
}

// DescribeConfigMap returns detailed information about a configmap
func (c *Client) DescribeConfigMap(ctx context.Context, name string) (*corev1.ConfigMap, error) {
	configMap, err := c.clientset.CoreV1().ConfigMaps(c.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get configmap %s: %w", name, err)
	}
	return configMap, nil
}

// GetNamespace returns the current namespace
func (c *Client) GetNamespace() string {
	return c.namespace
}
