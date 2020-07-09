package divingbell

import (
	"fmt"

	"k8s.io/klog"

	"github.com/tdaines42/diving-bell/internal/pkg/util"
)

// ClusterStatus status of all the nodes in the cluster
type ClusterStatus struct {
	Managers []nodeStatus
	Workers  []nodeStatus
}

type nodeStatus struct {
	HostName     string
	Provisioned  bool
	Joined       bool
	KubectlState string
	SkubaState   string
}

func checkNodesReady(kubeconfig string) {
	cmd := fmt.Sprintf("kubectl --kubeconfig=%s wait --for=condition=ready nodes --all --timeout 5m", kubeconfig)

	if util.RunShell(cmd) != true {
		klog.Fatalln("Failed waiting for the nodes to be ready!")
	}

}

func checkPodsReady(kubeconfig string) {
	cmd := fmt.Sprintf("kubectl --kubeconfig=%s wait --for=condition=ready pods -n kube-system --all --timeout 5m", kubeconfig)

	if util.RunShell(cmd) != true {
		klog.Fatalln("Failed waiting for the pods to be ready!")
	}

}

// CheckClusterReady check that cluster is ready
func CheckClusterReady(kubeconfig string) {
	checkNodesReady(kubeconfig)
	checkPodsReady(kubeconfig)
}
