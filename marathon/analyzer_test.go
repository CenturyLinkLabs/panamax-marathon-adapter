package marathon

import (
	"testing"

	"github.com/jbdalido/gomarathon"
	"github.com/centurylinklabs/panamax-marathon-adapter/api"
	"github.com/stretchr/testify/assert"
)

func TestPreConditionEmpty(t *testing.T) {
	var svc = api.Service{Name: "Foo", Command: "echo"}
	var ctx = NewContext()

	app := CreateAppDeployment(&svc, nil)
	app.name = "TestEmpty"

	status := app.preFn(app.application, &ctx)
	assert.Equal(t, OK, status)
}

func TestPreConditionNotFound(t *testing.T) {
	var ctx = NewContext()

	var svc = api.Service{Name: "Foo", Command: "echo"}
	var link = api.Link{Name: "test", Alias: "bar"}
	svc.Links = []*api.Link{&link}
	app := CreateAppDeployment(&svc, nil)
	app.name = "TestLink"
	status := app.preFn(app.application, &ctx)

	var svc1 = api.Service{Name: "Food", Command: "echo"}
	app1 := CreateAppDeployment(&svc1, nil)
	app1.name = "TestEmpty"
	status1 := app1.preFn(app1.application, &ctx)

	assert.Equal(t, WAIT, status)
	assert.Equal(t, OK, status1)
}

func TestPreConditionFound(t *testing.T) {
	var ctx = NewContext()
	var fooMap = make(map[string]string)
	fooMap["PORT_3306_TCP_PORT"] = "3000"
	ctx.AddKey("foo", fooMap)

	var svc = api.Service{Name: "Bar", Command: "echo"}
	var link = api.Link{Name: "foo", Alias: "bar"}
	svc.Links = []*api.Link{&link}
	app := CreateAppDeployment(&svc, nil)
	app.name = "TestLinked"
	status := app.preFn(app.application, &ctx)

	assert.Equal(t, OK, status)
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
	app := CreateAppDeployment(&svc,nil)
	app.name = "TestNotLinked"
	status := app.preFn(app.application, &ctx)

	assert.Equal(t, WAIT, status)
}

func TestDockerVar(t *testing.T) {
	var ctx = NewContext()
	var fooMap = make(map[string]string)
	fooMap["PORT_3306_TCP_PORT"] = "3000"
	ctx.AddKey("foo", fooMap)

	var svc = api.Service{Name: "Bar", Command: "echo"}
	var link = api.Link{Name: "foo", Alias: "bar"}
	svc.Links = []*api.Link{&link}
	app := CreateAppDeployment(&svc,nil)
	app.name = "TestLinked"
	status := app.preFn(app.application, &ctx)

	assert.Equal(t, OK, status)
	assert.Equal(t, "3000", app.application.Env["BAR_PORT_3306_TCP_PORT"])
}



func TestBuildDeployment(t *testing.T) {
	var svc = api.Service{Name: "Foo", Command: "echo"}
	var ctx = NewContext()
	client := new(mockClient)
	resp := new(gomarathon.Response)
	app := CreateAppDeployment(&svc, client)
	app.name = "TestEmpty"

	client.On("CreateApp", app.application).Return(resp)

	status := app.deployFn(app.application, &ctx)

	assert.Equal(t, OK, status)
	client.AssertExpectations(t)
}

func TestPostAction(t *testing.T) {
	client := new(mockClient)
	resp := new(gomarathon.Response)
	var svc = api.Service{Name: "Foo", Command: "echo"}
	var ctx = NewContext()

	app := CreateAppDeployment(&svc, client)
	resp.App = app.application
	app.name = "TestEmpty"

	client.On("GetAppTasks", app.application.ID).Return(resp)
	client.On("GetApp", app.application.ID).Return(resp)
	status := app.postFn(app.application, &ctx)

	assert.Equal(t, OK, status)
	client.AssertExpectations(t)
}
