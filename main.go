package geddit

import (
	"encoding/json"
	"fmt"
)

// Reddit API methods
const (
	APPClear   = "/api/clear_sessions"
	APIComment = "/api/comment"
	APIDelete  = "/api/del"
	APIListing = "/r/%s.json"
	APILogin   = "/api/login"
	APIMe      = "/api/me.json"
	APISubmit  = "/api/submit"
)

const (
	StrReddit = "www.reddit.com"
	StrCookie = "reddit_session"
)

// AuthConfig is used when authenticating a Reddit session using User and Password.
// When calling Login if a session cookie already is set it will be discarded regardless
// of authentication status.
type Authconfig struct {
	User     string // Reddit username
	Password string // Authentication password
}

// jsonApiReply is only used for json api replies
// for return codes.
type jsonAPIReply struct {
	Data      json.RawMessage `json:"data"` // A data structure formatted based on kind
	Errors    [][]string      `json:"errors,omitempty"`
	ID        string          `json:"id"`   // Item identifier, e.g. "c3v7f8u"
	Kind      string          `json:"kind"` // Kind denotes the item's type.
	Name      string          `json:"name"` // Fullname of item, e.g. "t1_c3v7f8u"
	Ratelimit float64         `json:"ratelimit,omitempty"`
}

// BuildURL returns a URI for API-method. If 'secure' is true the scheme will be set to https.
func BuildURL(method string, secure bool) string {
	scheme := "http"
	if secure {
		scheme += "s"
	}

	return fmt.Sprintf("%s://%s%s", scheme, StrReddit, method)
}
