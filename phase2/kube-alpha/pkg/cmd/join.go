package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/disco"
	"github.com/spf13/cobra"
)

func NewCmdJoin(out io.Writer, config *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join",
		Short: "Run this on other servers to join an existing cluster.",
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
				p2pDisco.Join(strings.Split(args[0], ","))
			} else if config.disco == "token" {
				// TODO should be different in 'join' case
				tokenDisco := disco.NewTokenDisco()
				token := tokenDisco.Init()
				fmt.Println("Your token is:", token)
				fmt.Println("Please pass this into `join` on the other nodes.")
			}
		},
	}
	return cmd
}
