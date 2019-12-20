package main

import (
	"fmt"
	"os/user"

	"k8s.io/klog"

	"github.com/SUSE/skuba/pkg/skuba/actions/node/bootstrap"
	"github.com/SUSE/skuba/internal/pkg/skuba/deployments"
	"github.com/SUSE/skuba/internal/pkg/skuba/deployments/ssh"
	cluster "github.com/SUSE/skuba/pkg/skuba/actions/cluster/init"
	
)

func main() {
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

	// Bootstrap the first master
	bootstrapConfiguration := deployments.BootstrapConfiguration{
		KubeadmExtraArgs: map[string]string{"ignore-preflight-errors": ""},
	}

	target := ssh.Target{}
	target.user = "sles"
	target.targetName = "10.17.2.0"
	target.sudo = true
	target.port = 22

	role := deployments.MasterRole
	d := target.GetDeployment("testing-master-0", &role)
	if err := bootstrap.Bootstrap(bootstrapConfiguration, d); err != nil {
		klog.Fatalf("error bootstrapping node: %s", err)
	}
}
