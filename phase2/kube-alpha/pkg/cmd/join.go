package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

func NewCmdJoin(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join",
		Short: "Run this on other servers to join an existing cluster.",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	return cmd
}
