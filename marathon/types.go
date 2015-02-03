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


// A state function type used to define the states within
// deployment workflow.
type stateFn func(*deployment, *context) stateFn

// The context is a helper construct used to define a
// map of maps. The context is used by the deployment
// workflow to share data between steps.
type context struct {
	values map[string]map[string]string
}

// The context AddKey function is a helper to insert maps into the
// context under a provided key.
func (c *context) AddKey(key string, values map[string]string) {
	c.values[key] = values
}

// Creates a new context with empty values.
func NewContext() context {
	var ctx context
	ctx.values = make(map[string]map[string]string)

	return ctx
}

// A deployment status structure used by the deployment
// process to determine success or failure of a task.
type status struct {
	code    int
	message string
}

// The deployment structure defines a deployment task for
// the workflow. A deployment task requires a name and a
// startingState.
type deployment struct {
	name          string
	status        status
	reqs          map[string]string
	startingState stateFn
	client        gomarathonClientAbstractor
	application   *gomarathon.Application
	submitted     time.Time
}

// A constructor method for a deployment task.
//
// Given a service and a gomarathon client creates a complete deployment task structure.
// The service is investigated to build a requirements list and the starting state is
// set to the requirementState.
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

// The deploymentGroup gathers and manages a collection of deployment
// tasks.
type deploymentGroup struct {
	id          string
	deployments []deployment
}

// Determines if a deployment group is done.
//
// A group is done only when all of the deployment tasks have been completed.
func (g *deploymentGroup) Done() bool {
	completed := true
	for _, deployment := range g.deployments {
		completed = completed && (deployment.status.code == OK)
	}

	return completed
}

// Determines if deployment group has failed.
//
// A group has failed if any of the deployment tasks have failed.
func (g *deploymentGroup) Failed() bool {
	failed := false

	for _, deployment := range g.deployments {
		failed = failed || (deployment.status.code == FAIL) || (deployment.status.code == TIMEOUT)
	}

	return failed
}
