package rego

import (
	"bytes"
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

func TestAPIError(t *testing.T) {
	var tests = []struct {
		json      string
		hasError  bool
		ratelimit int
	}{
		{"{\"json\": {\"ratelimit\": 600, \"errors\": [[\"RATELIMIT\", \"Test\",     \"vdelay\"]]}}", true, 599},
		{"{\"json\": {\"errors\": [[\"WRONG_PASSWORD\", \"Test\",     \"passwd\"]]}}", true, 0},
		{"{\"json\": {totallybroken:  }}", true, 0},
		{"{\"json\": { }}", false, 0},
	}

	for _, test := range tests {
		json, err := getJSON(bytes.NewBufferString(test.json))
		if err != nil {
			if !test.hasError {
				fmt.Printf("getJSON() failed but should not. (%s)\n", err)
				t.Fail()
			}
			switch err.(type) {
			case APIError:
				actualErr := fmt.Sprintf("%s: %s", json.Errors[0][0], json.Errors[0][1])
				if err.Error() != actualErr {
					fmt.Printf("Error is '%s', should be '%s'\n", err, actualErr)
					t.Fail()
				}
				if err.(APIError).IsRatelimited() && test.ratelimit == 0 {
					fmt.Printf("IsRatelimited returns '%t', expeced ratelimit is '%d'\n", err.(APIError).IsRatelimited(), test.ratelimit)
					t.Fail()
				}
				if int(err.(APIError).Duration().Seconds()) != test.ratelimit {
					fmt.Printf("Duration is %d, expected %d\n", int(err.(APIError).Duration().Seconds()), test.ratelimit)
					t.Fail()
				}
			}
		} else {
			if test.hasError {
				t.Fail()
				fmt.Printf("getJSON() did not failed but should\n")
			}
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
