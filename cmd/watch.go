package cmd

import (
	"time"

	"github.com/spf13/viper"

	"github.com/laetificat/pricewatcher-api/internal/log"
	"github.com/laetificat/pricewatcher-api/internal/watcher"
	"github.com/spf13/cobra"
)

var (
	watchCmd = &cobra.Command{
		Use:   "watch",
		Short: "Run all the watchers",
		Run: func(cmd *cobra.Command, args []string) {
			checkWatchers()
			ticker := time.NewTicker(viper.GetDuration("watcher.timeout") * time.Minute)
			for range ticker.C {
				checkWatchers()
			}
		},
	}
)

func registerWatchCmd() {
	watchCmd.PersistentFlags().DurationP(
		"timeout",
		"t",
		10*time.Minute,
		"the amount of minutes to wait before checking",
	)

	if err := viper.BindPFlag("watcher.timeout", watchCmd.PersistentFlags().Lookup("timeout")); err != nil {
		log.Fatal(err.Error())
	}

	rootCmd.AddCommand(watchCmd)
}

func checkWatchers() {
	log.Debug("Checking if queues need to be filled...")
	err := watcher.RunAll()
	if err != nil {
		log.Panic(err)
	}
}
