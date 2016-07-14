package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

func NewCmdInit(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Run this on the first server you deploy onto.",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	return cmd
}
