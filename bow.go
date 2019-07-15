package main

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func RunBow(ctx context.Context, config *Config) error {
	restconfig, err := NewRestConfig(config.KubeConfig)
	if err != nil {
		return err
	}
	clientset, err := NewClientset(restconfig)
	if err != nil {
		return err
	}

	rest, err := restclient.RESTClientFor(restconfig)
	if err != nil {
		return err
	}

	pods, err := clientset.CoreV1().Pods(config.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			fmt.Fprintf(os.Stderr, "Pod %s is not ready: %s\n", pod.Name, pod.Status.Phase)
			continue
		}

		c := pod.Spec.Containers[0].Name
		fmt.Fprintf(os.Stderr, "==== %s/%s ====\n", pod.Name, c)
		req := rest.Post().
			Resource("pods").
			Name(pod.Name).
			Namespace(pod.Namespace).
			SubResource("exec")
		req.VersionedParams(&corev1.PodExecOptions{
			Container: c,
			Command:   config.Command,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, ParameterCodec)

		exec, err := remotecommand.NewSPDYExecutor(restconfig, "POST", req.URL())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create SPDY executor: %v\n", err)
			continue
		}
		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  nil,
			Stdout: os.Stdout,
			Stderr: os.Stdin,
			Tty:    false,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to exec command: %v\n", err)
			continue
		}
	}
	return nil
}
