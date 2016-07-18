package disco

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type WeaveDisco struct {
}

func getLogFile(logfile string) *os.File {
	f, err := os.OpenFile(
		fmt.Sprintf("%s.log", logfile),
		os.O_RDWR|os.O_CREATE, 0666,
	)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return f
}

// TODO - support multi master
// TODO - take arguments (master, worker, or master-and-worker)

// TODO support retries
func logCommand(logname, cmd string, args ...string) {
	// TODO don't create a log instance every time this gets run
	thisLog := log.New(getLogFile(logname), "", 0)
	thisLog.Printf("Starting to run %s %s...", cmd, args)
	output, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		thisLog.Printf("Execution failed with %s", err)
	} else {
		thisLog.Printf("Execution succeeded")
	}
	thisLog.Printf("OUTPUT FOLLOWS\n==============")
	thisLog.Print(string(output))
	if err != nil {
		os.Stderr.Write([]byte(fmt.Sprintf("Command failed, what to do?\nCommand: %s %s\nError: %s\nOutput: %s\n", cmd, args, output, err)))
		os.Exit(-1)
	}
}

// TODO retry_until_success, probably

func (w *WeaveDisco) Init(peers []string) {
	w.welcomeText()
	w.installNetwork()
	fmt.Println(`
Bootstrapping will now block until all servers join the
network.  Please run:

    kube join <ip1>,<...>,<ipN>

On all the other servers you want in your initial cluster,
giving the IP addresses of all the servers, and then wait
for up to 2 minutes for network bootstrapping to
complete...
`)
	w.launchWeave(peers)
	fmt.Println("Bootstrap network successfully created!")
	// TODO maybe this belongs in a different interface?
	//
	// we only do this on the first node for now
	w.initPKI()
	w.bootMaster()
}

func (w *WeaveDisco) Join(peers []string) {
	w.welcomeText()
	w.installNetwork()
	fmt.Println("Joining bootstrap network...")
	w.launchWeave(peers)
	fmt.Println("done!")
	w.doWorkerPKI()
	w.bootWorker()
}

func (w *WeaveDisco) welcomeText() {
	fmt.Println("================================")
	fmt.Println("  Kubernetes cluster bootstrap  ")
	fmt.Println("================================")
}

func (w *WeaveDisco) installNetwork() {
	fmt.Println("Installing bootstrap network...")
	logCommand("0001_bootstrap_install_weave",
		"curl", "-L", "git.io/weave", "-o", "/usr/local/bin/weave")
	logCommand("0002_bootstrap_chmod_weave",
		"chmod", "+x", "/usr/local/bin/weave")
	log.Println("done!")
}

func (w *WeaveDisco) launchWeave(peers []string) {
	args := []string{"launch"}
	args = append(args, peers...)
	logCommand("0003_bootstrap_launch_weave",
		"/usr/local/bin/weave", args...)
	hostname, _ := os.Hostname()
	logCommand("0004_bootstrap_weave_expose",
		"/usr/local/bin/weave", "expose", "-h", hostname+".weave.local")
}

// TODO we can probably implement toobox helper instead of having all of this in here
const loadFromTmpPkiServer = "set -o pipefail ; until docker --host=unix:///var/run/weave/weave.sock run --rm=true weaveworks/kubernetes-anywhere:toolbox curl --silent --fail tmp-pki-server/%s.dkr | docker load ; do sleep 1 ; done"

func (w *WeaveDisco) initPKI() {
	logCommand("0005_setup_pki",
		"docker", "run", "--rm=true", "--volume=/var/run/docker.sock:/docker.sock",
		"weaveworks/kubernetes-anywhere:toolbox", "create-pki-containers")

	logCommand("0006_create_tmp_pki_dir", "mkdir", "-p", "/tmp/pki")
	logCommand("0007_dump_toolbox_image",
		"docker", "save", "-o", "/tmp/pki/toolbox.dkr", "kubernetes-anywhere:toolbox-pki")
	logCommand("0008_dump_scheduler_image",
		"docker", "save", "-o", "/tmp/pki/scheduler.dkr", "kubernetes-anywhere:scheduler-pki")
	logCommand("0009_dump_controller_manager_image",
		"docker", "save", "-o", "/tmp/pki/controller-manager.dkr", "kubernetes-anywhere:controller-manager-pki")
	logCommand("0010_dump_kubelet_image",
		"docker", "save", "-o", "/tmp/pki/kubelet.dkr", "kubernetes-anywhere:kubelet-pki")
	logCommand("0011_dump_proxy_image",
		"docker", "save", "-o", "/tmp/pki/proxy.dkr", "kubernetes-anywhere:proxy-pki")
	logCommand("0012_dump_apiserver_image", "docker", "save", "-o", "/tmp/pki/apiserver.dkr", "kubernetes-anywhere:apiserver-pki")
	logCommand("0013_allow_nginx_read_access", "chmod", "o+r", "-R", "/tmp/pki")
	logCommand("0014_start_nginx",
		"docker", "--host=unix:///var/run/weave/weave.sock", "run",
		"--name=tmp-pki-server", "--volume=/tmp/pki:/usr/share/nginx/html:ro", "--detach=true", "nginx")
	logCommand("0015_init_toolbox_pki_container", "docker", "run", "--name=kube-toolbox-pki", "kubernetes-anywhere:toolbox-pki")
	logCommand("0016_init_apiserver_pki_container", "docker", "run", "--name=kube-apiserver-pki", "kubernetes-anywhere:apiserver-pki")
	logCommand("0017_init_scheduler_pki_container", "docker", "run", "--name=kube-scheduler-pki", "kubernetes-anywhere:scheduler-pki")
	logCommand("0018_init_controller_manager_pki_container", "docker", "run", "--name=kube-controller-manager-pki", "kubernetes-anywhere:controller-manager-pki")
}

