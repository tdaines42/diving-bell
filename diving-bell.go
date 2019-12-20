package main

import (
	"fmt"
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

func bootstrapControlPlane() {
	out, err := exec.Command("skuba", "node", "bootstrap", "--user", "sles", "--sudo", "--target", "10.17.2.0", "testing-master-0").Output()
	if err != nil {
		klog.Fatal(err)
	}
	klog.Info(out)
}

func main() {
	initCluster()
	bootstrapControlPlane()
}
