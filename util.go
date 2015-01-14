package rego

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
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

func timeFromNumber(n json.Number) time.Time {
	var s, u int
	parts := strings.Split(n.String(), ".")
	s, _ = strconv.Atoi(parts[0]) // Seconds
	if len(parts) > 1 {
		u, _ = strconv.Atoi(parts[1]) // Microseconds
	}
	return time.Unix(int64(s), int64(u))
}
