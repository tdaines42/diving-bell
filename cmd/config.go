package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config related actions",
	Long:  ``,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	if write {
	// 		divingbell.UpdateClusterConfigFile(args[0], args[1])
	// 	} else if retrieve {
	// 		divingbell.RetrieveClusterConfig(args[0])
	// 	} else if update {
	// 		divingbell.StoreClusterConfig(args[0], args[1])
	// 	} else {
	// 		fmt.Println(divingbell.ClusterConfigYamlString(args[0], args[1]))
	// 	}
	// },
}
