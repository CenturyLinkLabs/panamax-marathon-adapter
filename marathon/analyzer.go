package marathon

import (
	"log"
	"time"
	"math/rand"

	"github.com/centurylinklabs/panamax-marathon-adapter/api"
)

func buildRequirements(service api.Service) stateFn {
	var reqs = make(map[string]string)
	links := service.Links
	for i := range links {
		reqs[links[i].Name] = links[i].Alias
	}

	return func(a *app) int {
		if (len(reqs) == 0) {
			return OK
		} else {
			log.Printf("links non zero")
			return OK
		}
	}

}

func buildDeployment(service api.Service) stateFn {
	return func(a *app) int {
		time.Sleep(time.Duration(rand.Intn(100) * 100) * time.Millisecond)
		return OK
	}
}

func buildPostActions(service api.Service) stateFn {
	return func(a *app) int {
		return OK
	}
}

func CreateAppDeployment(service api.Service) app {
	log.Printf("Building Application Deployment")
	var converter = new(MarathonConverter)
	var application app
	application.name = "APP"
	application.preFn = buildRequirements(service)
	application.deployFn = buildDeployment(service)
	application.postFn = buildPostActions(service)
	application.application = converter.convertToApp(&service)

	return application
}

