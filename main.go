package main

import (
	"fmt"
	"os"
	"github.com/centurylinklabs/panamax-marathon-adapter/api"
	"github.com/centurylinklabs/panamax-marathon-adapter/marathon"
)

func main() {
	var endpoint = ""
	if endpoint = os.Getenv("MARATHON_ENDPOINT"); endpoint == "" {
		fmt.Println("Error: Invalid endpoint url. Set env. var. 'MARATHON_ENDPOINT' correctly. ")
		os.Exit(1)
	}
	marathonAdapter := marathon.NewMarathonAdapter(endpoint)
	api.ListenAndServe(marathonAdapter)
}
