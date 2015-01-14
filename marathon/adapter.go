package marathon

import (
	"github.com/centurylinklabs/panamax-marathon-adapter/api"
	"github.com/jbdalido/gomarathon"
	"log"
)

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

type gomarathonClientAbstractor interface {
	ListApps() (*gomarathon.Response, error)
	GetApp(string) (*gomarathon.Response, error)
	CreateGroup(*gomarathon.Group) (*gomarathon.Response, error)
}

type marathonAdapter struct {
	client gomarathonClientAbstractor
	conv   PanamaxServiceConverter
}

func NewMarathonAdapter(endpoint string) *marathonAdapter {
	adapter := new(marathonAdapter)
	adapter.client = newClient(endpoint)
	adapter.conv = new(MarathonConverter)
	return adapter
}

func (m *marathonAdapter) GetServices() ([]*api.Service, *api.Error) {
	var apiErr *api.Error

	response, err := m.client.ListApps()
	if err != nil {
		apiErr = api.NewError(0, err.Error())
	}
	return m.conv.convertToServices(response.Apps), apiErr
}

func (m *marathonAdapter) GetService(id string) (*api.Service, *api.Error) {
	var apiErr *api.Error

	response, err := m.client.GetApp(id)
	if err != nil {
		apiErr = api.NewError(0, err.Error())
	}
	return m.conv.convertToService(response.App), apiErr
}

func (m *marathonAdapter) CreateServices(services []*api.Service) ([]*api.Service, *api.Error) {
	var apiErr *api.Error
	group := new(gomarathon.Group)

	group.ID = "pmx"
	group.Apps = m.conv.convertToApps(services)

	_, err := m.client.CreateGroup(group)
	if err != nil {
		apiErr = api.NewError(0, err.Error())
	}
	return make([]*api.Service, 0), apiErr
}

func (m *marathonAdapter) UpdateService(s *api.Service) *api.Error {
	return nil
}

func (m *marathonAdapter) DestroyService(id string) *api.Error {
	return nil
}
