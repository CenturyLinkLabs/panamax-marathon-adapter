package marathon

import (
	"time"

	"github.com/CenturyLinkLabs/gomarathon"
	"github.com/CenturyLinkLabs/panamax-marathon-adapter/api"
)

const (
	DEPLOY = iota
	OK
	FAIL
	TIMEOUT
)

const (
	DEPLOY_TIMEOUT = time.Minute * 10
)

type stateFn func(*deployment, *context) stateFn

type context struct {
	values map[string]map[string]string
}

func (c *context) AddKey(key string, values map[string]string) {
	c.values[key] = values
}

func NewContext() context {
	var ctx context
	ctx.values = make(map[string]map[string]string)

	return ctx
}

type status struct {
	code    int
	message string
}

type deployment struct {
	name          string
	status        status
	reqs          map[string]string
	startingState stateFn
	client        gomarathonClientAbstractor
	application   *gomarathon.Application
	submitted     time.Time
}

func createDeployment(service *api.Service, client gomarathonClientAbstractor) deployment {
	var converter = new(MarathonConverter)
	var deployment deployment

	var reqs = make(map[string]string)
	links := service.Links
	for i := range links {
		reqs[links[i].Name] = links[i].Alias
	}

	deployment.name = service.Name
	deployment.reqs = reqs
	deployment.client = client
	deployment.startingState = requirementState
	deployment.status = status{code: DEPLOY}
	deployment.application = converter.convertToApp(service)

	return deployment
}

type deploymentGroup struct {
	id          string
	deployments []deployment
}

func (g *deploymentGroup) Done() bool {
	completed := true
	for _, deployment := range g.deployments {
		completed = completed && (deployment.status.code == OK)
	}

	return completed
}

func (g *deploymentGroup) Failed() bool {
	failed := false

	for _, deployment := range g.deployments {
		failed = failed || (deployment.status.code == FAIL) || (deployment.status.code == TIMEOUT)
	}

	return failed
}
