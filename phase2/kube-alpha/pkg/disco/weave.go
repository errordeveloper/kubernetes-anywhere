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
	log.Printf("Starting to run %s %s...", cmd, args)
	output, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		thisLog.Printf("Execution failed with %s", err)
		log.Printf("Execution failed with %s", err)
	} else {
		thisLog.Printf("Execution succeeded")
		log.Printf("Execution succeeded")
	}
	thisLog.Printf("OUTPUT FOLLOWS\n==============")
	thisLog.Print(string(output))
	if err != nil {
		panic(fmt.Sprintf("Command failed, what to do?\nCommand: %s %s\nError: %s\nOutput: %s", cmd, args, output, err))
	}
}

func (*WeaveDisco) Bootstrap(peers []string) {
	// weave launch $peers
	// TODO retry_until_success, probably
	logCommand("0001_bootstrap_install_weave",
		"curl", "-L", "git.io/weave", "-o", "/usr/local/bin/weave")
	logCommand("0002_bootstrap_chmod_weave",
		"chmod", "+x", "/usr/local/bin/weave")
	args := []string{"launch"}
	args = append(args, peers...)
	logCommand("0003_bootstrap_launch_weave",
		"/usr/local/bin/weave", args...)
	hostname, _ := os.Hostname()
	log.Println("Waiting for cluster memers to join...", peers)
	logCommand("0003_bootstrap_weave_expose",
		"/usr/local/bin/weave", "expose", "-h", hostname+".weave.local")
}

func NewWeaveDisco() P2PDiscovery {
	return &WeaveDisco{}
}
