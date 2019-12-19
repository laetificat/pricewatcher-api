package main

import (
	"github.com/laetificat/pricewatcher/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
