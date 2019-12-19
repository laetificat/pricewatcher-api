package cmd

import (
	"fmt"
	"log"

	"github.com/laetificat/slogger/pkg/slogger"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "pricewatcher",
		Short: "A tool to manage your price watchers",
		Long:  `Pricewatcher is a CLI tool that manages your price watcher API for your horde of watchers.`,
	}
)

// Execute executes the root command.
func Execute() error {
	registerRootCmd()
	registerWatchCmd()
	registerWebserverCmd()
	registerRemoveCmd()
	registerListCmd()
	registerAddCmd()
	return rootCmd.Execute()
}

func registerRootCmd() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config",
		"c",
		"",
		"config file (default is $HOME/.pricewatcher-api.toml)",
	)

	rootCmd.PersistentFlags().StringP(
		"verbose",
		"v",
		"info",
		"the minimum log verbosity level",
	)

	if err := viper.BindPFlag("log.minimum_level", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		slogger.Fatal(err.Error())
	}
	viper.SetDefault("log.minimum_level", "info")

	rootCmd.PersistentFlags().StringP(
		"database",
		"d",
		"watchers.db",
		"the path to the database file",
	)

	if err := viper.BindPFlag("database_file", rootCmd.PersistentFlags().Lookup("database")); err != nil {
		slogger.Fatal(err.Error())
	}
	viper.SetDefault("database_file", "watchers.db")
}

/*
initConfig sets up and configures the application
*/
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".pricewatcher-api")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	viper.Set("version", "0.1.0")

	c := slogger.Config{
		Level: viper.GetString("log.minimum_level"),
		Sentry: slogger.SentryZapOptions{
			Enabled: viper.GetBool("log.sentry.enabled"),
			Dsn:     viper.GetString("log.sentry.dsn"),
		},
		Loggly: slogger.LogglyZapOptions{
			Enabled: viper.GetBool("log.loggly.enabled"),
			Token:   viper.GetString("log.loggly.token"),
		},
	}
	slogger.SetConfig(c)

	slogger.Debug(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
}
