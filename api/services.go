package api

import (
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/codegangsta/martini"
	"github.com/centurylinklabs/panamax-marathon-adapter/utils"
)


func getServices(enc utils.Encoder, adapter PanamaxAdapter) (int, string) {

	res := adapter.GetServices()
	data := utils.HandleError(enc.Encode(res))
	// Convert the apps that are being returned into Services
	return http.StatusOK, data
}

func getService(enc utils.Encoder, adapter PanamaxAdapter, params martini.Params) (int, string) {

	id := params["id"]
	if id == "" {
		// empty data
		return http.StatusNotFound, fmt.Sprintf("id %s - cannot be empty", id)
	}
	res := adapter.GetService(id)

	//data := utils.HandleError(enc.Encode(res))
	// Convert the apps that are being returned into Services
	return http.StatusOK, res.Id
}


func createService(adapter PanamaxAdapter, r *http.Request, enc utils.Encoder) (int, string) {

	var services []*Service

	json.NewDecoder(r.Body).Decode(&services)
	res := adapter.CreateServices(services)
	data := utils.HandleError(enc.Encode(res))

	return http.StatusOK, data
}
