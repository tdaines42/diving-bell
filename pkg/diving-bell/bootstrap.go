package divingbell

import (
	"fmt"
	"os"
	"path"
	"time"

	"k8s.io/klog"

	cluster "github.com/SUSE/skuba/pkg/skuba/actions/cluster/init"
	"github.com/tdaines42/diving-bell/internal/pkg/util"
)

func initCluster(clusterConfigDir string, clusterName string, controlPlaneTarget string, kubernetesVersion string, destroy bool) {
	if destroy {
		_, statErr := os.Stat(clusterConfigDir)
		if statErr == nil {
			os.RemoveAll(clusterConfigDir)
		}
	}

	// Init the cluster
	klog.Infof("Creating cluster %s\n", clusterName)

	initConfig, err := cluster.NewInitConfiguration(
		clusterConfigDir,
		"",
		controlPlaneTarget,
		kubernetesVersion,
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
}

func joinNode(node clusterNode, role string) {
	klog.Infof("Joining %s %s\n", node.HostName, node.Target)
	cmd := fmt.Sprintf("skuba node join --role %s --user %s --sudo --target %s %s", role, node.User, node.Target, node.HostName)
	if util.RunShell(cmd) != true {
		klog.Fatalf("Failed to join %s to the cluster!\n", node.HostName)
	}
}

func joinManagers(nodes []clusterNode, clusterName string, kubeconfig string) {
	for _, node := range nodes {
		joinNode(node, "master")
		CheckClusterReady(kubeconfig)
	}
}

func joinWorkers(nodes []clusterNode, clusterName string, kubeconfig string) {
	for _, node := range nodes {
		joinNode(node, "worker")
	}

	CheckClusterReady(kubeconfig)
}

func createInitialCluster(firstMaster clusterNode, firstWorker clusterNode, clusterName string, kubeconfig string) {
	bootstrapControlPlane(firstMaster, clusterName)
	joinNode(firstWorker, "worker")

	time.Sleep(5 * time.Second)
	CheckClusterReady(kubeconfig)
}

// BootstrapCluster Uses the config to bootstrap the cluster
func BootstrapCluster(config ClusterConfig, currentWorkingDir string, destroy bool, kubeconfig string) {
	clusterConfigDir := path.Join(currentWorkingDir, config.ClusterName)

	initCluster(clusterConfigDir, config.ClusterName, config.ControlPlaneTarget, config.KubernetesVersion, destroy)
	createInitialCluster(config.Managers[0], config.Workers[0], config.ClusterName, kubeconfig)

	// Join additional nodes
	if len(config.Managers) > 1 {
		joinManagers(config.Managers[1:len(config.Managers)], config.ClusterName, kubeconfig)
	}

	if len(config.Workers) > 1 {
		joinWorkers(config.Workers[1:len(config.Workers)], config.ClusterName, kubeconfig)
	}
}
