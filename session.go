package rego

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

var (
	ErrBadCookie = errors.New("Bad cookie")
)

type RateLimit struct {
	Remaining int
	Reset     int
	Used      int
}

// Session is an active Reddit session that initially is unauthenticated. An authenticated
// session can be set up using Session.Login or Session.SetCookie.
type Session struct {
	client    *http.Client
	Cookie    string // Session cookie (empty if not logged in)
	modhash   string
	RateLimit RateLimit // RateLimit usage is updated on each API request
	useragent string
	lock      sync.Mutex
}

// NewSession creates an unauthenticated Reddit session
func NewSession(ua string) *Session {
	return &Session{
		client:    &http.Client{},
		useragent: ua,
	}
}

// Me returns Account type populated with data for the currently
// authenticated user.  This is equivalent to using Session.User()
// and providing the authenticated username.
func (s *Session) Me() (*Account, error) {
	resp, err := s.get(buildURL(apiMe, true), nil)
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

	return &account, nil
}

// User returns Account type populated with data for user u.
func (s *Session) User(u string) (*Account, error) {
	url := fmt.Sprintf(buildURL(apiUserAbout, true), u)
	resp, err := s.get(url, nil)
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

	return &account, nil
}

// Listing returns a paginated Listing wrapped in a Page type
func (s *Session) Listing(sub string) *Page {
	p := Page{}
	p.s = s
	p.url = fmt.Sprintf(buildURL(apiListing, true), sub)
	return &p
}

// Comment posts a reply to parent post p using the raw text t.
// A successfull post will return the new comments fullname id.
func (s *Session) Comment(p string, t string) (*CommentResult, error) {
	v := url.Values{"api_type": {"json"}}
	v.Set("thing_id", p)
	v.Set("text", t)

	resp, err := s.post(buildURL(apiComment, true), v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

// SetCookie authenticates the current session using a pre-authenticated cookie
func (s *Session) SetCookie(c string) error {
	s.Cookie = c
	acct, err := s.Me()
	if err != nil {
		return err
	}
	if len(acct.Modhash) == 0 {
		return ErrBadCookie
	}

	s.modhash = acct.Modhash
	return nil
}

// Login authenticates the current session using username and password
func (s *Session) Login(u string, p string) error {
	// Clear cookie and modhash before sending request
	s.Cookie = ""
	s.modhash = ""

	return s.authenticate(u, p)
}

func (s *Session) authenticate(u string, p string) error {
	v := url.Values{"api_type": {"json"}}
	v.Set("user", u)
	v.Set("passwd", p)

	resp, err := s.post(buildURL(apiLogin, true), v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
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
	s.Cookie = fmt.Sprintf("%s=%s", strCookie, reply.Cookie)
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
		v = url.Values{}
	}
	req, err := http.NewRequest("POST", u, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header = s.httpHeaders()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return s.do(req)
}

func (s *Session) do(req *http.Request) (*http.Response, error) {
	s.lock.Lock()
	resp, err := s.client.Do(req)
	s.lock.Unlock()
	if err == nil {
		// Atoi errors can safely be ignored as 0 will be returned on bad input
		// and we will not need to check that each header really exists
		s.RateLimit.Used, _ = strconv.Atoi(resp.Header.Get("X-Ratelimit-Used"))
		s.RateLimit.Reset, _ = strconv.Atoi(resp.Header.Get("X-Ratelimit-Reset"))
		s.RateLimit.Remaining, _ = strconv.Atoi(resp.Header.Get("X-Ratelimit-Remaining"))
	}
	return resp, err
}
