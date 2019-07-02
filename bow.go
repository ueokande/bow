package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func RunBow(ctx context.Context, config *Config) error {
	clientset, err := NewClientset(config.KubeConfig)
	if err != nil {
		return err
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	return nil
}
