package divingbell

import (
	"fmt"
	"path"
	"os"

	"k8s.io/klog"

	"github.com/tdaines42/diving-bell/internal/pkg/util"
)

func checkNodesReady(clusterName string) {
	// Find current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		klog.Errorln(err)
		os.Exit(1)
	}

	cmd := fmt.Sprintf("kubectl --kubeconfig=%s wait --for=condition=ready nodes --all --timeout 5m", path.Join(cwd, clusterName, "admin.conf"))

	if util.RunShell(cmd) != true {
		klog.Fatalln("Failed waiting for the nodes to be ready!")
	}

}

func checkPodsReady(clusterName string) {
	// Find current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		klog.Errorln(err)
		os.Exit(1)
	}

	cmd := fmt.Sprintf("kubectl --kubeconfig=%s wait --for=condition=ready pods -n kube-system --all --timeout 5m", path.Join(cwd, clusterName, "admin.conf"))

	if util.RunShell(cmd) != true {
		klog.Fatalln("Failed waiting for the pods to be ready!")
	}
	
}

// CheckClusterReady check that cluster is ready
func CheckClusterReady(clusterName string) {
	checkNodesReady(clusterName)
	checkPodsReady(clusterName)
}
