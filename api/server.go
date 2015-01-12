package api

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"github.com/codegangsta/martini"
	"github.com/centurylinklabs/panamax-marathon-adapter/utils"
)

// The one and only martini instance.
var mServer *martini.Martini
var adapterInstance PanamaxAdapter

func init() {
	mServer = martini.New()
	// Setup middleware
	mServer.Use(martini.Recovery())
	mServer.Use(martini.Logger())
	mServer.Use(mapEncoder)
	mServer.Use(adapter)
	// Setup routes
	r := martini.NewRouter()
	r.Get(`/services`, getServices)
	r.Get(`/services/:id`, getService)
	r.Post(`/services`, createService)

	// Add the router action
	mServer.Action(r.Handle)
}

// The regex to check for the requested format (allows an optional trailing
// slash)
var rxExt = regexp.MustCompile(`(\.(?:json))\/?$`)

// MapEncoder intercepts the request's URL, detects the requested format,
// and injects the correct encoder dependency for this request. It rewrites
// the URL to remove the format extension, so that routes can be defined
// without it.
func mapEncoder(c martini.Context, w http.ResponseWriter, r *http.Request) {
	// Get the format extension
	matches := rxExt.FindStringSubmatch(r.URL.Path)
	ft := ".json"
	if len(matches) > 1 {
		// Rewrite the URL without the format extension
		l := len(r.URL.Path) - len(matches[1])
		if strings.HasSuffix(r.URL.Path, "/") {
			l--
		}
		r.URL.Path = r.URL.Path[:l]
		ft = matches[1]
	}
	// Inject the requested encoder
	switch ft {
	// Add cases for other formats
	default:
		c.MapTo(utils.JsonEncoder{}, (*utils.Encoder)(nil))
		w.Header().Set("Content-Type", "application/json")
	}
}

func adapter(c martini.Context, w http.ResponseWriter, r *http.Request) {
	c.Map(adapterInstance)
}

func ListenAndServe(theAdapter PanamaxAdapter) {
        adapterInstance = theAdapter
	err := http.ListenAndServe(":8001", mServer)
	if	err != nil {
		log.Fatal(err)
	}
}
