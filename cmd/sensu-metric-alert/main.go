package main

import (
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stderr)

	if len(os.Args) == 1 {
		println("usage: sensu-metric-alert -m <metric> --lt|--gt|--ne <crit value>")
		os.Exit(2)
	}
}
