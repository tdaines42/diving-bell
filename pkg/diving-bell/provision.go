package divingbell

import (
	"k8s.io/klog"

	"github.com/tdaines42/diving-bell/internal/pkg/util"
)

// ProvisionCluster provision a cluster using terraform
func ProvisionCluster(workspacePath string) {
	if util.RunShellAt("terraform init -input=false", workspacePath).ExitCode != 0 {
		klog.Fatalln("Failed terraform init!")
	}
	if util.RunShellAt("terraform apply -input=false -auto-approve", workspacePath).ExitCode != 0 {
		klog.Fatalln("Failed terraform apply!")
	}
}
