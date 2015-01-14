package api

import (
	"encoding/json"
	"net/http"

	"github.com/codegangsta/martini"
)

func sanitizeErrorCode(code int) int {
	if http.StatusText(code) == "" {
		return http.StatusInternalServerError
	}
	return code
}

func getServices(e encoder, adapter PanamaxAdapter) (int, string) {
	data, err := adapter.GetServices()
	if err != nil {
		return sanitizeErrorCode(err.Code), err.Message
	}

	return http.StatusOK, e.Encode(data)
}

func getService(e encoder, adapter PanamaxAdapter, params martini.Params) (int, string) {
	id := params["id"]

	data, err := adapter.GetService(id)
	if err != nil {
		return sanitizeErrorCode(err.Code), err.Message
	}

	return http.StatusOK, e.Encode(data)
}

func createService(e encoder, adapter PanamaxAdapter, r *http.Request) (int, string) {
	var services []*Service
	json.NewDecoder(r.Body).Decode(&services)

	res, err := adapter.CreateServices(services)
	if err != nil {
		return sanitizeErrorCode(err.Code), err.Message
	}

	return http.StatusCreated, e.Encode(res)

}

func updateService(adapter PanamaxAdapter, params martini.Params, r *http.Request) (int, string) {
	service := new(Service)
	id := params["id"]

	json.NewDecoder(r.Body).Decode(&service)
	service.Id = id

	err := adapter.UpdateService(service)
	if err != nil {
		return sanitizeErrorCode(err.Code), err.Message
	}

	return http.StatusNoContent, ""
}

func deleteService(adapter PanamaxAdapter, params martini.Params) (int, string) {
	id := params["id"]

	err := adapter.DestroyService(id)
	if err != nil {
		return sanitizeErrorCode(err.Code), err.Message
	}

	return http.StatusNoContent, ""
}
