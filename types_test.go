package rego

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func Test_Edited(t *testing.T) {
	var tests = []struct {
		json   string
		status bool
		date   time.Time
	}{
		{"{\"edited\": false }", false, time.Unix(0, 0)},
		{"{\"edited\": true }", true, time.Unix(0, 0)},
		{"{\"edited\": 1234567890.0 }", true, time.Unix(1234567890, 0)},
	}

	var data = struct {
		Edited Edited
	}{}
	for _, test := range tests {
		err := json.NewDecoder(strings.NewReader(test.json)).Decode(&data)
		if err != nil {
			t.Error(err)
			continue
		}
		if test.status != data.Edited.Status || test.date != data.Edited.UTC {
			t.Errorf("Got: %t, %s, Wanted: %t, %s", data.Edited.Status, data.Edited.UTC, test.status, test.date)
		}
	}
}
