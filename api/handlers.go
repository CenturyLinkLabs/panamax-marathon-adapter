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

func createServices(e encoder, adapter PanamaxAdapter, r *http.Request) (int, string) {
	var services []*Service
	json.NewDecoder(r.Body).Decode(&services)

	res, err := adapter.CreateServices(services)
	if err != nil {
		return sanitizeErrorCode(err.Code), err.Message
	}

	return http.StatusCreated, e.Encode(res)
}

func updateService(adapter PanamaxAdapter, params martini.Params, r *http.Request) (int, string) {
	return http.StatusNotImplemented, ""
}

func deleteService(adapter PanamaxAdapter, params martini.Params) (int, string) {
	id := params["id"]

	err := adapter.DestroyService(id)
	if err != nil {
		return sanitizeErrorCode(err.Code), err.Message
	}

	return http.StatusNoContent, ""
}

func getMetadata(e encoder, adapter PanamaxAdapter) (int, string) {

	data := &Metadata{Version: VERSION, Type: "marathon"}

	return http.StatusOK, e.Encode(data)
}
