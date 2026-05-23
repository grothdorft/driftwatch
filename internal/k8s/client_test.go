package k8s

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/fake"
)

func TestGetResource_Found(t *testing.T) {
	gvr := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	obj := &unstructured.Unstructured{}
	obj.SetName("my-deploy")
	obj.SetNamespace("default")
	obj.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"})

	scheme := runtime.NewScheme()
	fakeDyn := fake.NewSimpleDynamicClient(scheme, obj)

	c := &Client{dynamic: fakeDyn}

	got, err := c.GetResource(context.Background(), gvr, "default", "my-deploy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.GetName() != "my-deploy" {
		t.Errorf("expected name my-deploy, got %s", got.GetName())
	}
}

func TestGetResource_NotFound(t *testing.T) {
	gvr := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	scheme := runtime.NewScheme()
	fakeDyn := fake.NewSimpleDynamicClient(scheme)

	c := &Client{dynamic: fakeDyn}

	_, err := c.GetResource(context.Background(), gvr, "default", "missing")
	if err == nil {
		t.Fatal("expected error for missing resource, got nil")
	}
}
