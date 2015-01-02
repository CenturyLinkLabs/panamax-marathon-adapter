package api

import (
	"fmt"
	"net/http"
	"github.com/codegangsta/martini"
	"github.com/centurylinklabs/panamax-marathon-adapter/marathon"
	"github.com/centurylinklabs/panamax-marathon-adapter/utils"
)

func GetServices(enc utils.Encoder) (int, string) {
	// get Marathon apps
	res, err := marathon.ListApplications()
	if err != nil {
		return http.StatusNotFound, fmt.Sprintf("%s", err)
	}
	data := utils.HandleError(enc.Encode(res))
	// Convert the apps that are being returned into Services
	return http.StatusOK, data
}

func GetService(enc utils.Encoder, params martini.Params) (int, string) {
	id := params["id"]
	if id == "" {
		// empty data
		return http.StatusNotFound, fmt.Sprintf("id %s - cannot be empty", id)
	}
	// get Marathon app by id
	res, err := marathon.GetApplication(id)
	if err != nil {
		return http.StatusNotFound, fmt.Sprintf("%s", err)
	}
	data := utils.HandleError(enc.Encode(res))
	// Convert the apps that are being returned into Services
	return http.StatusOK, data
}


