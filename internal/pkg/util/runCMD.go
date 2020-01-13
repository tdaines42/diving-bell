package util

import (
	"bufio"
	"os/exec"
	"strings"

	"k8s.io/klog"
)

// RunShell runs a shell cmd
func RunShell(shellCmd string) bool {
	args := strings.Fields(shellCmd)
	cmd := exec.Command(args[0], args[1:len(args)]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		klog.Fatal(err)
	}
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		klog.Infoln(scanner.Text())
	}

	cmd.Wait()

	klog.Infoln()

	return cmd.ProcessState.ExitCode() == 0
}
