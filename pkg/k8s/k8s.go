package k8s

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewRestConfig returns new rest.Config from kubeconfig path
func NewRestConfig(kubeconfig string) (*rest.Config, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	if len(kubeconfig) > 0 {
		rules.Precedence = []string{kubeconfig}
	}
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{
		ClusterDefaults: clientcmd.ClusterDefaults,
	}
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		rules,
		overrides,
	)

	cconfig, err := config.ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get client config")
	}

	// See https://github.com/kubernetes/kubernetes/blob/c6eb9a8ed51f5c63cb351e2a4c13494bf5c303a2/pkg/kubectl/cmd/util/kubectl_match_version.go
	cconfig.GroupVersion = &schema.GroupVersion{Group: "", Version: "v1"}
	if cconfig.APIPath == "" {
		cconfig.APIPath = "/api"
	}
	if cconfig.NegotiatedSerializer == nil {
		cconfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	}
	return cconfig, nil
}

// NewClientset returns a new kubernetes.Clientset from the reset.Config
func NewClientset(config *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create clientset")
	}

	return clientset, nil
}
