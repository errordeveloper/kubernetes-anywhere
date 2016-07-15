package main

import (
	"log"
	"os"
	"os/user"

	"github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/cmd"
)

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("Unable to find which user we're running as.", err)
	}
	if user.Uid != "0" {
		log.Fatalf("Please run me as root!\n")
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
