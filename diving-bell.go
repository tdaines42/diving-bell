package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
	"k8s.io/klog"

	cluster "github.com/SUSE/skuba/pkg/skuba/actions/cluster/init"
)

type node struct {
	User     string
	Target   string
	HostName string `yaml:"hostName"`
}
type clusterConfig struct {
	ClusterName        string `yaml:"clusterName"`
	ControlPlaneTarget string `yaml:"controlPlaneTarget"`
	Managers           []node
	Workers            []node
}

func newClusterConfig() *clusterConfig {
	var config clusterConfig
	reader, _ := os.Open("cluster-config.yaml")
	buf, _ := ioutil.ReadAll(reader)
	yaml.Unmarshal(buf, &config)

	return &config
}

func initCluster(clusterName string, controlPlaneTarget string) {
	// Get current user
	usr, err := user.Current()
	if err != nil {
		klog.Fatalf("getting current user failed: %s", err)
	}

	// Init the cluster
	klog.Infof("Creating cluster %s\n", clusterName)

	initConfig, err := cluster.NewInitConfiguration(
		path.Join(usr.HomeDir, clusterName),
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
	klog.Infoln()
}

func bootstrapControlPlane(firstMaster node) {
	klog.Infof("Bootstrapping %s %s\n", firstMaster.HostName, firstMaster.Target)
	cmd := fmt.Sprintf("skuba node bootstrap --user %s --sudo --target %s %s", firstMaster.User, firstMaster.Target, firstMaster.HostName)
	runShell(cmd)
}

func joinNodes(nodes []node, role string) {

	for _, node := range nodes {
		node := node
		klog.Infof("Joining %s %s\n", node.HostName, node.Target)
		cmd := fmt.Sprintf("skuba node join --role %s --user %s --sudo --target %s %s", role, node.User, node.Target, node.HostName)
		runShell(cmd)
		time.Sleep(10 * time.Second)
	}
}

func main() {
	config := newClusterConfig()

	initCluster(config.ClusterName, config.ControlPlaneTarget)
	bootstrapControlPlane(config.Managers[0])

	if len(config.Managers) > 1 {
		joinNodes(config.Managers[1:len(config.Managers)], "master")
	}

	joinNodes(config.Workers, "worker")
}
