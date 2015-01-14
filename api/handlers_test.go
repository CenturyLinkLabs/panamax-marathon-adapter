package api

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testEncoder = new(jsonEncoder)

type MockAdapter struct {
	returnError *Error
}

func (e MockAdapter) SetReturnCode(code int) {

}

func (e MockAdapter) GetServices() ([]*Service, *Error) {
	return nil, e.returnError

}
func (e MockAdapter) GetService(string) (*Service, *Error) {
	return nil, e.returnError
}
func (e MockAdapter) CreateServices([]*Service) ([]*Service, *Error) {
	return nil, e.returnError
}
func (e MockAdapter) UpdateService(*Service) *Error {
	return e.returnError
}
func (e MockAdapter) DestroyService(string) *Error {
	return e.returnError
}

func newMockAdapter(code int, message string) *MockAdapter {
	adapter := new(MockAdapter)
	adapter.returnError = NewError(code, message)

	return adapter
}

func TestCodeOutOfRange(t *testing.T) {
	code, _ := getServices(testEncoder, newMockAdapter(9090, ""))
	assert.Equal(t, http.StatusInternalServerError, code)
}

func TestSuccessfulGetServices(t *testing.T) {
	code, _ := getServices(testEncoder, newMockAdapter(200, ""))

	assert.Equal(t, http.StatusOK, code)
}

func TestSuccessfulGetService(t *testing.T) {
	params := map[string]string{
		"id": "test",
	}

	code, _ := getService(testEncoder, newMockAdapter(200, ""), params)

	assert.Equal(t, http.StatusOK, code)
}

func TestSuccessfulUpdateService(t *testing.T) {
	req, _ := http.NewRequest("PUT", "http://localhost", strings.NewReader("{}"))
	params := map[string]string{
		"id": "test",
	}
	code, _ := updateService(newMockAdapter(204, ""), params, req)

	assert.Equal(t, http.StatusNoContent, code)
}

func TestSuccessfulDeleteService(t *testing.T) {
	params := map[string]string{
		"id": "test",
	}

	code, _ := deleteService(newMockAdapter(204, ""), params)

	assert.Equal(t, http.StatusNoContent, code)
}

func TestGetServicesError(t *testing.T) {
	adapter := newMockAdapter(500, "internal error")
	code, _ := getServices(testEncoder, adapter)

	assert.Equal(t, http.StatusInternalServerError, code)
}

func TestServiceNotFound(t *testing.T) {
	adapter := newMockAdapter(404, "service not found")
	params := map[string]string{
		"id": "test",
	}

	code, body := getService(testEncoder, adapter, params)

	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, "service not found", body)

}

func TestUpdateServiceNotFound(t *testing.T) {
	req, _ := http.NewRequest("PUT", "http://localhost", strings.NewReader("{}"))
	adapter := newMockAdapter(404, "service not found")
	params := map[string]string{
		"id": "test",
	}

	code, body := updateService(adapter, params, req)

	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, "service not found", body)
}

func TestUpdateServiceInvalidState(t *testing.T) {
	req, _ := http.NewRequest("PUT", "http://localhost", strings.NewReader("{}"))
	adapter := newMockAdapter(400, "invalid service state")
	params := map[string]string{
		"id": "test",
	}

	code, body := updateService(adapter, params, req)

	assert.Equal(t, http.StatusBadRequest, code)
	assert.Equal(t, "invalid service state", body)
}

func TestDeleteServiceNotFound(t *testing.T) {
	adapter := newMockAdapter(404, "service not found")
	params := map[string]string{
		"id": "test",
	}

	code, body := deleteService(adapter, params)

	assert.Equal(t, http.StatusNotFound, code)
	assert.Equal(t, "service not found", body)
}
