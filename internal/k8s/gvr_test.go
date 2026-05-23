package k8s

import (
	"testing"
)

func TestGVRForKind_Known(t *testing.T) {
	cases := []struct {
		kind     string
		resource string
		group    string
	}{
		{"Deployment", "deployments", "apps"},
		{"Service", "services", ""},
		{"ConfigMap", "configmaps", ""},
		{"Ingress", "ingresses", "networking.k8s.io"},
		{"ClusterRole", "clusterroles", "rbac.authorization.k8s.io"},
	}

	for _, tc := range cases {
		t.Run(tc.kind, func(t *testing.T) {
			gvr, err := GVRForKind(tc.kind)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gvr.Resource != tc.resource {
				t.Errorf("resource: want %q, got %q", tc.resource, gvr.Resource)
			}
			if gvr.Group != tc.group {
				t.Errorf("group: want %q, got %q", tc.group, gvr.Group)
			}
		})
	}
}

func TestGVRForKind_Unknown(t *testing.T) {
	_, err := GVRForKind("Frobnitz")
	if err == nil {
		t.Fatal("expected error for unknown kind, got nil")
	}
}
