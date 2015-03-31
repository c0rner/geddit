package rego

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
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

func Test_timeFromNumber(t *testing.T) {
	var tests = []struct {
		number json.Number
		time   time.Time
	}{
		{"12345.6789", time.Unix(12345, 6789)},
		{"12345", time.Unix(12345, 0)},
		{"123oops456.789", time.Unix(0, 789)},
		{"123oops456", time.Unix(0, 0)},
	}

	for _, test := range tests {
		if timeFromNumber(test.number) != test.time {
			fmt.Printf("%s != %s\n", timeFromNumber(test.number), test.time)
			t.Fail()
		}
	}
}
