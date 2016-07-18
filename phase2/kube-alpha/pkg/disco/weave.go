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

func (*WeaveDisco) Bootstrap(peers []string) {
	// weave launch $peers
	// TODO retry_until_success, probably
	fmt.Println("================================")
	fmt.Println("  Kubernetes cluster bootstrap  ")
	fmt.Println("================================")
	fmt.Println("Installing bootstrap network...")
	logCommand("0001_bootstrap_install_weave",
		"curl", "-L", "git.io/weave", "-o", "/usr/local/bin/weave")
	logCommand("0002_bootstrap_chmod_weave",
		"chmod", "+x", "/usr/local/bin/weave")
	log.Println("done!")
	args := []string{"launch"}
	args = append(args, peers...)
	fmt.Println(`
Bootstrapping will now block until all servers join the bootstrap network.
Please run 'kube join <ip1>,<...>,<ipN>' on all the servers you want in your
initial cluter, and then wait for up to 2 minutes for bootstrapping to
complete...`)
	logCommand("0003_bootstrap_launch_weave",
		"/usr/local/bin/weave", args...)
	hostname, _ := os.Hostname()
	logCommand("0004_bootstrap_weave_expose",
		"/usr/local/bin/weave", "expose", "-h", hostname+".weave.local")
	fmt.Println("Bootstrap network successfully created!")
}

func NewWeaveDisco() P2PDiscovery {
	return &WeaveDisco{}
}
