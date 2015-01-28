package marathon

import (
	"github.com/CenturyLinkLabs/panamax-marathon-adapter/api"
	"github.com/jbdalido/gomarathon"
	"github.com/stretchr/testify/assert"
	"testing"
)

var mConverter = new(MarathonConverter)

func TestConvertToServices(t *testing.T) {

	app1 := gomarathon.Application{ID: "foo"}
	app2 := gomarathon.Application{ID: "bar"}
	services := mConverter.convertToServices([]*gomarathon.Application{&app1, &app2})

	assert.Equal(t, 2, len(services))
}

func TestConvertToService(t *testing.T) {

	application := gomarathon.Application{ID: "/foo"}

	service := mConverter.convertToService(&application)

	assert.Equal(t, "/foo", service.Name)
	assert.Equal(t, "foo", service.Id)
}

func TestConvertToApps(t *testing.T) {

	service1 := api.Service{Name: "foo"}
	service2 := api.Service{Name: "bar"}

	apps := mConverter.convertToApps([]*api.Service{&service1, &service2})

	assert.Equal(t, 2, len(apps))
}

func TestConvertToApp(t *testing.T) {

	service := api.Service{Name: "FOO", Command: "echo", Deployment: api.Deployment{Count: 5}}
	app := mConverter.convertToApp(&service)

  service1 := api.Service{Name: "FOO", Command: "echo"}
	app1 := mConverter.convertToApp(&service1)

	assert.Equal(t, "foo", app.ID)
	assert.Equal(t, "echo", app.Cmd)
	assert.Equal(t, 0.5, app.CPUs)
	assert.Equal(t, 5, app.Instances)
	assert.Equal(t, 1, app1.Instances)

}

func TestBuildEnvMap(t *testing.T) {

	env := api.Environment{Variable: "VARIABLE", Value: "VALUE"}

	envs := buildEnvMap([]*api.Environment{&env})

	assert.Equal(t, envs["VARIABLE"], "VALUE")
}

func TestBuildDockerContainer(t *testing.T) {
	ports := []*api.Port{&api.Port{ContainerPort: 3000}}

	service := api.Service{Name: "foo", Source: "centurylink/panamax", Ports: ports}

	container := buildDockerContainer(&service)
	portMappings := buildPortMappings(service.Ports)

	assert.Equal(t, container.Docker.PortMappings, portMappings)
}

func TestBuildPortMappings(t *testing.T) {

	port1 := api.Port{ContainerPort: 3000}
	port2 := api.Port{ContainerPort: 3001}
	port3 := api.Port{ContainerPort: 3002, Protocol: "UDP"}

	mappings := buildPortMappings([]*api.Port{&port1, &port2, &port3})

	assert.Equal(t, 3, len(mappings))
	assert.Equal(t, 3000, mappings[0].ContainerPort)
	assert.Equal(t, 0, mappings[0].HostPort)
	assert.Equal(t, "tcp", mappings[0].Protocol)
	assert.Equal(t, "udp", mappings[2].Protocol)
}
