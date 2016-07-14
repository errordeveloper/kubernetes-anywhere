package main

import (
	"os"

	"github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
