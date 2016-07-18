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
	w.setupCertificateRepo()
}

func (w *WeaveDisco) Join(peers []string) {
	w.welcomeText()
	w.installNetwork()
	fmt.Println("Joining bootstrap network...")
	w.launchWeave(peers)
	fmt.Println("done!")
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

func (w *WeaveDisco) setupCertificateRepo() {
	logCommand("0005_setup_pki",
		"docker", "run", "-v", "/var/run/docker.sock:/docker.sock",
		"weaveworks/kubernetes-anywhere:toolbox", "create-pki-containers",
	)
	/*
		* on first master
			mkdir /tmp/pki
			docker save -i /tmp/pki/toolbox.dkr kubernetes-anywhere:toolbox-pki
			docker save -i /tmp/pki/scheduler.dkr kubernetes-anywhere:scheduler-pki
			docker save -i /tmp/pki/controller-manager.dkr kubernetes-anywhere:controller-manager-pki
			docker save -i /tmp/pki/kubelet.dkr kubernetes-anywhere:kubelet-pki
			docker save -i /tmp/pki/proxy.dkr kubernetes-anywhere:proxy-pki
			docker save -i /tmp/pki/apiserver.dkr kubernetes-anywhere:apiserver-pki
			chmod o+r -R /tmp/pki
			docker --host=unix:///var/run/weave/weave.sock run --name=tmp-pki-server --volume=/tmp/pki:/usr/share/nginx/html:ro --detach=true nginx

		* on other masters
			docker --host=unix:///var/run/weave/weave.sock run weaveworks/kubernetes-anywhere:toolbox curl tmp-pki-server/toolbox.dkr | docker load
			docker --host=unix:///var/run/weave/weave.sock run weaveworks/kubernetes-anywhere:toolbox curl tmp-pki-server/apiserver.dkr | docker load
			docker --host=unix:///var/run/weave/weave.sock run weaveworks/kubernetes-anywhere:toolbox curl tmp-pki-server/scheduler.dkr | docker load
			docker --host=unix:///var/run/weave/weave.sock run weaveworks/kubernetes-anywhere:toolbox curl tmp-pki-server/controller-manager.dkr | docker load
			docker run --name=kube-toolbox-pki weaveworks/kubernetes-anywhere:toolbox-pki
			docker run --name=kube-apiserver-pki weaveworks/kubernetes-anywhere:apiserver-pki
			docker run --name=kube-scheduler-pki weaveworks/kubernetes-anywhere:scheduler-pki
			docker run --name=kube-controller-manager-pki weaveworks/kubernetes-anywhere:controller-manager-pki

		* on nodes
			docker --host=unix:///var/run/weave/weave.sock run weaveworks/kubernetes-anywhere:toolbox curl tmp-pki-server/toolbox.dkr | docker load
			docker --host=unix:///var/run/weave/weave.sock run weaveworks/kubernetes-anywhere:toolbox curl tmp-pki-server/kubelet.dkr | docker load
			docker --host=unix:///var/run/weave/weave.sock run weaveworks/kubernetes-anywhere:toolbox curl tmp-pki-server/proxy.dkr | docker load
			docker run --name=kube-toolbox-pki weaveworks/kubernetes-anywhere:toolbox-pki
			docker run --name=kube-proxy-pki weaveworks/kubernetes-anywhere:kubelet-pki
			docker run --name=kubelet-pki weaveworks/kubernetes-anywhere:proxy-pki
	*/
}

func NewWeaveDisco() P2PDiscovery {
	return &WeaveDisco{}
}
