package main

import (
	"fmt"
	"github.com/centurylinklabs/panamax-marathon-adapter/api"
	"github.com/centurylinklabs/panamax-marathon-adapter/marathon"
	"os"
)

func main() {
	var endpoint = ""
	if endpoint = os.Getenv("MARATHON_ENDPOINT"); endpoint == "" {
		fmt.Println("Error: Invalid endpoint url. Set env. var. 'MARATHON_ENDPOINT' correctly. ")
		os.Exit(1)
	}

	server := api.NewServer(marathon.NewMarathonAdapter(endpoint))
	server.Start()
}
