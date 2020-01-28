package divingbell

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"k8s.io/klog"

	cluster "github.com/SUSE/skuba/pkg/skuba/actions/cluster/init"
	"github.com/tdaines42/diving-bell/internal/pkg/util"
)

func initCluster(clusterName string, controlPlaneTarget string) {
	// Get current user
	usr, err := user.Current()
	if err != nil {
		klog.Fatalf("getting current user failed: %s", err)
	}

	clusterConfigDir := path.Join(usr.HomeDir, clusterName)

	_, statErr := os.Stat(clusterConfigDir)
	if statErr == nil {
		os.RemoveAll(clusterConfigDir)
	}
	// Init the cluster
	klog.Infof("Creating cluster %s\n", clusterName)

	initConfig, err := cluster.NewInitConfiguration(
		clusterConfigDir,
		"",
		controlPlaneTarget,
		"",
		false)
	if err != nil {
		klog.Fatalf("init failed due to error: %s", err)
	}

	if err = cluster.Init(initConfig); err != nil {
		klog.Fatalf("init failed due to error: %s", err)
	}
	klog.Infoln()
}

func bootstrapControlPlane(firstMaster clusterNode, clusterName string) {
	klog.Infof("Bootstrapping %s %s\n", firstMaster.HostName, firstMaster.Target)
	cmd := fmt.Sprintf("skuba node bootstrap --user %s --sudo --target %s %s", firstMaster.User, firstMaster.Target, firstMaster.HostName)
	if util.RunShell(cmd) != true {
		klog.Fatalln("Failed to bootstrap the first master!")
	}

	CheckClusterReady(clusterName)
}

func joinNode(node clusterNode, role string) {
	klog.Infof("Joining %s %s\n", node.HostName, node.Target)
	cmd := fmt.Sprintf("skuba node join --role %s --user %s --sudo --target %s %s", role, node.User, node.Target, node.HostName)
	if util.RunShell(cmd) != true {
		klog.Fatalf("Failed to join %s to the cluster!\n", node.HostName)
	}
}

func joinManagers(nodes []clusterNode, clusterName string) {
	for _, node := range nodes {
		joinNode(node, "master")
		CheckClusterReady(clusterName)
	}
}

func joinWorkers(nodes []clusterNode, clusterName string) {
	for _, node := range nodes {
		joinNode(node, "worker")
	}

	CheckClusterReady(clusterName)
}

// BootstrapCluster Uses the config to bootstrap the cluster
func BootstrapCluster(config ClusterConfig) {
	initCluster(config.ClusterName, config.ControlPlaneTarget)
	bootstrapControlPlane(config.Managers[0], config.ClusterName)

	if len(config.Managers) > 1 {
		joinManagers(config.Managers[1:len(config.Managers)], config.ClusterName)
	}

	joinWorkers(config.Workers, config.ClusterName)
}
