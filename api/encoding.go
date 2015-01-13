package api

import (
	"encoding/json"
	"log"
)

type jsonEncoder struct{}

// An Encoder implements an encoding format of values to be sent as response to
// requests on the API endpoints.
type encoder interface {
	Encode(v ...interface{}) (string, error)
}

// Because `panic`s are caught by martini's Recovery handler, it can be used
// to return server-side errors (500). Some helpful text message should probably
// be sent, although not the technical error (which is printed in the log).
func handleError(data string, err error) string {
	if err != nil {
		panic(err)
	}
	return data
}

// JsonEncoder is an Encoder that produces JSON-formatted responses.
func (_ jsonEncoder) Encode(v ...interface{}) (string, error) {
	var data interface{} = v
	if v == nil {
		// So that empty results produces `[]` and not `null`
		data = []interface{}{}
	} else if len(v) == 1 {
		data = v[0]
	}
	b, err := json.Marshal(data)
	log.Printf("%s", b)
	return string(b), err
}

