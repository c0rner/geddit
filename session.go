package geddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type RateLimit struct {
	Remaining int
	Reset     int
	Used      int
}

// Session is an active Reddit session that initially is unauthenticated. An authenticated
// session can be set up using Session.Login or by updating Session.Cookie manually.
type Session struct {
	client    *http.Client
	Cookie    string // Session cookie (empty if not logged in)
	modhash   string
	RateLimit RateLimit // RateLimit usage is updated on each API request
	useragent string
}

// NewSession creates an unauthenticated Reddit session
func NewSession(ua string) *Session {
	return &Session{
		client:    &http.Client{},
		useragent: ua,
	}
}

func (s *Session) Me() (*Account, error) {
	resp, err := s.get(BuildURL(APIMe, true), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	thing := Thing{}
	err = json.NewDecoder(resp.Body).Decode(&thing)
	if err != nil {
		return nil, err
	}

	if thing.Kind != TypeAccount {
		// TODO.. handle error
	}

	account := Account{}
	err = json.Unmarshal(thing.Data, &account)
	if err != nil {
		return nil, err
	}

	s.modhash = account.Modhash
	return &account, nil
}

// Listing returns a list of new Link items using supplied Paginator page.
// If page is nil the Reddit defaults are used.
func (s *Session) Listing(sub string) *Paginator {
	p := Paginator{}
	p.s = s
	p.url = fmt.Sprintf(BuildURL(APIListing, true), sub)
	return &p
}

// Comment posts a reply to parent post p using the raw text t returning
// new new comments fullname id
func (s *Session) Comment(p string, t string) (*CommentResult, error) {
	v := url.Values{"api_type": {"json"}}
	v.Set("thing_id", p)
	v.Set("text", t)

	resp, err := s.post(BuildURL(APIComment, true), v)
	/*
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		fmt.Printf("Data: %#v\n", buf.String())
	*/

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	r, err := getJSON(resp.Body)
	if err != nil {
		return nil, err
	}

	container := struct {
		Things []struct {
			Kind string
			Data json.RawMessage
		}
	}{}
	err = json.Unmarshal(r.Data, &container)
	if err != nil {
		return nil, err
	}

	cr := CommentResult{}
	err = json.Unmarshal(container.Things[0].Data, &cr) // FIXME, verify Things[] is not empty

	if err != nil {
		return nil, err
	}

	return &cr, nil
}

// Login authenticates the current session
func (s *Session) Login(ac *Authconfig) error {
	if ac == nil || len(ac.User) == 0 || len(ac.Password) == 0 {
		return errors.New("No authentixation credentials")
	}
	v := url.Values{"api_type": {"json"}}
	v.Set("user", ac.User)
	v.Set("passwd", ac.Password)

	s.Cookie = ""
	s.modhash = ""

	resp, err := s.post(BuildURL(APILogin, true), v)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	if err == nil {
		// FIXME, the needs to be refactored, possibly moved into a c.Do() wrapper
		// Atoi errors can safely be ignored as 0 will be returned on bad input
		// and we will not need to check that each header really exists
		s.RateLimit.Used, _ = strconv.Atoi(resp.Header.Get("X-Ratelimit-Used"))
		s.RateLimit.Reset, _ = strconv.Atoi(resp.Header.Get("X-Ratelimit-Reset"))
		s.RateLimit.Remaining, _ = strconv.Atoi(resp.Header.Get("X-Ratelimit-Remaining"))
	}

	r, err := getJSON(resp.Body)
	if err != nil {
		return err
	}

	reply := struct {
		NeedHttps bool   `json:"need_https"`
		Modhash   string `json:"modhash"`
		Cookie    string `json:"cookie"`
	}{}
	err = json.Unmarshal(r.Data, &reply)
	if err != nil {
		return err
	}

	// Get the session cookie and modhash from response
	/*
		for _, c := range resp.Cookies() {
			if c.Name == StrCookie {
				s.Cookie = c.String()
			}
		}
	*/
	s.Cookie = fmt.Sprintf("%s=%s", StrCookie, reply.Cookie)
	s.modhash = reply.Modhash

	return nil
}

func (s *Session) httpHeaders() http.Header {
	h := http.Header{}
	h.Set("User-Agent", s.useragent)
	if len(s.Cookie) != 0 {
		h.Set("Cookie", s.Cookie)
	}
	if len(s.modhash) != 0 {
		h.Set("X-Modhash", s.modhash)
	}
	return h
}

func (s *Session) get(u string, v url.Values) (*http.Response, error) {
	var values string
	if v != nil {
		values = fmt.Sprintf("?%s", v.Encode())
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", u, values), nil)
	if err != nil {
		return nil, err
	}
	req.Header = s.httpHeaders()
	return s.client.Do(req)
}

func (s *Session) post(u string, v url.Values) (*http.Response, error) {
	if v == nil {
		return nil, errors.New("No values supplied")
	}
	req, err := http.NewRequest("POST", u, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header = s.httpHeaders()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return s.client.Do(req)
}

// getJSON is a convenience function used by all JSON API methods
func getJSON(rc io.ReadCloser) (*jsonAPIReply, error) {
	r := struct {
		JSON jsonAPIReply `json:"json"`
	}{}
	err := json.NewDecoder(rc).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r.JSON, newAPIError(&r.JSON)
}
