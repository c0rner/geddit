package rego

import (
	"encoding/json"
	"fmt"
	"io"
)

// buildURL returns a URI for API-method. If 'secure' is true the scheme will be set to https.
func buildURL(method string, secure bool) string {
	scheme := "http"
	if secure {
		scheme += "s"
	}

	return fmt.Sprintf("%s://%s%s", scheme, strReddit, method)
}

// getJSON is a convenience function used by all JSON API methods
func getJSON(rc io.Reader) (*jsonAPIReply, error) {
	r := struct {
		JSON jsonAPIReply `json:"json"`
	}{}
	err := json.NewDecoder(rc).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r.JSON, newAPIError(&r.JSON)
}
