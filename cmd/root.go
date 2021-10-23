package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "cereal",
		Short: "A serial monitor for environmental data",
		Long:  `A serial monitor for environmental data created by microcontrollers`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cereal.yaml)")

}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		fmt.Printf("Attempting to use config file: %s\n", cfgFile)
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in current directory with name ".cereal" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cereal")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
