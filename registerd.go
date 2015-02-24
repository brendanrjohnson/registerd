package main

import (
	"flag"
	"fmt"
	"os"

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
}
