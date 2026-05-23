package k8s

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// knownGVR maps "Kind" to its GroupVersionResource for common types.
var knownGVR = map[string]schema.GroupVersionResource{
	"Deployment":             {Group: "apps", Version: "v1", Resource: "deployments"},
	"StatefulSet":            {Group: "apps", Version: "v1", Resource: "statefulsets"},
	"DaemonSet":              {Group: "apps", Version: "v1", Resource: "daemonsets"},
	"ReplicaSet":             {Group: "apps", Version: "v1", Resource: "replicasets"},
	"Service":                {Group: "", Version: "v1", Resource: "services"},
	"ConfigMap":              {Group: "", Version: "v1", Resource: "configmaps"},
	"Secret":                 {Group: "", Version: "v1", Resource: "secrets"},
	"ServiceAccount":         {Group: "", Version: "v1", Resource: "serviceaccounts"},
	"Ingress":                {Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"},
	"ClusterRole":            {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"},
	"ClusterRoleBinding":     {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterrolebindings"},
	"Role":                   {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"},
	"RoleBinding":            {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings"},
	"HorizontalPodAutoscaler": {Group: "autoscaling", Version: "v2", Resource: "horizontalpodautoscalers"},
}

// GVRForKind returns the GroupVersionResource for a given Kubernetes Kind.
// Returns an error if the kind is not in the built-in registry.
func GVRForKind(kind string) (schema.GroupVersionResource, error) {
	gvr, ok := knownGVR[kind]
	if !ok {
		return schema.GroupVersionResource{}, fmt.Errorf("unknown kind %q: register it in knownGVR or use --gvr flag", kind)
	}
	return gvr, nil
}
