package main

import (
	"github.com/laetificat/pricewatcher-api/cmd"
	"github.com/laetificat/pricewatcher-api/internal/log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Panic(err)
	}
}
