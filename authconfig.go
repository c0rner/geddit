package geddit

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// AuthConfig is used to authenticating a Reddit session using User/Password or Cookie.
// If Cookie is set it takes precedence and User/Password will not be used.
type Authconfig struct {
	Cookie   string // Pre-authenticated cookie
	Password string // Authentication password
	User     string // Reddit username
}

// LoadFromPath parses a JSON file into Authconfig
func (a *Authconfig) LoadFromPath(p string) error {
	buf, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, a)
	if err != nil {
		return err
	}

	return nil
}

// SaveToPath writes Authconfig as JSON to path
func (a *Authconfig) SaveToPath(p string) error {
	// Always overwrite, no questions asked
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	// Make sure file is readable by the running
	// user only before we write session data
	err = f.Chmod(0600)
	if err != nil {
		return err
	}

	buf, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}
	_, err = f.Write(buf)
	return err
}
