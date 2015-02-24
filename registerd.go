package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	if printVersion {
		fmt.Printf("registerd %s\n", Version)
		os.Exit(0)
	}
}
