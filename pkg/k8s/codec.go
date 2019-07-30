package k8s

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	util "k8s.io/apimachinery/pkg/util/runtime"
)

// ParameterCodec is a default parameter codec for Kubernetes API
var ParameterCodec = func() runtime.ParameterCodec {
	var scheme = runtime.NewScheme()
	util.Must(corev1.AddToScheme(scheme))
	return runtime.NewParameterCodec(scheme)
}()
