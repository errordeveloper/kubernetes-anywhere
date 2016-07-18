package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func NewCmdToolbox(out io.Writer, config *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "toolbox",
		Short: "Give me a shell where 'kubectl' is available",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("docker --host=unix:///var/run/weave/weave.sock run --tty --interactive --volumes-from=kube-toolbox-pki weaveworks/kubernetes-anywhere:toolbox")
		},
	}
	return cmd
}
