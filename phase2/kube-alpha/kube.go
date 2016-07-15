package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/cmd"
)

func main() {
	user, err := exec.Command("whoami").Output()
	if err != nil {
		log.Fatalf("Unable to find which user we're running as: %s", err)
	}
	if string(user) != "root\n" {
		log.Fatalf("Please run me as root!\n")
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
