package cmd

import (
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/discovery_providers"
	. "github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/interfaces"
	"github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/launchers"
	"github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/pki"
	"github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/util"
)

func NewCmdInit(out io.Writer, config *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Run this on the first server you deploy onto.",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO obviously this not ideal, we should have better flags
			// to reflect the shape of ClusterInfo
			info := ClusterInfo{MasterIPs: strings.Split(args[0], ",")}
			var (
				_disco    FirstMasterDiscoveryProvider
				_pki1     FirstMasterPKI
				_pki2     AnyMasterPKI
				_launcher FirstMasterLauncher
			)

			util.PrintMessage(`
				  Bootstrapping will now block until all servers join the
				  network.  Please run:

				      kube join <ip1>,<...>,<ipN>

				  On all the other servers you want in your initial cluster,
				  giving the IP addresses of all the servers, and then wait
				  for up to 2 minutes for network bootstrapping to
				  complete...
			`)

			_disco = discovery_providers.NewWeaveFirstMasterDiscoveryProvider(info)

			_disco.Setup()
			_disco.Launch()

			util.PrintMessage("Bootstrap network successfully created!")
			util.PrintOkay()

			util.PrintMessage("Initializing PKI...")

			_pki1 = pki.NewContainerizedWeaveFirstMasterPKI(info)

			_pki1.Init()
			_pki1.Publish()

			_pki2 = pki.NewContainerizedWeaveAnyMasterPKI(info)

			_pki2.Init()

			util.PrintOkay()

			util.PrintMessage("Booting master...")

			_launcher = launchers.NewContainerizedWeaveFirstMasterLauncher(info)

			_launcher.Launch()

			util.PrintOkay()

			// Alternative version, we could change these interfaces to be more like:
			/*
				info := ClusterInfo{MasterIPs: strings.Split(args[0], ",")}
				var (
					_disco    DiscoveryProvider
					_pki      PKI
					_launcher Launcher
				)

				_disco = discovery_providers.NewWeaveDiscoveryProvider(info)

				_disco.FirstMasterSetup()
				_disco.FirstMasterLaunch()

				_pki = pki.NewContainerizedWeavePKI(info)

				_pki.FirstMasterInit()
				_pki.FirstMasterPublish()

				_launcher = launchers.NewContainerizedWeaveLauncher(info)

				_launcher.FirstMasterLaunch()
			*/
		},
	}
	return cmd
}
