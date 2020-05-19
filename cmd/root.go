package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog"

	"github.com/tdaines42/diving-bell/internal/pkg/util"
)

var cfgFile string
var kubernetesVersion string
var currentWorkingDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "diving-bell",
	Short: "Manage a k8s cluster using kubectl, terraform, and skuba",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	//
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
	// Find current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		klog.Errorln(err)
		os.Exit(1)
	}

	currentWorkingDir = cwd

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is {cwd}/.diving-bell.yaml)")
	rootCmd.PersistentFlags().BoolVar(&util.Debug, "debug", false, "run in debug mode")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find current working directory.
		cwd, err := os.Getwd()
		if err != nil {
			klog.Errorln(err)
			os.Exit(1)
		}

		// Search config in cwd with name ".diving-bell" (without extension).
		viper.AddConfigPath(cwd)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".diving-bell")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		klog.Infoln("Using config file:", viper.ConfigFileUsed())
	} else {
		// Try to create config file
		if err := viper.SafeWriteConfig(); err != nil {
			klog.Fatalln(err)
		}

		// Read config again
		if err := viper.ReadInConfig(); err != nil {
			klog.Fatalln(err)
		}

		klog.Infoln("Created new config file:", viper.ConfigFileUsed())
	}

}
