package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ueokande/bow/pkg/bow"
)

func init() {
}

type params struct {
	container  string
	namespace  string
	kubeconfig string
	nohosts    bool
}

func main() {
	params := params{
		namespace: "default",
	}

	cmd := &cobra.Command{}
	cmd.Use = "bow [OPTIONS] POD_SELECTOR -- COMMAND ARGS..."
	cmd.Short = "Exec a command on multiple pods from Kubernetes"

	cmd.Flags().StringVarP(&params.container, "container", "c", params.container, "Container name when multiple containers in pod")
	cmd.Flags().StringVarP(&params.namespace, "namespace", "n", params.namespace, "Kubernetes namespace to use. Default to namespace configured in Kubernetes context")
	cmd.Flags().StringVarP(&params.kubeconfig, "kubeconfig", "", params.kubeconfig, " Path to kubeconfig file to use")
	cmd.Flags().BoolVarP(&params.nohosts, "no-hosts", "", params.nohosts, "Do not print hosts")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		narg := len(args)
		if narg < 2 {
			return cmd.Help()
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		config := bow.Config{
			Namespace:  params.namespace,
			KubeConfig: params.kubeconfig,
			Query:      args[0],
			Container:  params.container,
			Command:    args[1:],
			NoHosts:    params.nohosts,
		}

		err := bow.RunBow(ctx, &config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return nil
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
