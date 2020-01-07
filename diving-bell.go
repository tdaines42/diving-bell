package main

import (
	"k8s.io/klog"

	"github.com/tdaines42/diving-bell/cmd"
)

func main() {
	klog.InitFlags(nil)
	cmd.Execute()
}
