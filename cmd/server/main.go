package main

import (
	"flag"

	"github.com/limanmys/cloud-manager/internal/server"
)

func main() {
	runType := flag.String("type", "client", "run type")
	flag.Parse()
	if *runType == "admin" {
		server.RunAdmin(false)
	} else if *runType == "client" {
		server.RunClient()
	} else if *runType == "admintest" {
		server.RunAdmin(true)
	}
}
