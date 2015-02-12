package marathon

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/CenturyLinkLabs/gomarathon"
)

func loadDockerVars(ctx *context, reqs map[string]string) map[string]string {
	var docker = make(map[string]string)

	for k, alias := range reqs {
		for name, value := range ctx.values[strings.ToLower(k)] {
			key := fmt.Sprintf("%s_%s", strings.ToUpper(alias), strings.ToUpper(name))
			docker[key] = value
		}
	}

	return docker
}

// The requirement gathering state.
//
// Validates the requirements listed in the deployment are found in the
// shared context. If the requirements are not found then stay in the
// requirement gathering state. If the requirements are found then load
// the docker variables and return the deployment state.
func requirementState(deployment *deployment, ctx *context) stateFn {
	if len(deployment.reqs) == 0 {
		return deploymentState
	} else {
		found := true
		for k, _ := range deployment.reqs {

			if ctx.values[strings.ToLower(k)] == nil {
				found = false
			}
		}
		if !found {
			return requirementState
		} else {
			log.Printf("Deployment Requirements: %s", deployment.application.ID)
			dockers := loadDockerVars(ctx, deployment.reqs)
			for key, value := range dockers {
				deployment.application.Env[key] = value
			}
			return deploymentState
		}
	}
}

// The deployment state.
//
// Using the client on the deployment task desploy the application.
// If the deployment was successful return the post action state.
// If the deployment failed then set fail code and return nil.
func deploymentState(deployment *deployment, ctx *context) stateFn {
	_, err := deployment.client.CreateApp(deployment.application)
	time.Sleep(2000 * time.Millisecond)
	if err != nil {
		deployment.status.code = FAIL
		deployment.status.message = fmt.Sprintf("%s", err)
		log.Printf("Error Deployment: %s - %s", deployment.application.ID, err)
		return nil
	}
	return postActionState
}

func createDockerMapping(app *gomarathon.Application, name string) map[string]string {
	var min_port = 0
	var host = ""
	var docker = make(map[string]string)

	docker[fmt.Sprintf("NAME")] = fmt.Sprintf("%s", name)

	if len(app.Tasks) > 0 {
		host = app.Tasks[0].Host
		if len(app.Tasks[0].Ports) > 0 {
			min_port = app.Tasks[0].Ports[0]
			mappings := app.Container.Docker.PortMappings
			protocol := strings.ToUpper(mappings[0].Protocol)

			docker[fmt.Sprintf("PORT")] = fmt.Sprintf("%s://%s:%d", protocol, host, min_port)

			for i := range mappings {
				servicePort := app.Tasks[0].Ports[i]
				containerPort := mappings[i].ContainerPort
				protocol := strings.ToUpper(mappings[i].Protocol)

				docker[fmt.Sprintf("PORT_%d_%s", containerPort, protocol)] = fmt.Sprintf("%s://%s:%d", protocol, host, servicePort)
				docker[fmt.Sprintf("PORT_%d_%s_PROTO", containerPort, protocol)] = fmt.Sprintf("%s", protocol)
				docker[fmt.Sprintf("PORT_%d_%s_ADDR", containerPort, protocol)] = host
				docker[fmt.Sprintf("PORT_%d_%s_PORT", containerPort, protocol)] = fmt.Sprintf("%d", servicePort)
			}
		}
	}

	return docker
}

// The post action state.
//
// Find all the network mappings on this application and create the
// required docker environment variables. These variables are set into
// the shared context under this application's identifier.
func postActionState(deployment *deployment, ctx *context) stateFn {
	application := deployment.application
	_, name := splitServiceId(application.ID[1:], "/")

	res, _ := deployment.client.GetApp(application.ID)
	if res.App.TasksRunning > 0 {
		log.Printf("Post Deployment: %s", deployment.application.ID)
		appRes, err := deployment.client.GetApp(application.ID)
		if err != nil {
			deployment.status.code = FAIL
			deployment.status.message = fmt.Sprintf("%s", err)
			return nil
		}

		mappings := createDockerMapping(appRes.App, name)
		if len(mappings) > 0 {
			log.Printf("Adding mappings : %s", strings.ToLower(name))
			ctx.values[strings.ToLower(name)] = mappings
		}

		deployment.status.code = OK
		deployment.status.message = fmt.Sprintf("Successful deployment: %s", deployment.application.ID)
		log.Printf("Successful deployment: %s", deployment.application.ID)

		return nil
	}
	time.Sleep(2000 * time.Millisecond)
	return postActionState
}
