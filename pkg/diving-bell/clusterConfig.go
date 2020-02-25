package divingbell

import (
	"encoding/json"
	"fmt"
	"os/user"
	"path"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/klog"

	"github.com/tdaines42/diving-bell/internal/pkg/util"
)

type clusterNode struct {
	User     string
	Target   string
	HostName string
}

// ClusterConfig struct representing a config
type ClusterConfig struct {
	ClusterName            string `yaml:"clusterName"`
	ControlPlaneTarget     string `yaml:"controlPlaneTarget"`
	KubernetesVersion      string `yaml:"kubernetesVersion"`
	TerraformWorkspacePath string `yaml:"terraformWorkspacePath"`
	Managers               []clusterNode
	Workers                []clusterNode
}

// ClusterNodes node info
type ClusterNodes struct {
	LoadBalancer string
	Managers     []clusterNode
	Workers      []clusterNode
}

type outputMap struct {
	Sensitive bool
	Type      string
	Value     map[string]string
}

type terraformOutput struct {
	IPLoadBalancer outputMap `json:"ip_load_balancer"`
	IPMasters      outputMap `json:"ip_masters"`
	IPWorkers      outputMap `json:"ip_workers"`
}

// ClusterConfigYamlString generate a string of the config
func ClusterConfigYamlString(clusterName string, kubernetesVersion string, terraformWorkspacePath string) string {
	UpdateClusterConfig(clusterName, kubernetesVersion, terraformWorkspacePath)

	bs, err := yaml.Marshal(viper.AllSettings())
	if err != nil {
		klog.Fatalf("unable to marshal config to YAML: %v", err)
	}

	return string(bs)
}

// RetrieveClusterConfig retrieve the config from the cluster
func RetrieveClusterConfig(clusterName string) {
	usr, err := user.Current()
	if err != nil {
		klog.Fatalf("getting current user failed: %s", err)
	}

	cmd := fmt.Sprintf("kubectl --kubeconfig=%s get configmap diving-bell -o jsonpath='{.data.\\.diving-bell\\.yaml}'", path.Join(usr.HomeDir, clusterName, "admin.conf"))
	out := util.RunShellOutput(cmd)
	if out.Error != nil {
		klog.Fatalln(out.Error)
	}
	fmt.Println(out.Output[1 : len(out.Output)-1])
}

// StoreClusterConfig store the config in the cluster as a config map
func StoreClusterConfig(clusterName string) {
	usr, err := user.Current()
	if err != nil {
		klog.Fatalf("getting current user failed: %s", err)
	}

	cmd := fmt.Sprintf("kubectl --kubeconfig=%s delete configmap diving-bell", path.Join(usr.HomeDir, clusterName, "admin.conf"))
	util.RunShellOutput(cmd)
	cmd = fmt.Sprintf("kubectl --kubeconfig=%s create configmap diving-bell --from-file=%s", path.Join(usr.HomeDir, clusterName, "admin.conf"), viper.ConfigFileUsed())
	out := util.RunShellOutput(cmd)

	if out.Error != nil {
		klog.Fatalln(out.Error)
	}
}

// UpdateClusterConfig update the config
func UpdateClusterConfig(clusterName string, kubernetesVersion string, terraformWorkspacePath string) {
	clusterNodes := clusterNodesFromTerraform(clusterName, terraformWorkspacePath)

	viper.Set("clusterName", clusterName)
	viper.Set("kubernetesVersion", kubernetesVersion)
	viper.Set("terraformWorkspacePath", terraformWorkspacePath)
	viper.Set("controlPlaneTarget", clusterNodes.LoadBalancer)
	viper.Set("managers", clusterNodes.Managers)
	viper.Set("workers", clusterNodes.Workers)
}

// UpdateClusterConfigFile update the config file
func UpdateClusterConfigFile(clusterName string, kubernetesVersion string, terraformWorkspacePath string) {
	UpdateClusterConfig(clusterName, kubernetesVersion, terraformWorkspacePath)
	viper.WriteConfig()
}

func clusterNodesFromTerraform(clusterName string, terraformWorkspacePath string) *ClusterNodes {
	var tOutput terraformOutput
	var clusterNodes ClusterNodes

	outputResults := util.RunShellAt("terraform output -json", terraformWorkspacePath)

	if outputResults.ExitCode != 0 {
		klog.Fatalln("Failed terraform output!")
	}

	json.Unmarshal([]byte(outputResults.Output), &tOutput)

	for _, value := range tOutput.IPLoadBalancer.Value {
		clusterNodes.LoadBalancer = value
		break
	}

	for key, value := range tOutput.IPMasters.Value {
		node := clusterNode{User: "sles", Target: value, HostName: key}
		clusterNodes.Managers = append(clusterNodes.Managers, node)
	}

	for key, value := range tOutput.IPWorkers.Value {
		node := clusterNode{User: "sles", Target: value, HostName: key}
		clusterNodes.Workers = append(clusterNodes.Workers, node)
	}

	return &clusterNodes
}
