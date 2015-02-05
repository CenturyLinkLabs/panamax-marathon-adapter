package marathon

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/CenturyLinkLabs/panamax-marathon-adapter/api"
)

// Deployment interface for building a group deployment and
// issuing the actual deployment via provided client.
type Deployer interface {
	BuildDeploymentGroup([]*api.Service, gomarathonClientAbstractor) *deploymentGroup
	DeployGroup(*deploymentGroup, time.Duration) status
}

type MarathonDeployer struct {
	uidGenerator func() string
}

func newMarathonDeployer() *MarathonDeployer {
	var deployer MarathonDeployer

	deployer.uidGenerator = func() string {
		num := make([]byte, 4)
		rand.Read(num)

		return fmt.Sprintf("%X", num)
	}

	return &deployer
}

// BuildDeploymentGroup converts a list of api Services into a deployment group.
//
// Generates the unique identifier for the marathon group, investigates the
// services dependencies, and converts the services into a list of gomarathon Applications.
func (m MarathonDeployer) BuildDeploymentGroup(services []*api.Service, client gomarathonClientAbstractor) *deploymentGroup {
	var deployments = make([]deployment, len(services))
	g := m.generateUniqueUID(client)

	dependents := m.findDependencies(services)
	for i := range services {
		if dependents[services[i].Name] != 0 {
			services[i].Deployment.Count = 1
		}

		m.prepareServiceForDeployment(g, services[i])
		deployments[i] = createDeployment(services[i], client)
	}

	myGroup := new(deploymentGroup)
	myGroup.deployments = deployments

	return myGroup
}

// DeployGroup
//
// Manages a group of deployment structures as a single deployment.
// It uses a deployment channel and a timeout channel to determine
// if the overall deployment was successful, failed, or was unable to complete
// within a given duration.
func (m MarathonDeployer) DeployGroup(myGroup *deploymentGroup, timeout time.Duration) status {
	log.Printf("Deploying Group: %s", myGroup.id)

	var ctx = NewContext()

	// set up timeout channel
	timeoutChannel := timeoutChannel(timeout)

	// set up deployment channel
	deploymentChannel := deployGroupChannel(myGroup, &ctx)

	for {
		select {
		case <-timeoutChannel:
			log.Printf("Deployment timed out")
			ctx.signal = TIMEOUT
			return status{code: TIMEOUT}
		case <-deploymentChannel:
			if myGroup.Done() {
				return status{code: OK}
			}

			if myGroup.Failed() {
				log.Printf("Deployment Failed")
				return status{code: FAIL}
			}
		}
	}

}

func (m MarathonDeployer) prepareServiceForDeployment(group string, service *api.Service) {
	var serviceName = sanitizeServiceName(service.Name)

	service.Id = fmt.Sprintf("%s.%s", group, serviceName)
	service.Name = fmt.Sprintf("/%s/%s", group, serviceName)
	service.ActualState = "deploying"
}

func (m MarathonDeployer) findDependencies(services []*api.Service) map[string]int {
	var deps = make(map[string]int)
	for s := range services {
		for l := range services[s].Links {
			deps[services[s].Links[l].Name] = 1
		}
	}

	return deps
}

func (m MarathonDeployer) generateUniqueUID(client gomarathonClientAbstractor) string {
	//generate a random number expressed in hex
	uid := m.uidGenerator()

	if _, err := client.GetGroup(uid); err == nil {
		uid = m.generateUniqueUID(client)
	}

	return uid
}

func timeoutChannel(duration time.Duration) chan bool {
	// make a timeout channel
	timeout := make(chan bool)
	go func() {
		select {
		case <-time.After(duration):
    			timeout <- true
		}
	}()
	return timeout
}

func deployGroupChannel(myGroup *deploymentGroup, ctx *context) chan status {
	deploymentChannel := make(chan status, len(myGroup.deployments))
	for i := 0; i < len(myGroup.deployments); i++ {
		go deploy(deploymentChannel, &myGroup.deployments[i], ctx)
	}

	return deploymentChannel
}

func deploy(done chan status, deployment *deployment, ctx *context) {
	log.Printf("Starting Deployment: %s", deployment.name)

	for state := deployment.startingState; state != nil; {
		if ctx.signal == TIMEOUT {
			state = nil
		} else {
			state = state(deployment, ctx)
		}
	}

	done <- deployment.status
}
