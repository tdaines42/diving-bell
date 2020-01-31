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

// GenerateClusterConfig generate a struct of the cluster config
func GenerateClusterConfig(clusterName string, terraformWorkspacePath string) *ClusterConfig {
	var clusterConfig ClusterConfig
	clusterNodes := clusterNodesFromTerraform(clusterName, terraformWorkspacePath)

	clusterConfig.ClusterName = clusterName
	clusterConfig.ControlPlaneTarget = clusterNodes.LoadBalancer
	clusterConfig.TerraformWorkspacePath = terraformWorkspacePath
	clusterConfig.Managers = clusterNodes.Managers
	clusterConfig.Workers = clusterNodes.Workers

	return &clusterConfig
}

// UpdateClusterConfig update the config file
func UpdateClusterConfig(clusterName string, terraformWorkspacePath string) {
	clusterNodes := clusterNodesFromTerraform(clusterName, terraformWorkspacePath)

	viper.Set("clusterName", clusterName)
	viper.Set("terraformWorkspacePath", terraformWorkspacePath)
	viper.Set("controlPlaneTarget", clusterNodes.LoadBalancer)
	viper.Set("managers", clusterNodes.Managers)
	viper.Set("workers", clusterNodes.Workers)
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
