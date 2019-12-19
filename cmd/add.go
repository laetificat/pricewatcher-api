package cmd

import (
	"fmt"
	"strings"

	"github.com/laetificat/slogger/pkg/slogger"

	"github.com/laetificat/pricewatcher/internal/helper"
	"github.com/laetificat/pricewatcher/internal/watcher"
	"github.com/spf13/cobra"
)

var (
	domain string
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new price watcher",
		Long:  `Add a new price watcher to keep an eye on a price.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				if err := addDomain(args[0], domain); err != nil {
					slogger.Fatal(err.Error())
				}
			} else {
				_ = cmd.Help()
			}
		},
	}
)

func registerAddCmd() {
	addCmd.PersistentFlags().StringVar(&domain, "domain", "", "define the domain, for example: bol.com, ebay.nl, coolblue.nl, etc")

	rootCmd.AddCommand(addCmd)
}

func addDomain(url, domain string) error {
	if domain == "" {
		var err error
		domain, err = helper.GuessDomain(url)
		if err != nil {
			if strings.EqualFold(helper.NoSupportedDomainFoundErrorMessage, err.Error()) {
				slogger.Info(fmt.Sprintf("Domain '%s' is not supported", domain))
			}
		}
	}

	if helper.IsSupported(domain) {
		return watcher.Add(domain, url)
	}

	slogger.Info(fmt.Sprintf("Domain '%s' is not supported", domain))

	return nil
}
