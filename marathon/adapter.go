package marathon // import "github.com/CenturyLinkLabs/panamax-marathon-adapter/marathon"

import (
	"log"
	"fmt"
	"strings"

	"github.com/CenturyLinkLabs/gomarathon"
	"github.com/CenturyLinkLabs/panamax-marathon-adapter/api"
	"github.com/satori/go.uuid"
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
	DeleteApp(string) (*gomarathon.Response, error)
	DeleteGroup(string) (*gomarathon.Response, error)
}

type marathonAdapter struct {
	client gomarathonClientAbstractor
	conv   PanamaxServiceConverter
	generateUID func() string
}

func NewMarathonAdapter(endpoint string) *marathonAdapter {
	adapter := new(marathonAdapter)
	adapter.client = newClient(endpoint)
	adapter.conv = new(MarathonConverter)
	adapter.generateUID = func() string { return fmt.Sprintf("%s",uuid.NewV4()) }
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

	response, err := m.client.GetApp(sanitizeServiceId(id))
	if err != nil {
		apiErr = api.NewError(0, err.Error())
	}
	return m.conv.convertToService(response.App), apiErr
}

func (m *marathonAdapter) CreateServices(services []*api.Service) ([]*api.Service, *api.Error) {
	var apiErr *api.Error
	group := new(gomarathon.Group)

	group.ID = m.generateUID()
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
	var apiErr *api.Error
	group, _ := splitServiceId(id)

	_, err := m.client.DeleteApp(sanitizeServiceId(id))
	if err != nil {
		apiErr = api.NewError(0, err.Error())
	}

	m.client.DeleteGroup(group) // Remove group if possible we dont care about error or return.

	return apiErr
}

// Split the service string into 2 parts part[0] is group part[1] is service
func splitServiceId(serviceId string) (string, string) {
	var group, service string

	parts := strings.Split(serviceId, ".")
	if len(parts) == 2 {
		group = parts[0]
		service = parts[1]
	} else {
		service = parts[0]
	}
	return group, service
}

func sanitizeServiceId(id string) string {
	group, service := splitServiceId(id)
	return fmt.Sprintf("%s/%s", group, service)
}


