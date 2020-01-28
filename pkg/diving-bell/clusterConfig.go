package divingbell

import (
	"encoding/json"

	"github.com/spf13/viper"
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

type outputList struct {
	Sensitive bool
	Type      string
	Value     []string
}
type outputString struct {
	Sensitive bool
	Type      string
	Value     string
}

type terraformOutput struct {
	HostnameMasters outputList   `json:"hostnames_masters"`
	HostnameWorkers outputList   `json:"hostnames_workers"`
	IPLoadBalancer  outputString `json:"ip_load_balancer"`
	IPMasters       outputList   `json:"ip_masters"`
	IPWorkers       outputList   `json:"ip_workers"`
}

// GenerateClusterConfig generate a struct of the cluster config
func GenerateClusterConfig(clusterName string, terraformWorkspacePath string) *ClusterConfig {
	var clusterConfig ClusterConfig
	clusterNodes := clusterNodesFromTerraform(terraformWorkspacePath)

	clusterConfig.ClusterName = clusterName
	clusterConfig.ControlPlaneTarget = clusterNodes.LoadBalancer
	clusterConfig.TerraformWorkspacePath = terraformWorkspacePath
	clusterConfig.Managers = clusterNodes.Managers
	clusterConfig.Workers = clusterNodes.Workers

	return &clusterConfig
}

// UpdateClusterConfig update the config file
func UpdateClusterConfig(clusterName string, terraformWorkspacePath string) {
	clusterNodes := clusterNodesFromTerraform(viper.GetString("terraformWorkspacePath"))

	viper.Set("controlPlaneTarget", clusterNodes.LoadBalancer)
	viper.Set("managers", clusterNodes.Managers)
	viper.Set("workers", clusterNodes.Workers)
	viper.WriteConfig()
}

// ClusterNodesFromTerraform get cluster nodes from terraform output
func clusterNodesFromTerraform(terraformWorkspacePath string) *ClusterNodes {
	var tOutput terraformOutput
	var clusterNodes ClusterNodes

	outputResults := util.RunShellAt("terraform output -json", terraformWorkspacePath)

	if outputResults.ExitCode != 0 {
		klog.Fatalln("Failed terraform output!")
	}

	json.Unmarshal([]byte(outputResults.Output), &tOutput)
	clusterNodes.LoadBalancer = tOutput.IPLoadBalancer.Value

	for i := 0; i < len(tOutput.IPMasters.Value) && i < len(tOutput.HostnameMasters.Value); i++ {
		node := clusterNode{User: "sles", Target: tOutput.IPMasters.Value[i], HostName: tOutput.HostnameMasters.Value[i]}
		clusterNodes.Managers = append(clusterNodes.Managers, node)
	}

	for i := 0; i < len(tOutput.IPWorkers.Value) && i < len(tOutput.HostnameWorkers.Value); i++ {
		node := clusterNode{User: "sles", Target: tOutput.IPWorkers.Value[i], HostName: tOutput.HostnameWorkers.Value[i]}
		clusterNodes.Workers = append(clusterNodes.Workers, node)
	}

	return &clusterNodes
}
