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
				_pki      FirstMasterPKI
				_launcher FirstMasterLauncher
			)

			_disco = discovery_providers.NewWeaveFirstMasterDiscoveryProvider(info)

			_disco.Setup()
			_disco.Launch()

			_pki = pki.NewContainerizedWeaveFirstMasterPKI(info)

			_pki.Init()
			_pki.Publish()

			_launcher = launchers.NewContainerizedWeaveFirstMasterLauncher(info)

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
