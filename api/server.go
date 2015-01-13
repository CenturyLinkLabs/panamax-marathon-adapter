package api

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"github.com/codegangsta/martini"
)

type martiniServer struct {
	svr *martini.Martini
}

func NewServer(adapterInst PanamaxAdapter) (*martiniServer) {
	s := martini.New()

	// Setup middleware
	s.Use(martini.Recovery())
	s.Use(martini.Logger())
	s.Use(mapEncoder)
	s.Use(func (c martini.Context, w http.ResponseWriter, r *http.Request) {
		c.Map(adapterInst)
	})
	// Setup routes
	router := martini.NewRouter()
	router.Get(`/services`, getServices)
	router.Get(`/services/:id`, getService)
	router.Post(`/services`, createService)
	router.Delete(`/services/:id`, deleteService)
	// Add the router action
	s.Action(router.Handle)
	server := martiniServer{svr: s}

	return &server
}

func (m *martiniServer) Start() {
	err := http.ListenAndServe(":8001", m.svr)
	if	err != nil {
		log.Fatal(err)
	}
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
		c.MapTo(jsonEncoder{}, (*encoder)(nil))
		w.Header().Set("Content-Type", "application/json")
	}
}
