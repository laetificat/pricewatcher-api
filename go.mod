module github.com/laetificat/pricewatcher

go 1.13

require (
	github.com/fatih/structs v1.1.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/laetificat/slogger v0.1.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.1
	go.etcd.io/bbolt v1.3.3
	golang.org/x/text v0.3.2 // indirect
)

replace github.com/laetificat/slogger => ../slogger
