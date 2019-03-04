package main

import (
	"flag"
	"log"
)

var (
	serverMode bool
	clientMode bool
	mongoPort  string
	grpcPort   string
)

func main() {

	// for better loggin in crash
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.BoolVar(&serverMode, "s", false, "run as the server. Runs as client by default.")
	flag.StringVar(&mongoPort, "m", "27017", "change local port for Mongo")
	flag.StringVar(&grpcPort, "g", "50051", "change local port for gRPC")
	flag.Parse()

	if serverMode {
		NewServer().Run()
	} else {
		NewClient().Run()
	}
}
