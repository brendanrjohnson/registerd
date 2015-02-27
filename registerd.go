package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/brendanrjohnson/registerd/backends"
	"github.com/kelseyhightower/confd/log"
)

func main() {
	flag.Parse()
	if printVersion {
		fmt.Printf("registerd %s\n", Version)
		os.Exit(0)
	}
	if err := initConfig(); err != nil {
		log.Fatal(err.Error())
	}
	log.Notice("Starting registerd")
	storeClient, err := backends.New(backendsConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	templateConfig.StoreClient = storeClient
	if onetime {
		if err := templ
	}
}
