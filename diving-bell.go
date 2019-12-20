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
	HostName string
}
type clusterConfig struct {
	ClusterName        string
	ControlPlaneTarget string
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

func bootstrapControlPlane(firstMaster node) {
	klog.Infof("Bootstrapping %+v\n", firstMaster)
	cmd := fmt.Sprintf("skuba node bootstrap --user %s --sudo --target %s %s", firstMaster.User, firstMaster.Target, firstMaster.HostName)
	runShell(cmd)
}

func joinNodes(nodes []node) {

	for _, node := range nodes {
		klog.Infof("Joining %+v\n", node)
		cmd := fmt.Sprintf("skuba node join --user %s --sudo --target %s %s", node.User, node.Target, node.HostName)
		runShell(cmd)
		time.Sleep(10 * time.Second)
	}
}

func main() {
	config := newClusterConfig()

	klog.Infof("Config %+v\n", config)

	initCluster(config.ClusterName, config.ControlPlaneTarget)
	bootstrapControlPlane(config.Managers[0])

	if len(config.Managers) > 1 {
		joinNodes(config.Managers[1:len(config.Managers)])
	}

	joinNodes(config.Workers)
}
