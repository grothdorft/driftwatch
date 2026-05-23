package k8s

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// Client wraps a dynamic Kubernetes client for fetching live resources.
type Client struct {
	dynamic dynamic.Interface
}

// NewClient creates a Client from the given kubeconfig path.
// If kubeconfigPath is empty, in-cluster config is attempted.
func NewClient(kubeconfigPath string) (*Client, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("building kubeconfig: %w", err)
	}
	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating dynamic client: %w", err)
	}
	return &Client{dynamic: dyn}, nil
}

// GetResource fetches a live resource by GVR, namespace, and name.
// Use namespace="" for cluster-scoped resources.
func (c *Client) GetResource(
	ctx context.Context,
	gvr schema.GroupVersionResource,
	namespace, name string,
) (*unstructured.Unstructured, error) {
	var ri dynamic.ResourceInterface
	if namespace == "" {
		ri = c.dynamic.Resource(gvr)
	} else {
		ri = c.dynamic.Resource(gvr).Namespace(namespace)
	}
	obj, err := ri.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("getting %s/%s: %w", namespace, name, err)
	}
	return obj, nil
}
