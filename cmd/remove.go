package cmd

import (
	"log"
	"strconv"

	"github.com/laetificat/pricewatcher/internal/watcher"
	"github.com/spf13/cobra"
)

var removeAll bool

var (
	removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove a price watcher",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				if err := removeWatcher(args[0]); err != nil {
					log.Panic(err)
				}
			} else {
				if removeAll {
					if err := watcher.RemoveAll(); err != nil {
						log.Panic(err)
					}
				} else {
					_ = cmd.Help()
				}
			}
		},
	}
)

func registerRemoveCmd() {
	removeCmd.PersistentFlags().BoolVarP(&removeAll, "all", "a", false, "remove all the watchers")

	rootCmd.AddCommand(removeCmd)
}

func removeWatcher(id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return watcher.Remove(idInt)
}
