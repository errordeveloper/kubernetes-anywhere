package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

func NewCmdToolbox(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "toolbox",
		Short: "Give me a shell where 'kubectl' is available",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	return cmd
}
