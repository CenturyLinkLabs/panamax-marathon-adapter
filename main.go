package main

import (
	"fmt"
	"github.com/centurylinklabs/panamax-marathon-adapter/api"
	"github.com/centurylinklabs/panamax-marathon-adapter/marathon"
	"os"
	"flag"
)

func main() {
	//set up command line flag
	var cmd_endpoint string
	flag.StringVar(&cmd_endpoint, "endpoint", "", "Marathon Endpoint")

	endpoint := os.Getenv("MARATHON_ENDPOINT");

	if endpoint == "" {
		flag.Parse();
		endpoint = cmd_endpoint;
	}

	if  endpoint == "" {
		fmt.Println("Error: Invalid endpoint url. Set env. var. 'MARATHON_ENDPOINT' correctly. ")
		os.Exit(1)
	}

	server := api.NewServer(marathon.NewMarathonAdapter(endpoint))
	server.Start()
}
