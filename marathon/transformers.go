package marathon

import (
	"strings"

	"github.com/centurylinklabs/panamax-marathon-adapter/api"
	"github.com/jbdalido/gomarathon"
)

// Split the service string into 2 parts part[0] is group part[1] is service
func splitServiceId(serviceId string, del string) (string, string) {
	var group, service string

	parts := strings.Split(serviceId, del)
	if len(parts) == 2 {
		group = parts[0]
		service = parts[1]
	} else {
		service = parts[0]
	}
	return group, service
}

type PanamaxServiceConverter interface {
	convertToServices([]*gomarathon.Application) []*api.Service
	convertToService(*gomarathon.Application) *api.Service
	convertToApps([]*api.Service) []*gomarathon.Application
	convertToApp(*api.Service) *gomarathon.Application
}

type MarathonConverter struct {
}

func (c *MarathonConverter) convertToServices(apps []*gomarathon.Application) []*api.Service {
	services := make([]*api.Service, len(apps))

	for i := range apps {
		services[i] = c.convertToService(apps[i])
	}

	return services
}

func (c *MarathonConverter) convertToService(app *gomarathon.Application) *api.Service {
	service := new(api.Service)

	service.ActualState = api.StartedStatus
	service.Id = (strings.Replace(app.ID, "/", ".", -1))[1:]
	service.Name = app.ID

	return service
}

func (c *MarathonConverter) convertToApps(services []*api.Service) []*gomarathon.Application {
	apps := make([]*gomarathon.Application, len(services))
	for i := range services {
		apps[i] = c.convertToApp(services[i])
	}

	return apps
}

func (c *MarathonConverter) convertToApp(service *api.Service) *gomarathon.Application {
	app := new(gomarathon.Application)

	// set count to 1 for services with no deployment count specifier
	var count int = 1;
	if service.Deployment.Count > 0 {
		count = service.Deployment.Count
	}

	app.ID = strings.ToLower(service.Name)
	app.Cmd = service.Command
	app.CPUs = 0.5
	app.Env = buildEnvMap(service.Environment)
	app.Instances = count
	app.Mem = 1024
	app.Container = buildDockerContainer(service)

	return app
}

func buildEnvMap(env []*api.Environment) map[string]string {
	envs := make(map[string]string)
	for i := range env {
		envs[env[i].Variable] = env[i].Value
	}

	return envs
}

func buildDockerContainer(service *api.Service) *gomarathon.Container {
	container := new(gomarathon.Container)
	container.Type = "DOCKER"

	docker := new(gomarathon.Docker)
	docker.Image = service.Source
	docker.Network = "BRIDGE"
	docker.PortMappings = buildPortMappings(service.Ports)

	container.Docker = docker

	return container
}

func buildPortMappings(ports []*api.Port) []*gomarathon.PortMapping {
	mappings := make([]*gomarathon.PortMapping, len(ports))

	for i := range ports {
		mapping := new(gomarathon.PortMapping)
		mapping.ContainerPort = ports[i].ContainerPort
		mapping.HostPort = 0

		proto := "tcp"

		if ports[i].Protocol != "" {
			proto = strings.ToLower(ports[i].Protocol)
		}

		mapping.Protocol = proto

		mappings[i] = mapping
	}

	return mappings
}
