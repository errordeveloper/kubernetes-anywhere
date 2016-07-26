package launchers

import (
	. "github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/interfaces"
	"github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/util"
)

type ContainerizedWeaveFirstMasterKubernetesLauncher struct {
	Info ClusterInfo
}

type ContainerizedWeaveAnyMasterKubernetesLauncher struct {
	Info ClusterInfo
}

type ContainerizedWeaveNodeKubernetesLauncher struct {
	Info ClusterInfo
}

func NewContainerizedWeaveFirstMasterLauncher(info ClusterInfo) FirstMasterLauncher {
	return &ContainerizedWeaveFirstMasterKubernetesLauncher{Info: info}
}

func NewContainerizedWeaveAnyMasterLauncher(info ClusterInfo) AnyMasterLauncher {
	return &ContainerizedWeaveAnyMasterKubernetesLauncher{Info: info}
}

func NewContainerizedWeaveNodeLauncher(info ClusterInfo) NodeLauncher {
	return &ContainerizedWeaveNodeKubernetesLauncher{Info: info}
}

func (w *ContainerizedWeaveFirstMasterKubernetesLauncher) Launch() {
	// TODO in HA setting we will need to figure which etcdX we are starting...
	util.LogCommand("0020_start_etcd",
		"docker", "--host=unix:///var/run/weave/weave.sock", "run", "--detach=true", "--name=etcd1", "weaveworks/kubernetes-anywhere:etcd")
	util.LogCommand("0021_start_apiserver",
		"docker", "--host=unix:///var/run/weave/weave.sock", "run", "--detach=true", "--name=kube-apiserver", "--volumes-from=kube-apiserver-pki", "weaveworks/kubernetes-anywhere:apiserver")
	util.LogCommand("0022_start_controller_manager",
		"docker", "--host=unix:///var/run/weave/weave.sock", "run", "--detach=true", "--name=kube-controller-manager", "--volumes-from=kube-controller-manager-pki", "weaveworks/kubernetes-anywhere:controller-manager")
	util.LogCommand("0023_start_scheduler",
		"docker", "--host=unix:///var/run/weave/weave.sock", "run", "--detach=true", "--name=kube-scheduler", "--volumes-from=kube-scheduler-pki", "weaveworks/kubernetes-anywhere:scheduler")
}

func (w *ContainerizedWeaveAnyMasterKubernetesLauncher) Launch() {
}

func (w *ContainerizedWeaveNodeKubernetesLauncher) Launch() {
	util.LogCommand("0020_setup_kubelet_volumes",
		"docker", "run", "--rm=true", "--volume=/:/rootfs", "--volume=/var/run/docker.sock:/docker.sock", "weaveworks/kubernetes-anywhere:toolbox", "setup-kubelet-volumes")
	util.LogCommand("0021_start_kubelet",
		"docker", "--host=unix:///var/run/weave/weave.sock", "run", "--detach=true", "--name=kubelet", "--privileged=true", "--net=host", "--pid=host", "--volumes-from=kubelet-volumes", "--volumes-from=kubelet-pki", "weaveworks/kubernetes-anywhere:kubelet")
	util.LogCommand("0022_start_proxy", "docker", "--host=unix:///var/run/weave/weave.sock", "run", "--detach", "--name=kube-proxy", "--privileged=true", "--net=host", "--pid=host", "--volumes-from=kube-proxy-pki", "weaveworks/kubernetes-anywhere:proxy")
}
