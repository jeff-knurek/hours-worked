package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "hours-worked",
	Short: "Track minutes that the user is active",
	Long: `Track minutes that the user is active
Usage of track:
hours-worked [COMMAND] [OPTIONS]

Hoours worked increments the minutes which the user is active in the UI.
Configuration options can be set in the config file, and are created by default on the first run.

It's best to setup the app to auto-start with the "track" command when user first logs in.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hours-worked/config.toml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in home directory with name ".hours-worked"
		viper.AddConfigPath(home)
		viper.SetConfigName(".hours-worked/config")
		viper.SetConfigType("toml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		// else, create a new config file
		viper.Set("tracking_file", filepath.Join(home, ".hours-worked", "tracked.json"))
		if err := viper.SafeWriteConfig(); err != nil {
			fmt.Println("problem with creating a new config file", err)
		}
	}
}
