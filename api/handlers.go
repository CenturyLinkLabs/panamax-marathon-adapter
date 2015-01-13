package api

import (
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/codegangsta/martini"
)


func getServices(enc encoder, adapter PanamaxAdapter) (int, string) {

	res, _ := adapter.GetServices()
	data := handleError(enc.Encode(res))
	// Convert the apps that are being returned into Services
	return http.StatusOK, data
}

func getService(enc encoder, adapter PanamaxAdapter, params martini.Params) (int, string) {

	id := params["id"]
	if id == "" {
		// empty data
		return http.StatusNotFound, fmt.Sprintf("id %s - cannot be empty", id)
	}
	res, _ := adapter.GetService(id)

	//data := utils.HandleError(enc.Encode(res))
	// Convert the apps that are being returned into Services
	return http.StatusOK, res.Id
}


func createService(enc encoder, adapter PanamaxAdapter, r *http.Request) (int, string) {

	var services []*Service

	json.NewDecoder(r.Body).Decode(&services)
	res, _ := adapter.CreateServices(services)
	data := handleError(enc.Encode(res))

	return http.StatusOK, data
}

func deleteService(enc encoder, adapter PanamaxAdapter, params martini.Params) (int, string) {

	id := params["id"]
	if id == "" {
		// empty data
		return http.StatusNotFound, fmt.Sprintf("id %s - cannot be empty", id)
	}
	adapter.DestroyService(id)

	//data := utils.HandleError(enc.Encode(res))
	// Convert the apps that are being returned into Services
	return http.StatusOK, ""
}
