package cmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/laetificat/pricewatcher/internal/queue"
	"github.com/laetificat/pricewatcher/internal/watcher"
	"github.com/laetificat/pricewatcher/internal/web/api"
	"github.com/laetificat/pricewatcher/internal/web/middleware"
	"github.com/laetificat/slogger/pkg/slogger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	webserverCmd = &cobra.Command{
		Use:   "webserver",
		Short: "Start the webserver",
		Run: func(cmd *cobra.Command, args []string) {
			registerQueues()
			runWebserver()
		},
	}
)

func registerWebserverCmd() {
	webserverCmd.PersistentFlags().StringP(
		"address",
		"a",
		"http://localhost:8080",
		"the ip address to bind the webserver on",
	)

	if err := viper.BindPFlag("webserver.address", webserverCmd.PersistentFlags().Lookup("address")); err != nil {
		slogger.Fatal(err.Error())
	}
	viper.SetDefault("webserver.address", "http://localhost:8080")

	rootCmd.AddCommand(webserverCmd)
}

/*
registerQueues registers queues based on the supported domains that are supported
*/
func registerQueues() {
	for _, v := range watcher.SupportedDomains {
		queueName := queue.GetNameForDomain(v)

		slogger.Debug(
			fmt.Sprintf("creating queue '%s'", queueName),
		)

		if err := queue.Create(queueName); err != nil {
			slogger.Fatal(err.Error())
		}
	}
}

/*
runWebserver registers the routes, adds middlewares and starts listening on the given address and port
*/
func runWebserver() {
	slogger.Info("Starting webserver...")

	router := httprouter.New()

	slogger.Debug("Registering routes...")
	api.RegisterHomeHandler(router)
	api.RegisterWatcherHandler(router)
	api.RegisterPriceHandler(router)
	api.RegisterQueueHandler(router)

	routerWithMiddleWare := middleware.NewLogMiddleWare(router)

	slogger.Info(
		fmt.Sprintf("Server started on %s", viper.GetString("webserver.address")),
	)

	address := strings.ReplaceAll(
		strings.ReplaceAll(
			viper.GetString("webserver.address"),
			"https://",
			"",
		),
		"http://",
		"",
	)

	slogger.Fatal(
		http.ListenAndServe(address, routerWithMiddleWare).Error(),
	)
}