func (w *WeaveDisco) doMasterPKI() {
	logCommand("0005_load_toolbox_pki_image", "sh", "-c", fmt.Sprintf(loadFromTmpPkiServer, "toolbox"))
	logCommand("0006_load_apiserver_pki_image", "sh", "-c", fmt.Sprintf(loadFromTmpPkiServer, "apiserver"))
	logCommand("0007_load_scheduler_pki_image", "sh", "-c", fmt.Sprintf(loadFromTmpPkiServer, "scheduler"))
	logCommand("0008_load_controller_manager_pki_image", "sh", "-c", fmt.Sprintf(loadFromTmpPkiServer, "controller-manager"))
	logCommand("0009_init_toolbox_pki_container", "docker", "run", "--name=kube-toolbox-pki", "kubernetes-anywhere:toolbox-pki")
	logCommand("0010_init_apiserver_pki_container", "docker", "run", "--name=kube-apiserver-pki", "kubernetes-anywhere:apiserver-pki")
	logCommand("0011_init_scheduler_pki_container", "docker", "run", "--name=kube-scheduler-pki", "kubernetes-anywhere:scheduler-pki")
	logCommand("0012_init_controller_manager_pki_container", "docker", "run", "--name=kube-controller-manager-pki", "kubernetes-anywhere:controller-manager-pki")
}

func (w *WeaveDisco) doWorkerPKI() {
	logCommand("0005_load_toolbox_pki_image", "sh", "-c", fmt.Sprintf(loadFromTmpPkiServer, "toolbox"))
	logCommand("0006_load_kubelet_pki_image", "sh", "-c", fmt.Sprintf(loadFromTmpPkiServer, "kubelet"))
	logCommand("0007_load_proxy_pki_image", "sh", "-c", fmt.Sprintf(loadFromTmpPkiServer, "proxy"))
	logCommand("0008_init_toolbox_pki_container", "docker", "run", "--name=kube-toolbox-pki", "kubernetes-anywhere:toolbox-pki")
	logCommand("0009_init_proxy_pki_container", "docker", "run", "--name=kube-proxy-pki", "kubernetes-anywhere:proxy-pki")
	logCommand("0010_init_kubelet_pki_container", "docker", "run", "--name=kubelet-pki", "kubernetes-anywhere:kubelet-pki")
}

func (w *WeaveDisco) bootMaster() {
	logCommand("0020_start_etcd",
		"docker", "--host=unix:///var/run/weave/weave.sock", "run", "--detach=true", "--name=etcd1", "weaveworks/kubernetes-anywhere:etcd")
	logCommand("0021_start_apiserver",
		"docker", "--host=unix:///var/run/weave/weave.sock", "run", "--detach=true", "--name=kube-apiserver", "--volumes-from=kube-apiserver-pki", "weaveworks/kubernetes-anywhere:apiserver")
	logCommand("0022_start_controller_manager",
		"docker", "--host=unix:///var/run/weave/weave.sock", "run", "--detach=true", "--name=kube-controller-manager", "--volumes-from=kube-controller-manager-pki", "weaveworks/kubernetes-anywhere:controller-manager")
	logCommand("0023_start_scheduler",
		"docker", "--host=unix:///var/run/weave/weave.sock", "run", "--detach=true", "--name=kube-scheduler", "--volumes-from=kube-scheduler-pki", "weaveworks/kubernetes-anywhere:scheduler")
}

func (w *WeaveDisco) bootWorker() {
	logCommand("0020_setup_kubelet_volumes",
		"docker", "run", "--rm=true", "--volume=/:/rootfs", "--volume=/var/run/docker.sock:/docker.sock", "weaveworks/kubernetes-anywhere:toolbox", "setup-kubelet-volumes")
	logCommand("0021_start_kubelet",
		"docker", "--host=unix:///var/run/weave/weave.sock", "run", "--detach=true", "--name=kubelet", "--privileged=true", "--net=host", "--pid=host", "--volumes-from=kubelet-volumes", "--volumes-from=kubelet-pki", "weaveworks/kubernetes-anywhere:kubelet")
	logCommand("0022_start_proxy",
		"docker", "--host=unix:///var/run/weave/weave.sock", "run", "--detach", "--name=kube-proxy", "--privileged=true", "--net=host", "--pid=host", "--volumes-from=kube-proxy-pki", "weaveworks/kubernetes-anywhere:proxy")
}

func NewWeaveDisco() P2PDiscovery {
	return &WeaveDisco{}
}
