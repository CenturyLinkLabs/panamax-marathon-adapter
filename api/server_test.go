package api

import (
	"testing"
	"io"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"

)

var testServer *httptest.Server

func init() {
	martini := NewServer(new(NoOPAdapter))
	testServer = httptest.NewServer(martini.svr)
}


type NoOPAdapter struct {
}

func (NoOPAdapter) GetServices() ([]*Service, *Error) {
	return make([]*Service, 0), nil
}
func (NoOPAdapter) GetService(string) (*Service, *Error) {
	service := new(Service)
	return service, nil
}
func (NoOPAdapter) CreateServices([]*Service)([]*Response, *Error) {
	return make([]*Response, 0), nil
}
func (NoOPAdapter) UpdateService(*Service) (*Error) {
	return nil
}
func (NoOPAdapter) DestroyService(string) (*Error) {
	return nil
}


func TestGetServices(t *testing.T) {
	res, _ := http.Get(fmt.Sprintf("%s/services",testServer.URL))

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetService(t *testing.T) {
	res, _ := http.Get(fmt.Sprintf("%s/services/1",testServer.URL))

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPostService(t *testing.T) {
	var body io.Reader
	res, _ := http.Post(fmt.Sprintf("%s/services",testServer.URL), "", body)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestDeleteService(t *testing.T) {
	var body io.Reader
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/services/1",testServer.URL), body)
	res, _ := http.DefaultClient.Do(req)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestNoRoute(t *testing.T) {
	res, _ := http.Get(fmt.Sprintf("%s/nothere",testServer.URL))

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}






