package pki

import (
	"fmt"
	. "github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/interfaces"
	"github.com/kubernetes/kubernetes-anywhere/phase2/kube-alpha/pkg/util"
)

type ContainerizedWeaveFirstMasterPKI struct {
	Info ClusterInfo
}

type ContainerizedWeaveAnyMasterPKI struct {
	Info ClusterInfo
}

type ContainerizedWeaveNodePKI struct {
	Info ClusterInfo
}

func NewContainerizedWeaveFirstMasterPKI(info ClusterInfo) FirstMasterPKI {
	return &ContainerizedWeaveFirstMasterPKI{Info: info}
}

func NewContainerizedWeaveAnyMasterPKI(info ClusterInfo) AnyMasterPKI {
	return &ContainerizedWeaveAnyMasterPKI{Info: info}
}

func NewContainerizedWeaveNodePKI(info ClusterInfo) NodePKI {
	return &ContainerizedWeaveNodePKI{Info: info}
}

// TODO we can probably implement toobox helper instead of having all of this in here
const loadFromTmpPkiServer = "set -o pipefail ; until docker --host=unix:///var/run/weave/weave.sock run --rm=true weaveworks/kubernetes-anywhere:toolbox curl --silent --fail tmp-pki-server/%s.dkr | docker load ; do sleep 1 ; done"

func (w *ContainerizedWeaveFirstMasterPKI) Init() {
	util.LogCommand("0005_setup_pki",
		"docker", "run", "--rm=true", "--volume=/var/run/docker.sock:/docker.sock",
		"weaveworks/kubernetes-anywhere:toolbox", "create-pki-containers")
}

func (w *ContainerizedWeaveFirstMasterPKI) Publish() {
	util.LogCommand("0006_create_tmp_pki_dir", "mkdir", "-p", "/tmp/pki")
	util.LogCommand("0007_dump_toolbox_image",
		"docker", "save", "-o", "/tmp/pki/toolbox.dkr", "kubernetes-anywhere:toolbox-pki")
	util.LogCommand("0008_dump_scheduler_image",
		"docker", "save", "-o", "/tmp/pki/scheduler.dkr", "kubernetes-anywhere:scheduler-pki")
	util.LogCommand("0009_dump_controller_manager_image",
		"docker", "save", "-o", "/tmp/pki/controller-manager.dkr", "kubernetes-anywhere:controller-manager-pki")
	util.LogCommand("0010_dump_kubelet_image",
		"docker", "save", "-o", "/tmp/pki/kubelet.dkr", "kubernetes-anywhere:kubelet-pki")
	util.LogCommand("0011_dump_proxy_image",
		"docker", "save", "-o", "/tmp/pki/proxy.dkr", "kubernetes-anywhere:proxy-pki")
	util.LogCommand("0012_dump_apiserver_image", "docker", "save", "-o", "/tmp/pki/apiserver.dkr", "kubernetes-anywhere:apiserver-pki")
	util.LogCommand("0013_allow_nginx_read_access", "chmod", "o+r", "-R", "/tmp/pki")
	util.LogCommand("0014_start_nginx",
		"docker", "--host=unix:///var/run/weave/weave.sock", "run",
		"--name=tmp-pki-server", "--volume=/tmp/pki:/usr/share/nginx/html:ro", "--detach=true", "nginx")
}

func (w *ContainerizedWeaveAnyMasterPKI) Fetch() {
	util.LogCommand("0005_load_toolbox_pki_image", "sh", "-c", fmt.Sprintf(loadFromTmpPkiServer, "toolbox"))
	util.LogCommand("0006_load_apiserver_pki_image", "sh", "-c", fmt.Sprintf(loadFromTmpPkiServer, "apiserver"))
	util.LogCommand("0007_load_scheduler_pki_image", "sh", "-c", fmt.Sprintf(loadFromTmpPkiServer, "scheduler"))
	util.LogCommand("0008_load_controller_manager_pki_image", "sh", "-c", fmt.Sprintf(loadFromTmpPkiServer, "controller-manager"))
}

func (w *ContainerizedWeaveAnyMasterPKI) Init() {
	util.LogCommand("0009_init_toolbox_pki_container", "docker", "run", "--name=kube-toolbox-pki", "kubernetes-anywhere:toolbox-pki")
	util.LogCommand("0010_init_apiserver_pki_container", "docker", "run", "--name=kube-apiserver-pki", "kubernetes-anywhere:apiserver-pki")
	util.LogCommand("0011_init_scheduler_pki_container", "docker", "run", "--name=kube-scheduler-pki", "kubernetes-anywhere:scheduler-pki")
	util.LogCommand("0012_init_controller_manager_pki_container", "docker", "run", "--name=kube-controller-manager-pki", "kubernetes-anywhere:controller-manager-pki")
}

func (w *ContainerizedWeaveNodePKI) Fetch() {
	util.LogCommand("0005_load_toolbox_pki_image", "sh", "-c", fmt.Sprintf(loadFromTmpPkiServer, "toolbox"))
	util.LogCommand("0006_load_kubelet_pki_image", "sh", "-c", fmt.Sprintf(loadFromTmpPkiServer, "kubelet"))
	util.LogCommand("0007_load_proxy_pki_image", "sh", "-c", fmt.Sprintf(loadFromTmpPkiServer, "proxy"))
}

func (w *ContainerizedWeaveNodePKI) Init() {
	util.LogCommand("0008_init_toolbox_pki_container", "docker", "run", "--name=kube-toolbox-pki", "kubernetes-anywhere:toolbox-pki")
	util.LogCommand("0009_init_proxy_pki_container", "docker", "run", "--name=kube-proxy-pki", "kubernetes-anywhere:proxy-pki")
	util.LogCommand("0010_init_kubelet_pki_container", "docker", "run", "--name=kubelet-pki", "kubernetes-anywhere:kubelet-pki")
}
