package k8s

import (
	"os"
	"path/filepath"
)

// DefaultKubeconfigPath returns default kubeconfig path ($HOME/.kube/config)
func DefaultKubeconfigPath() string {
	var homedir string
	if h := os.Getenv("HOME"); h != "" {
		homedir = h
	} else {
		homedir = os.Getenv("USERPROFILE") // windows
	}
	return filepath.Join(homedir, ".kube", "config")
}
