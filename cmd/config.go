package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/klog"

	divingbell "github.com/tdaines42/diving-bell/pkg/diving-bell"
)

var write bool

func init() {
	configCmd.Flags().BoolVar(&write, "write", false, "Write the config to the local file")
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config [cluster name] [terraform workspace path]",
	Short: "Print the config to the console",
	Long:  ``,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if write {
			divingbell.UpdateClusterConfig(args[0], args[1])
		} else {
			bs, err := yaml.Marshal(divingbell.GenerateClusterConfig(args[0], args[1]))
			if err != nil {
				klog.Fatalf("unable to marshal config to YAML: %v", err)
			}

			fmt.Println(string(bs))
		}
	},
}
