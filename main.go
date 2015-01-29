package main // import "github.com/CenturyLinkLabs/panamax-marathon-adapter"

import (
	"flag"
	"fmt"
	"github.com/CenturyLinkLabs/panamax-marathon-adapter/api"
	"github.com/CenturyLinkLabs/panamax-marathon-adapter/marathon"
	"os"
)

func main() {
	//set up command line flag
	var cmd_endpoint string
	flag.StringVar(&cmd_endpoint, "endpoint", "", "Marathon Endpoint")

	endpoint := os.Getenv("MARATHON_ENDPOINT")

	if endpoint == "" {
		flag.Parse()
		endpoint = cmd_endpoint
	}

	if endpoint == "" {
		fmt.Println("Error: Invalid endpoint url. Set env. var. 'MARATHON_ENDPOINT' correctly. ")
		os.Exit(1)
	}

	server := api.NewServer(marathon.NewMarathonAdapter(endpoint))
	server.Start()
}
