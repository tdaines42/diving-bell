package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog"

	divingbell "github.com/tdaines42/diving-bell/pkg/diving-bell"
)

func init() {
	rootCmd.AddCommand(bootstrapCmd)
}

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap the cluster using skuba",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var config divingbell.ClusterConfig

		err := viper.Unmarshal(&config)
		if err != nil {
			klog.Fatalf("unable to decode into struct, %v", err)
		}

		divingbell.BootstrapCluster(config, false)
	},
}
