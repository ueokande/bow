package bow

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/ueokande/bow/pkg/k8s"
	"github.com/ueokande/bow/pkg/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// RunBow runs bow application with the config
func RunBow(ctx context.Context, config *Config) error {
	restconfig, err := k8s.NewRestConfig(config.KubeConfig)
	if err != nil {
		return err
	}
	clientset, err := k8s.NewClientset(restconfig)
	if err != nil {
		return err
	}

	rest, err := restclient.RESTClientFor(restconfig)
	if err != nil {
		return err
	}

	resp, err := clientset.CoreV1().Pods(config.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	pods := filterPods(resp.Items, config.Query)

	loggerFactory := log.NewLoggerFactory(config.NoHosts)

	var wg sync.WaitGroup
	for _, pod := range pods {
		pod := pod
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			continue
		}
		c := pod.Spec.Containers[0].Name
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
		}, k8s.ParameterCodec)

		exec, err := remotecommand.NewSPDYExecutor(restconfig, "POST", req.URL())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create SPDY executor: %v\n", err)
			continue
		}

		logger := loggerFactory.NewLogger(pod.Name)
		wg.Add(1)
		r, w := io.Pipe()
		go func() {
			err = exec.Stream(remotecommand.StreamOptions{
				Stdin:  nil,
				Stdout: w,
				Stderr: w,
				Tty:    false,
			})
			if err != nil {
				w.CloseWithError(err)
			}
			w.Close()
		}()
		go func() {
			defer wg.Done()
			s := bufio.NewScanner(r)
			for s.Scan() {
				logger.Println(s.Text())
			}
			if err := s.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Scan exit with", err)
			}
		}()
	}
	wg.Wait()
	return nil
}

func filterPods(pods []corev1.Pod, query string) []corev1.Pod {
	match := func(p corev1.Pod) bool {
		if strings.Contains(p.Name, query) {
			return true
		}
		for _, v := range p.Labels {
			if v == query {
				return true
			}
		}
		return false
	}

	var filtered []corev1.Pod
	for _, p := range pods {
		if match(p) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}
