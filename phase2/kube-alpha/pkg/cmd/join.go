package cmd

import (
	_ "fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/discovery_providers"
	. "github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/interfaces"
	"github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/launchers"
	"github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/pki"
)

func NewCmdJoin(out io.Writer, config *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join",
		Short: "Run this on other servers to join an existing cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			info := ClusterInfo{MasterIPs: strings.Split(args[0], ",")}
			var (
				_disco    NodeDiscoveryProvider
				_pki      NodePKI
				_launcher NodeLauncher
			)

			_disco = discovery_providers.NewWeaveNodeDiscoveryProvider(info)

			_disco.Setup()
			_disco.Launch()

			_pki = pki.NewContainerizedWeaveNodePKI(info)

			_pki.Fetch()
			_pki.Init()

			_launcher = launchers.NewContainerizedWeaveNodeLauncher(info)

			_launcher.Launch()

			// Alternative version, we could change these interfaces to be more like:
			/*
				info := ClusterInfo{MasterIPs: strings.Split(args[0], ",")}
				var (
					_disco    DiscoveryProvider
					_pki      PKI
					_launcher Launcher
				)

				_disco = discovery_providers.NewWeaveDiscoveryProvider(info)

				_disco.NodeSetup()
				_disco.NodeLaunch()

				_pki = pki.NewContainerizedWeavePKI(info)

				_pki.NodeFetch()
				_pki.NodeInit()

				_launcher = launchers.NewContainerizedWeaveLauncher(info)

				_launcher.NodeLaunch()

				// Or may be even like this:

				_launcher = launchers.NewContainerizedWeaveLauncher(info).Node()

				_launcher.Launch()

			*/
		},
	}
	return cmd
}
