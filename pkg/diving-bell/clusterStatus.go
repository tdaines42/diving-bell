package divingbell

import (
	"fmt"
	"path"
	"os/user"

	"k8s.io/klog"

	"github.com/tdaines42/diving-bell/internal/pkg/util"
)

func checkNodesReady(clusterName string) {
	usr, err := user.Current()
	if err != nil {
		klog.Fatalf("getting current user failed: %s", err)
	}

	cmd := fmt.Sprintf("kubectl --kubeconfig=%s wait --for=condition=ready nodes --all --timeout 5m", path.Join(usr.HomeDir, clusterName, "admin.conf"))

	if util.RunShell(cmd) != true {
		klog.Fatalln("Failed waiting for the nodes to be ready!")
	}

}

func checkPodsReady(clusterName string) {
	usr, err := user.Current()
	if err != nil {
		klog.Fatalf("getting current user failed: %s", err)
	}

	cmd := fmt.Sprintf("kubectl --kubeconfig=%s wait --for=condition=ready pods -n kube-system --all --timeout 5m", path.Join(usr.HomeDir, clusterName, "admin.conf"))

	if util.RunShell(cmd) != true {
		klog.Fatalln("Failed waiting for the pods to be ready!")
	}
	
}

// CheckClusterReady check that cluster is ready
func CheckClusterReady(clusterName string) {
	checkNodesReady(clusterName)
	checkPodsReady(clusterName)
}
