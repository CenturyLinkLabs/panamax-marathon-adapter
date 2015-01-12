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

type marathonAdapter struct {
	client *gomarathon.Client
}

func NewMarathonAdapter(endpoint string) (*marathonAdapter) {
	adapter := new(marathonAdapter)
	adapter.client = newClient(endpoint)
	return adapter
}


func (m *marathonAdapter) GetServices() ([]*api.Service) {
	response, _ := m.client.ListApps()

	return convertApps(response.Apps)
}

func (m *marathonAdapter) GetService(id string) (*api.Service) {
	response, _ := m.client.GetApp(id)

	return convertApp(response.App)
}

func (m *marathonAdapter) CreateServices(services []*api.Service) ([]*api.Response) {
	group := new(gomarathon.Group)
	//res := new(api.Response)

	group.ID = "pmx"
	group.Apps = convertServices(services)
	m.client.CreateGroup(group)

	return make([]*api.Response, 0)
}

func (m *marathonAdapter) UpdateService(s *api.Service) (bool) {
	return false;
}

func (m *marathonAdapter) DestroyService(s *api.Service) (bool) {
	return false;
}

