package divingbell

type clusterNode struct {
	User     string
	Target   string
	HostName string
}

// ClusterConfig struct representing a config
type ClusterConfig struct {
	ClusterName        string
	ControlPlaneTarget string
	Managers           []clusterNode
	Workers            []clusterNode
}
