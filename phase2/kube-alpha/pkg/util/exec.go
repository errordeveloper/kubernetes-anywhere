package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

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
func LogCommand(logname, cmd string, args ...string) {
	// TODO don't create a log instance every time this gets run
	thisLog := log.New(getLogFile(logname), "", 0)
	thisLog.Printf("Starting to run %s %s...", cmd, args)
	output, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		thisLog.Printf("Execution failed with %s", err)
	} else {
		thisLog.Printf("Execution succeeded")
	}
	thisLog.Println("OUTPUT FOLLOWS")
	thisLog.Println("==============")
	thisLog.Print(string(output))
	if err != nil {
		fmt.Errorf("Command failed, what to do?\nCommand: %s %s\nError: %s\nOutput: %s\n", cmd, args, output, err)
		os.Exit(-1)
	}
}

// TODO retry_until_success, probably
