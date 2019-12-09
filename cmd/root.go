package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitmoji",
	Short: "Gitmoji helper written in Go.",
	Long:  `gogitmoji helps you write git commit messages containing gitmoji!`,
	Run: func(cmd *cobra.Command, args []string) {
		commit()
	},
}

// AddCommand adds new, top-level command.
func AddCommand(command *cobra.Command) {
	rootCmd.AddCommand(command)
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gitmoji/config.yaml)")

	setHelpEmoji()
}

// Set the help emoji
func setHelpEmoji() {
	rootCmd.InitDefaultHelpCmd()
	commands := rootCmd.Commands()

	for i := 0; i < len(commands); i++ {
		command := commands[i]

		if command.Name() == "help" {
			command.Short = "📗  " + command.Short
			break
		}
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(path.Join(home, ".gitmoji"))
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("gitmoji")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// missing config file is ok
		default:
			log.Fatalf("Error reading '%s': %v", viper.ConfigFileUsed(), err)
		}
	}
}
