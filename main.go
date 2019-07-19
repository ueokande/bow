package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var homedir string

func init() {
	if h := os.Getenv("HOME"); h != "" {
		homedir = h
	} else {
		homedir = os.Getenv("USERPROFILE") // windows
	}
}

type Params struct {
	container  string
	namespace  string
	kubeconfig string
}

func main() {
	params := Params{
		namespace:  "default",
		kubeconfig: filepath.Join(homedir, ".kube", "config"),
	}

	cmd := &cobra.Command{}
	cmd.Use = "bow POD_SELECTOR"
	cmd.Short = "Run commands on multiple pods and containers from Kubernetes"

	cmd.Flags().StringVarP(&params.container, "container", "c", params.container, "Container name when multiple containers in pod")
	cmd.Flags().StringVarP(&params.namespace, "namespace", "n", params.namespace, "Kubernetes namespace to use. Default to namespace configured in Kubernetes context")
	cmd.Flags().StringVarP(&params.kubeconfig, "kubeconfig", "", params.kubeconfig, " Path to kubeconfig file to use")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		narg := len(args)
		if narg < 2 {
			return cmd.Help()
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		config := Config{
			ContextName: params.container,
			Namespace:   params.namespace,
			KubeConfig:  params.kubeconfig,
			Query:       args[0],
			Command:     args[1:],
		}

		err := RunBow(ctx, &config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return nil
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}