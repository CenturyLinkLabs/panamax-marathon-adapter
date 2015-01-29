package marathon

import (
	"testing"

	"github.com/CenturyLinkLabs/gomarathon"
	"github.com/CenturyLinkLabs/panamax-marathon-adapter/api"
	"github.com/stretchr/testify/assert"
)

func TestPreConditionEmpty(t *testing.T) {
	var svc = api.Service{Name: "Foo", Command: "echo"}
	var ctx = NewContext()

	deployment := createDeployment(&svc, nil)
	deployment.name = "TestEmpty"

	requirementState(&deployment, &ctx)
	assert.Equal(t, DEPLOY, deployment.status.code)
}

func TestPreConditionNotFound(t *testing.T) {
	var ctx = NewContext()

	var svc = api.Service{Name: "Foo", Command: "echo"}
	var link = api.Link{Name: "test", Alias: "bar"}
	svc.Links = []*api.Link{&link}
	deployment := createDeployment(&svc, nil)
	deployment.name = "TestLink"

	requirementState(&deployment, &ctx)

	assert.Equal(t, DEPLOY, deployment.status.code)
}

func TestPreConditionFound(t *testing.T) {
	var ctx = NewContext()
	var fooMap = make(map[string]string)
	fooMap["PORT_3306_TCP_PORT"] = "3000"
	ctx.AddKey("foo", fooMap)

	var svc = api.Service{Name: "Bar", Command: "echo"}
	var link = api.Link{Name: "foo", Alias: "bar"}
	svc.Links = []*api.Link{&link}
	deployment := createDeployment(&svc, nil)
	deployment.name = "TestLinked"
	requirementState(&deployment, &ctx)

	assert.Equal(t, DEPLOY, deployment.status.code)
}

func TestPreConditionFoundOnlyFew(t *testing.T) {
	var ctx = NewContext()
	var fooMap = make(map[string]string)
	fooMap["PORT_3306_TCP_PORT"] = "3000"
	ctx.AddKey("foo", fooMap)


	var svc = api.Service{Name: "Bar", Command: "echo"}
	var link = api.Link{Name: "foo", Alias: "bar"}
	var link2 = api.Link{Name: "foo2", Alias: "bar"}
	svc.Links = []*api.Link{&link, &link2}
	deployment := createDeployment(&svc,nil)
	deployment.name = "TestNotLinked"
	requirementState(&deployment, &ctx)

	assert.Equal(t, DEPLOY, deployment.status.code)
}

func TestDockerVar(t *testing.T) {
	var ctx = NewContext()
	var fooMap = make(map[string]string)
	fooMap["PORT_3306_TCP_PORT"] = "3000"
	ctx.AddKey("foo", fooMap)

	var svc = api.Service{Name: "Bar", Command: "echo"}
	var link = api.Link{Name: "foo", Alias: "bar"}
	svc.Links = []*api.Link{&link}
	deployment := createDeployment(&svc,nil)
	deployment.name = "TestLinked"
	requirementState(&deployment, &ctx)

	assert.Equal(t, DEPLOY, deployment.status.code)
	assert.Equal(t, "3000", deployment.application.Env["BAR_PORT_3306_TCP_PORT"])
}

func TestCreateDockerMapping(t *testing.T) {
	var port = gomarathon.PortMapping{ServicePort: 1111, ContainerPort: 5555, Protocol: "UDP"}

	mappings := []*gomarathon.PortMapping{&port}

	docker := createDockerMapping("141.10.10.141", mappings)
	assert.Equal(t, docker["PORT_5555_UDP"], "UDP://141.10.10.141:5555")
	assert.Equal(t, docker["PORT_5555_UDP_PROTO"], "UDP")
	assert.Equal(t, docker["PORT_5555_UDP_ADDR"], "141.10.10.141")
	assert.Equal(t, docker["PORT_5555_UDP_PORT"], "1111")
}

func TestBuildDeployment(t *testing.T) {
	var svc = api.Service{Name: "Foo", Command: "echo"}
	var ctx = NewContext()
	client := new(mockClient)
	resp := new(gomarathon.Response)
	deployment := createDeployment(&svc, client)
	deployment.name = "TestEmpty"

	client.On("CreateApp", deployment.application).Return(resp)

	deploymentState(&deployment, &ctx)

	assert.Equal(t, DEPLOY, deployment.status.code)
	client.AssertExpectations(t)
}

func TestPostAction(t *testing.T) {
	client := new(mockClient)
	resp := new(gomarathon.Response)
	task := new(gomarathon.Task)
	var svc = api.Service{Name: "Foo", Command: "echo"}
	var ctx = NewContext()

	deployment := createDeployment(&svc, client)
	task.Host = "1.2.3.4"
	resp.App = deployment.application
	resp.Tasks = []*gomarathon.Task{task}
	deployment.name = "TestEmpty"

	client.On("GetAppTasks", deployment.application.ID).Return(resp)
	client.On("GetApp", deployment.application.ID).Return(resp)
	postActionState(&deployment, &ctx)

	assert.Equal(t, OK, deployment.status.code)
	client.AssertExpectations(t)
}
