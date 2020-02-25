package cmd

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog"
)

var cfgFile string
var kubernetesVersion string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "diving-bell",
	Short: "Manage a k8s cluster using kubectl, terraform, and skuba",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	var config divingbell.ClusterConfig

	// 	err := viper.Unmarshal(&config)
	// 	if err != nil {
	// 		klog.Fatalf("unable to decode into struct, %v", err)
	// 	}

	// 	divingbell.BootstrapCluster(config)
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		klog.Errorln(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.diving-bell.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			klog.Errorln(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".diving-bell" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".diving-bell")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		klog.Infoln("Using config file:", viper.ConfigFileUsed())
	}
}
