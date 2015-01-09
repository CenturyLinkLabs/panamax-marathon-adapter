package marathon


import (
	"github.com/CenturyLinkLabs/panamax-marathon-adapter/api"
	"github.com/jbdalido/gomarathon"
)

func convertApps(apps []*gomarathon.Application) ([]*api.Service) {
	services := make([]*api.Service, len(apps))

	for i := range apps {
		services[i] = convertApp(apps[i])
	}

	return services
}

func convertApp(app *gomarathon.Application) (*api.Service) {
	service := new(api.Service)

	service.CurrentState = api.StartedStatus
	service.Id = app.ID

	return service
}

func convertServices(services []*api.Service) ([]*gomarathon.Application) {
	apps := make([]*gomarathon.Application, len(services))
	for i := range services {
		apps[i] = convertService(services[i])
	}

	return apps
}

func convertService(service *api.Service) (*gomarathon.Application) {
	app := new(gomarathon.Application)

	app.ID = service.Name
	app.Cmd = service.Command
	app.CPUs = 0.5
	app.Env = buildEnvMap(service.Environment)
	app.Instances = 1
	app.Mem = 1024
	app.Container = buildDockerContainer(service)

	return app
}

func buildEnvMap(env []*api.Environment) (map[string]string) {
	envs := make(map[string]string)
	for i := range env {
		envs[env[i].Variable] = env[i].Value
	}

	return envs
}

func buildDockerContainer(service *api.Service) (*gomarathon.Container) {
	container := new(gomarathon.Container)
	container.Type = "DOCKER"

	docker := new(gomarathon.Docker)
	docker.Image = service.Source
	docker.Network = "BRIDGE"
	docker.PortMappings = buildPortMappings(service.Ports)

	container.Docker = docker

	return container
}

func buildPortMappings(ports []*api.Port) ([]*gomarathon.PortMapping) {
	mappings := make([]*gomarathon.PortMapping, len(ports))

	for i := range ports {
		mapping := new(gomarathon.PortMapping)
		mapping.ContainerPort = ports[i].ContainerPort
		mapping.HostPort = 0
		mapping.Protocol = "tcp"

		mappings[i] = mapping
	}

	return mappings
}
