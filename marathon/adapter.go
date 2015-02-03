// Package marathon is the Marathon implementation for a Panamax Remote Adapter.
package marathon // import "github.com/CenturyLinkLabs/panamax-marathon-adapter/marathon"

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/CenturyLinkLabs/gomarathon"
	"github.com/CenturyLinkLabs/panamax-marathon-adapter/api"
)

// Creates a client connection to Marathon on the provided endpoint.
func newClient(endpoint string) *gomarathon.Client {
	url := endpoint
	if endpoint != "" {
		url = endpoint
	}
	log.Printf("Marathon Endpoint: %s", url)
	c, err := gomarathon.NewClient(url, nil)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func sanitizeServiceName(name string) string {
	name = strings.Replace(name, " ", "", -1)
	name = strings.Replace(name, "-", "", -1)
	name = strings.Replace(name, "_", "", -1)
	name = strings.Replace(name, ",", "", -1)
	return name
}

func sanitizeMarathonAppURL(id string) string {
	group, service := splitServiceId(id, ".")
	return fmt.Sprintf("%s/%s", strings.ToLower(group), strings.ToLower(service))
}

type gomarathonClientAbstractor interface {
	ListApps() (*gomarathon.Response, error)
	GetApp(string) (*gomarathon.Response, error)
	GetAppTasks(string) (*gomarathon.Response, error)
	CreateApp(*gomarathon.Application) (*gomarathon.Response, error)
	CreateGroup(*gomarathon.Group) (*gomarathon.Response, error)
	GetGroup(string) (*gomarathon.Response, error)
	DeleteApp(string) (*gomarathon.Response, error)
	DeleteGroup(string) (*gomarathon.Response, error)
}

type marathonAdapter struct {
	client   gomarathonClientAbstractor
	conv     PanamaxServiceConverter
	deployer Deployer
}

// Create an instance of the marathon adapter.
func NewMarathonAdapter(endpoint string) *marathonAdapter {
	adapter := new(marathonAdapter)
	adapter.client = newClient(endpoint)
	adapter.conv = new(MarathonConverter)
	adapter.deployer = newMarathonDeployer()

	return adapter
}

// Implementation of the PanamaxAdapter GetServices interface
func (m *marathonAdapter) GetServices() ([]*api.Service, *api.Error) {
	var apiErr *api.Error

	response, err := m.client.ListApps()
	if err != nil {
		apiErr = api.NewError(http.StatusNotFound, err.Error())
	}
	return m.conv.convertToServices(response.Apps), apiErr
}

// Implementation of the PanamaxAdapter GetService interface
func (m *marathonAdapter) GetService(id string) (*api.Service, *api.Error) {
	var apiErr *api.Error

	response, err := m.client.GetApp(sanitizeMarathonAppURL(id))
	if err != nil {
		apiErr = api.NewError(http.StatusNotFound, err.Error())
	}
	return m.conv.convertToService(response.App), apiErr
}

// Implementation of the PanamaxAdapter CreateServices interface
func (m *marathonAdapter) CreateServices(services []*api.Service) ([]*api.Service, *api.Error) {
	var apiErr *api.Error

	myGroup := m.deployer.BuildDeploymentGroup(services, m.client)
	status := m.deployer.DeployGroup(myGroup, DEPLOY_TIMEOUT)

	switch status.code {
	case FAIL:
		apiErr = api.NewError(http.StatusConflict, "Group deployment failed.")
	case TIMEOUT:
		apiErr = api.NewError(http.StatusInternalServerError, "Group deployment timed out.")
	}

	return services, apiErr
}

// Implementation of the PanamaxAdapter UpdateService interface
func (m *marathonAdapter) UpdateService(s *api.Service) *api.Error {
	return nil
}

// Implementation of the PanamaxAdapter DestroyService interface
func (m *marathonAdapter) DestroyService(id string) *api.Error {
	var apiErr *api.Error
	group, _ := splitServiceId(id, ".")

	_, err := m.client.DeleteApp(sanitizeMarathonAppURL(id))
	if err != nil {
		apiErr = api.NewError(http.StatusNotFound, err.Error())
	}

	m.client.DeleteGroup(group) // Remove group if possible we dont care about error or return.

	return apiErr
}
