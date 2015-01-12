package rego

import (
	"bytes"
	"fmt"
	"testing"
)

func Test_buildURL(t *testing.T) {
	//func Test_buildURL(method string, secure bool) string {
	var tests = []struct {
		method string
		secure bool
		result string
	}{
		{"/notls", false, "http://www.reddit.com/notls"},
		{"/withtls", true, "https://www.reddit.com/withtls"},
	}

	for _, test := range tests {
		s := buildURL(test.method, test.secure)
		if s != test.result {
			fmt.Printf("%s != %s\n", s, test.result)
			t.Fail()
		}
	}
}

func Test_getJSON(t *testing.T) {
	var tests = []struct {
		json      string
		err       string
		ratelimit bool
	}{
		{"{\"json\": {\"ratelimit\": 1234, \"errors\": [[\"RATELIMIT\", \"Test\",     \"vdelay\"]]}}", "RATELIMIT: Test", true},
		{"{\"json\": {\"errors\": [[\"WRONG_PASSWORD\", \"Test\",     \"passwd\"]]}}", "WRONG_PASSWORD: Test", false},
	}

	for _, test := range tests {
		_, err := getJSON(bytes.NewBufferString(test.json))
		if err.Error() != test.err {
			fmt.Printf("Error is '%s', should be '%s'\n", err, test.err)
			t.Fail()
		}

		if err.(APIError).IsRatelimited() != test.ratelimit {
			fmt.Printf("IsRatelimited '%t', should be '%t'\n", err.(APIError).IsRatelimited(), test.ratelimit)
			t.Fail()
		}
	}
	//func Test_getJSON(rc io.ReadCloser) (*jsonAPIReply, error) {
}
