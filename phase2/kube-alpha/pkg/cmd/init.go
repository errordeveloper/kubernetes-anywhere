package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/disco"
	"github.com/spf13/cobra"
)

func NewCmdInit(out io.Writer, config *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Run this on the first server you deploy onto.",
		Run: func(cmd *cobra.Command, args []string) {
			// XXX remove duplication
			// Maybe there's a nicer way of doing this
			if config.disco == "weave" || config.disco == "consul" {
				// these are p2p discos
				var p2pDisco disco.P2PDiscovery
				if config.disco == "weave" {
					p2pDisco = disco.NewWeaveDisco()
				} else if config.disco == "consul" {
					p2pDisco = disco.NewConsulDisco()
				}
				p2pDisco.Init(strings.Split(args[0], ","))
			} else if config.disco == "token" {
				tokenDisco := disco.NewTokenDisco()
				token := tokenDisco.Init()
				fmt.Println("Your token is:", token)
				fmt.Println("Please pass this into `join` on the other nodes.")
			}
		},
	}
	return cmd
}
