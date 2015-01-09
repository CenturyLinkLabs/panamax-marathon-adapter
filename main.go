package main

import (
	"github.com/CenturyLinkLabs/panamax-marathon-adapter/api"
	"github.com/CenturyLinkLabs/panamax-marathon-adapter/marathon"
)

func main() {
	marathonAdapter := marathon.NewMarathonAdapter("http://10.141.141.10:8080")
	api.ListenAndServe(marathonAdapter)
}
