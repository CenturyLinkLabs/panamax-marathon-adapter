package marathon

import (
	"log"
	"fmt"
	"time"
	"strings"

	"github.com/centurylinklabs/panamax-marathon-adapter/api"
	"github.com/jbdalido/gomarathon"
)

func loadDockerVars(ctx *context, reqs map[string]string) map[string]string {
	var docker = make(map[string]string)

	for k, alias := range reqs {
		for name, value := range ctx.values[k] {
			key := fmt.Sprintf("%s_%s", strings.ToUpper(alias), strings.ToUpper(name))
			docker[key] = value
		}
	}

	return docker

}

func buildRequirements(service *api.Service) stateFn {
	var reqs = make(map[string]string)
	links := service.Links
	for i := range links {
		reqs[links[i].Name] = links[i].Alias
	}

	return func(a *gomarathon.Application, ctx *context) int {
		if (len(reqs) == 0) {
			return OK
		} else {
			found := true
			for k, _ := range reqs {
				if (ctx.values[k] == nil) {
					found = false
				}
			}
			if (!found) {
				return WAIT
			} else {
				dockers := loadDockerVars(ctx, reqs)
				for key, value := range dockers {
					a.Env[key] = value
				}
				return OK
			}
		}
	}

}

func buildDeployment(service *api.Service, client gomarathonClientAbstractor) stateFn {
	return func(a *gomarathon.Application, ctx *context) int {
		log.Printf("Starting Deployment: %s", a.ID)
		_, err := client.CreateApp(a)
		time.Sleep(2000 * time.Millisecond)
		if err != nil {
			return FAIL
		}
		return OK
	}
}

func createDockerMapping(mappings []*gomarathon.PortMapping) map[string]string {
	var docker = make(map[string]string)

	for i := range(mappings) {
		docker[fmt.Sprintf("PORT_%d_TCP_ADDR", mappings[i].ContainerPort)] = "10.141.141.10"
		docker[fmt.Sprintf("PORT_%d_TCP_PORT", mappings[i].ContainerPort)] = fmt.Sprintf("%d",mappings[i].ServicePort)
	}
	return docker
}

func buildPostActions(service *api.Service, client gomarathonClientAbstractor) stateFn {
	_, name := splitServiceId(service.Name[1:], "/")

	return func(a *gomarathon.Application, ctx *context) int {
		log.Printf("Post Action for: %s", name)
		res, _ := client.GetAppTasks(a.ID)
		if len(res.Tasks) == 0 {
			appRes, err := client.GetApp(a.ID)
			if err != nil {
				return FAIL
			}
			mappings := createDockerMapping(appRes.App.Container.Docker.PortMappings)
			if len(mappings) > 0 {
				ctx.values[name] = mappings
			}

			return OK
		}
		return WAIT
	}
}

func CreateAppDeployment(service *api.Service, client gomarathonClientAbstractor) app {
	log.Printf("Building Application Deployment %s", service.Name )
	var converter = new(MarathonConverter)
	var application app
	application.name = service.Name
	application.preFn = buildRequirements(service)
	application.deployFn = buildDeployment(service, client)
	application.postFn = buildPostActions(service, client)
	application.application = converter.convertToApp(service)

	return application
}

