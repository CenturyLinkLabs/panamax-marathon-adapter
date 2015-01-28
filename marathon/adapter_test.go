package marathon

import (
	"github.com/CenturyLinkLabs/panamax-marathon-adapter/api"
	"github.com/jbdalido/gomarathon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// Mock the gomarathonClient type that implements the gomarathonClientAbstractor interface
type mockClient struct {
	mock.Mock
}

func (c *mockClient) ListApps() (*gomarathon.Response, error) {
	args := c.Mock.Called()
	if len(args) == 1 {
		return args.Get(0).(*gomarathon.Response), nil
	} else {
		return args.Get(0).(*gomarathon.Response), args.Error(1)
	}
}

func (c *mockClient) GetApp(id string) (*gomarathon.Response, error) {
	args := c.Mock.Called(id)
	if len(args) == 1 {
		return args.Get(0).(*gomarathon.Response), nil
	} else {
		return args.Get(0).(*gomarathon.Response), args.Error(1)
	}
}

func (c *mockClient) GetAppTasks(id string) (*gomarathon.Response, error) {
	args := c.Mock.Called(id)
	if len(args) == 1 {
		return args.Get(0).(*gomarathon.Response), nil
	} else {
		return args.Get(0).(*gomarathon.Response), args.Error(1)
	}
}

func (c *mockClient) CreateGroup(group *gomarathon.Group) (*gomarathon.Response, error) {
	args := c.Mock.Called(group)
	if len(args) == 1 {
		return args.Get(0).(*gomarathon.Response), nil
	} else {
		return args.Get(0).(*gomarathon.Response), args.Error(1)
	}
}

func (c *mockClient) DeleteApp(id string) (*gomarathon.Response, error) {
	args := c.Mock.Called(id)
	if len(args) == 1 {
		return args.Get(0).(*gomarathon.Response), nil
	} else {
		return args.Get(0).(*gomarathon.Response), args.Error(1)
	}
}

func (c *mockClient) DeleteGroup(id string) (*gomarathon.Response, error) {
	args := c.Mock.Called(id)
	if len(args) == 1 {
		return args.Get(0).(*gomarathon.Response), nil
	} else {
		return args.Get(0).(*gomarathon.Response), args.Error(1)
	}
}

func (c *mockClient) CreateApp(app *gomarathon.Application) (*gomarathon.Response, error) {
	args := c.Mock.Called(app)
	if len(args) == 1 {
		return args.Get(0).(*gomarathon.Response), nil
	} else {
		return args.Get(0).(*gomarathon.Response), args.Error(1)
	}
}

// Mock the MarathonConverter type that implements the PanamaxConverter interface
type mockConverter struct {
	mock.Mock
}

func (c *mockConverter) convertToServices(apps []*gomarathon.Application) []*api.Service {
	args := c.Mock.Called(apps)
	return args.Get(0).([]*api.Service)
}

func (c *mockConverter) convertToService(app *gomarathon.Application) *api.Service {
	args := c.Mock.Called(app)
	return args.Get(0).(*api.Service)
}

func (c *mockConverter) convertToApps(services []*api.Service) []*gomarathon.Application {
	args := c.Mock.Called(services)
	return args.Get(0).([]*gomarathon.Application)
}

func (c *mockConverter) convertToApp(service *api.Service) *gomarathon.Application {
	args := c.Mock.Called(service)
	return args.Get(0).(*gomarathon.Application)
}

// Tests
func TestMarathonAdapterImplementsPanamaxAdapterInterface(t *testing.T) {
	assert.Implements(t, (*api.PanamaxAdapter)(nil), new(marathonAdapter))
}

func TestMockClientImplementsGoMarathonClientAbstractorInterface(t *testing.T) {
	assert.Implements(t, (*gomarathonClientAbstractor)(nil), new(mockClient))
}

func TestMockConverterImplementsPanamaxServiceConverterInterface(t *testing.T) {
	assert.Implements(t, (*PanamaxServiceConverter)(nil), new(mockConverter))
}

func setup() (*mockClient, *mockConverter, *marathonAdapter) {
	testClient := new(mockClient)
	testConverter := new(mockConverter)
	adapter := new(marathonAdapter)
	adapter.client = testClient
	adapter.conv = testConverter
	adapter.generateUID = func() string { return "pmx" }

	return testClient, testConverter, adapter
}

func TestSuccessfulGetServices(t *testing.T) {

	// setup
	testClient, testConverter, adapter := setup()

	resp := new(gomarathon.Response)
	resp.Apps = make([]*gomarathon.Application, 0)
	services := make([]*api.Service, 0)

	// set expectations
	testClient.On("ListApps").Return(resp)
	testConverter.On("convertToServices", resp.Apps).Return(services)

	// call the code to be tested
	srvcs, err := adapter.GetServices()

	// assert if expectations are met
	assert.NoError(t, err)
	assert.Len(t, srvcs, 0)

	testClient.AssertExpectations(t)
	testConverter.AssertExpectations(t)

}

func TestSuccessfulGetService(t *testing.T) {

	// setup
	testClient, testConverter, adapter := setup()

	resp := new(gomarathon.Response)
	resp.App = &gomarathon.Application{ID: "foo"}
	service := &api.Service{Id: "foo", Name: "foo"}

	// set expectations
	testClient.On("GetApp", "/foo").Return(resp)
	testConverter.On("convertToService", resp.App).Return(service)

	// call the code to be tested
	srvc, err := adapter.GetService("foo")

	// assert if expectations are met
	assert.NoError(t, err)
	assert.IsType(t, new(api.Service), srvc)
	assert.Equal(t, service, srvc)

	testClient.AssertExpectations(t)
	testConverter.AssertExpectations(t)

}

/*
NOTE: Note sure how to test this with the channels now!

func TestSuccessfulCreateServices(t *testing.T) {

	// setup
	testClient, testConverter, adapter := setup()

	//resp := new(gomarathon.Response)
	services := make([]*api.Service, 1)
	services[0] = &api.Service{Name: "foo"}
	group := new(gomarathon.Group)
	group.ID = "pmx"
	group.Apps = make([]*gomarathon.Application, 0)

	// set expectations
	//testClient.On("CreateApp", group).Return(resp)
	testConverter.On("convertToApps", services).Return(group.Apps)

	// call the code to be tested
	_, err := adapter.CreateServices(services)

	// assert if expectations are met
	assert.NoError(t, err)

	testClient.AssertExpectations(t)
	testConverter.AssertExpectations(t)

}
*/

func TestSuccessfulDeleteService(t *testing.T) {

	// setup
	testClient, _, adapter := setup()

	resp := new(gomarathon.Response)

	// set expectations
	testClient.On("DeleteApp", "/foo").Return(resp)
	testClient.On("DeleteGroup", "").Return(resp)

	// call the code to be tested
	err := adapter.DestroyService("foo")

	// assert if expectations are met
	assert.NoError(t, err)

	testClient.AssertExpectations(t)

}
