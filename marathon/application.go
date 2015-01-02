package marathon

import (
	"github.com/jbdalido/gomarathon"
	"log"
)

const (
	apiEndpoint = "http://206.128.155.112:8080"
)

var client *gomarathon.Client

func init() {
	client = NewClient("")
}

func ListApplications() (*gomarathon.Response, error) {
	// List all apps
	return client.ListApps()
}

func GetApplication(id string) (*gomarathon.Response, error) {
	// Get an app
	return client.GetApp(id)
}

func NewClient(endpoint string) *gomarathon.Client {
	url := apiEndpoint
	if endpoint != "" {
		url = endpoint
	}
	log.Printf("Marathon API: %s", url)
	c, err := gomarathon.NewClient(url, nil)
	if err != nil {
		log.Fatal(err)
	}
	return c
}
