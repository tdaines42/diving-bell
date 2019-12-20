package main

import (
	"bufio"
	"fmt"
	"strings"
	"os/exec"
	"os/user"

	"k8s.io/klog"

	cluster "github.com/SUSE/skuba/pkg/skuba/actions/cluster/init"
)

func initCluster() {
	// Get current user
	usr, err := user.Current()
	if err != nil {
		klog.Fatalf("getting current user failed: %s", err)
	}

	// Init the cluster
	initConfig, err := cluster.NewInitConfiguration(
		fmt.Sprintf("%s/test-cluster", usr.HomeDir),
		"",
		"10.17.1.0",
		"",
		false)
	if err != nil {
		klog.Fatalf("init failed due to error: %s", err)
	}

	if err = cluster.Init(initConfig); err != nil {
		klog.Fatalf("init failed due to error: %s", err)
	}
}

func runShell(shellCmd string) {
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
}

func bootstrapControlPlane() {
	runShell("skuba node bootstrap --user sles --sudo --target 10.17.2.0 testing-master-0")
}

func joinNodes() {
	runShell("skuba node join --user sles --sudo --target 10.17.3.0 --role worker testing-worker-0")
	runShell("skuba node join --user sles --sudo --target 10.17.3.1 --role worker testing-worker-1")
}

func main() {
	initCluster()
	bootstrapControlPlane()
	joinNodes()
}
