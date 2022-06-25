package commands

import (
	"log"
	"os"
	"path/filepath"

	"github.com/caguiclajmg/tensordock-cli/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	client  *api.Client

	rootCmd = &cobra.Command{
		Use:          "tensordock-cli",
		Short:        "A brief description of your application",
		SilenceUsage: true,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	pflags := rootCmd.PersistentFlags()
	pflags.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tensordock.yml)")
	pflags.String("apiKey", "", "API key")
	pflags.String("apiToken", "", "API token")
	pflags.Bool("debug", false, "Enable debug mode")

	viper.BindPFlag("apiKey", pflags.Lookup("apiKey"))
	viper.BindPFlag("apiToken", pflags.Lookup("apiToken"))
	viper.BindPFlag("debug", pflags.Lookup("debug"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.SetConfigFile(filepath.Join(home, ".tensordock.yml"))
	}

	viper.SetDefault("serviceUrl", "https://console.tensordock.com/api")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("warning: config file %v not found", viper.ConfigFileUsed())
	}

	viper.AutomaticEnv()

	serviceUrl := viper.GetString("serviceUrl")
	apiKey := viper.GetString("apiKey")
	apiToken := viper.GetString("apiToken")
	debug := viper.GetBool("debug")

	client = api.NewClient(serviceUrl, apiKey, apiToken, debug)
}
