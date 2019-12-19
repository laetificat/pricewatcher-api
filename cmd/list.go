package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/laetificat/pricewatcher/internal/helper"
	"github.com/laetificat/pricewatcher/internal/watcher"
	"github.com/laetificat/slogger/pkg/slogger"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Show a list of watchers or domains",
		Long: `Shows a list of items based on the given argument, available lists are
- domains
- watchers`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {

				switch arg := args[0]; arg {
				case "domains":
					listDomains(os.Stdout)
				case "watchers":
					if err := listWatchers(map[string]string{}, os.Stdout); err != nil {
						slogger.Fatal(err.Error())
					}
				default:
					_ = cmd.Help()
				}
			} else {
				_ = cmd.Help()
			}
		},
	}
)

func registerListCmd() {
	rootCmd.AddCommand(listCmd)
}

func listDomains(writer io.Writer) {
	if _, err := writer.Write([]byte(fmt.Sprintln("Supported domains:"))); err != nil {
		fmt.Println(err)
	}

	for _, domain := range helper.GetSupportedDomains() {
		if _, err := writer.Write([]byte(fmt.Sprintf("- %s\n", domain))); err != nil {
			fmt.Println(err)
		}
	}
}

func listWatchers(filters map[string]string, writer io.Writer) error {
	watcherList, err := watcher.List(filters)
	if err != nil {
		return err
	}

	for _, v := range watcherList {
		if _, err := writer.Write([]byte(fmt.Sprintf("%+v\n", v))); err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
