package geddit

import (
	"encoding/json"
)

// Reddit API methods
const (
	apiClear   = "/api/clear_sessions"
	apiComment = "/api/comment"
	apiDelete  = "/api/del"
	apiListing = "/r/%s.json"
	apiLogin   = "/api/login"
	apiMe      = "/api/me.json"
	apiSubmit  = "/api/submit"
)

const (
	strReddit = "www.reddit.com"
	strCookie = "reddit_session"
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
