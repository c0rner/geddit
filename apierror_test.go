package rego

import (
	"bytes"
	"fmt"
	"testing"
)

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
					fmt.Printf("IsRatelimited returns '%t', expeced %t\n", err.(APIError).IsRatelimited(), test.ratelimit == 0)
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
