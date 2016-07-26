package discovery_providers

import (
	. "github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/interfaces"
	"github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/util"
	"os"
)

type WeaveFirstMasterDiscoveryProvider struct {
	Info ClusterInfo
}

type WeaveAnyMasterDiscoveryProvider struct {
	Info ClusterInfo
}

type WeaveNodeDiscoveryProvider struct {
	Info ClusterInfo
}

func NewWeaveFirstMasterDiscoveryProvider(info ClusterInfo) FirstMasterDiscoveryProvider {
	return &WeaveFirstMasterDiscoveryProvider{Info: info}
}

func NewWeaveAnyMasterDiscoveryProvider(info ClusterInfo) AnyMasterDiscoveryProvider {
	return &WeaveAnyMasterDiscoveryProvider{Info: info}
}

func NewWeaveNodeDiscoveryProvider(info ClusterInfo) NodeDiscoveryProvider {
	return &WeaveNodeDiscoveryProvider{Info: info}
}

func (w *WeaveFirstMasterDiscoveryProvider) Setup() {
	installWeave()
}

func (w *WeaveFirstMasterDiscoveryProvider) Launch() {
	// TODO: we should pass `--ipalloc-init consensus=1` here, if we
	// know that this cluster is only gonna have a single master
	launchWeave(nil, nil, append(w.Info.MasterIPs, w.Info.NodeIPs...))
}

func (w *WeaveAnyMasterDiscoveryProvider) Setup() {
	installWeave()
}

func (w *WeaveAnyMasterDiscoveryProvider) Launch() {
	// TODO: we should pass `--ipalloc-init consensus=X` here
	launchWeave(nil, nil, append(w.Info.MasterIPs, w.Info.NodeIPs...))
}

func (w *WeaveNodeDiscoveryProvider) Setup() {
	installWeave()
}

func (w *WeaveNodeDiscoveryProvider) Launch() {
	// TODO: most likelly this going to pass `--ipalloc-init observer`
	launchWeave(nil, nil, append(w.Info.MasterIPs, w.Info.NodeIPs...))
}

func installWeave() {
	// TODO: this can be done through Docker API
	util.LogCommand("0001_bootstrap_install_weave",
		"curl", "-L", "git.io/weave", "-o", "/usr/local/bin/weave")
	util.LogCommand("0002_bootstrap_chmod_weave",
		"chmod", "+x", "/usr/local/bin/weave")
	// this pulls containers from the registry, and it will also install
	// CNI config and wrappers, if /opt/cni and /etc/cni exist, we should
	// probably disable that, to avoid accidents, but such are unlikelly
	// at this point in time
	util.LogCommand("0003_bootstrap_weave_setup",
		"/usr/local/bin/weave", "setup")
}

func launchWeave(routerArgs, proxyArgs, peers []string) {
	// TODO: this can be done through Docker API
	// "weave <args>" => "docker run --rm --privileged --net=host -v /var/run/docker.sock:/var/run/docker.sock --pid=host -v /:/host -e HOST_ROOT=/host weaveworks/weaveexec:1.6.0 --local <args>"
	launchRouter := []string{"launch-router"}
	launchRouter = append(launchRouter, routerArgs...)
	launchRouter = append(launchRouter, peers...)
	launchProxy := []string{"launch-proxy"}
	launchProxy = append(launchProxy, proxyArgs...)
	util.LogCommand("0003_bootstrap_weave_launch_router",
		"/usr/local/bin/weave", launchRouter...)
	util.LogCommand("0003_bootstrap_weave_launch_proxy",
		"/usr/local/bin/weave", launchProxy...)
	hostname, _ := os.Hostname()
	util.LogCommand("0004_bootstrap_weave_expose",
		"/usr/local/bin/weave", "expose", "-h", hostname+".weave.local")
}
