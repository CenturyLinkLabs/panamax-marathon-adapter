package marathon

import (
	"log"
	"fmt"
	"time"
	"strings"

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

func requirementState(deployment *deployment, ctx *context) stateFn {
	if (len(deployment.reqs) == 0) {
		return deploymentState
	} else {
		found := true
		for k, _ := range deployment.reqs {
			if (ctx.values[k] == nil) {
				found = false
			}
		}
		if (!found) {
			return requirementState
		} else {
			dockers := loadDockerVars(ctx, deployment.reqs)
			for key, value := range dockers {
				deployment.application.Env[key] = value
			}
			return deploymentState
		}
	}

}

func deploymentState(deployment *deployment, ctx *context) stateFn {
	log.Printf("Starting Deployment: %s", deployment.application.ID)

	_, err := deployment.client.CreateApp(deployment.application)
	time.Sleep(2000 * time.Millisecond)
	if err != nil {
		deployment.status.code = FAIL
		deployment.status.message = fmt.Sprintf("%s",err)
		return nil
	}
	return postActionState

}

func createDockerMapping(host string, mappings []*gomarathon.PortMapping) map[string]string {
	var docker = make(map[string]string)

	for i := range(mappings) {
		docker[fmt.Sprintf("PORT_%d_TCP_ADDR", mappings[i].ContainerPort)] = host
		docker[fmt.Sprintf("PORT_%d_TCP_PORT", mappings[i].ContainerPort)] = fmt.Sprintf("%d",mappings[i].ServicePort)
	}
	return docker
}

func postActionState(deployment *deployment, ctx *context) stateFn {
	application := deployment.application
	_, name := splitServiceId(application.ID[1:], "/")

	res, _ := deployment.client.GetAppTasks(application.ID)
	if len(res.Tasks) != 0 {
		host := res.Tasks[0].Host
		appRes, err := deployment.client.GetApp(application.ID)
		if err != nil {
			deployment.status.code = FAIL
			deployment.status.message = fmt.Sprintf("%s",err)
			return nil
		}
		mappings := createDockerMapping(host, appRes.App.Container.Docker.PortMappings)
		if len(mappings) > 0 {
			ctx.values[name] = mappings
		}

		deployment.status.code = OK
		deployment.status.message = fmt.Sprintf("Successful deployment: %s", deployment.application.ID)
		log.Printf("Successful deployment: %s", deployment.application.ID)

		return nil
	}
	return postActionState

}

