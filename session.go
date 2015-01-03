package geddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Session is an active Reddit session that initially is unauthenticated. An authenticated
// session can be set up using Session.Login or by updating Session.Cookie manually.
type Session struct {
	client    *http.Client
	Cookie    string // Session cookie (empty if not logged in)
	modhash   string
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
func (s *Session) Comment(p string, t string) (string, error) {
	v := url.Values{"api_type": {"json"}}
	v.Set("thing_id", p)
	v.Set("text", t)

	resp, err := s.post(BuildURL(APIComment, true), &v)
	/*
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		fmt.Printf("Data: %#v\n", buf.String())
	*/

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	r, err := getJSON(resp.Body)

	type container struct {
		things []Thing
	}

	things := container{}
	err = json.Unmarshal(r.JSON.Data, &things)
	if err != nil {
		return "", err
	}

	fmt.Printf("Data: %#v\n", things)
	th := Thing{}
	for _, thing := range things.things {
		err = json.Unmarshal(thing.Data, &th)
		fmt.Printf("Id: %s\n", th.ID)
	}

	if err != nil {
		return "", err
	}

	return "", nil
}

// Login authenticates the current session
func (s *Session) Login(ac *Authconfig) error {
	if ac == nil || len(ac.User) == 0 || len(ac.Password) == 0 {
		return ErrNoCredentials
	}
	v := url.Values{"api_type": {"json"}}
	v.Set("user", ac.User)
	v.Set("passwd", ac.Password)

	s.Cookie = ""
	s.modhash = ""

	resp, err := s.post(BuildURL(APILogin, true), &v)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	r, err := getJSON(resp.Body)
	if err != nil {
		return err
	}

	type loginReply struct {
		NeedHttps bool   `json:"need_https"`
		Modhash   string `json:"modhash"`
		Cookie    string `json:"cookie"`
	}
	reply := loginReply{}
	err = json.Unmarshal(r.JSON.Data, &reply)
	if err != nil {
		return err
	}

	// Get the session cookie and modhash from response
	/*
		for _, c := range resp.Cookies() {
			fmt.Printf("%s\n", c)
			if c.Name == StrCookie {
				s.Cookie = c.String()
			}
		}
	*/
	s.Cookie = fmt.Sprintf("%s=%s", StrCookie, reply.Cookie)
	s.modhash = reply.Modhash

	fmt.Printf("Cookie: %s\nModhash: %s\n", s.Cookie, s.modhash)

	return nil
}

func (s *Session) get(u string, v *url.Values) (*http.Response, error) {
	var values string
	if v != nil {
		values = fmt.Sprintf("?%s", v.Encode())
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", u, values), nil)
	if err != nil {
		return nil, err
	}

	return s.doRequest(req)
}

func (s *Session) post(u string, v *url.Values) (*http.Response, error) {
	if v == nil {
		return nil, ErrNoValues
	}
	req, err := http.NewRequest("POST", u, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return s.doRequest(req)
}

func (s *Session) doRequest(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", s.useragent)
	if len(s.Cookie) != 0 {
		r.Header.Set("Cookie", s.Cookie)
	}
	if len(s.modhash) != 0 {
		r.Header.Set("X-Modhash", s.modhash)
	}
	c := &http.Client{}
	return c.Do(r)
}

// getJSON is a convenience function used by all API methods using api_type=json
func getJSON(rc io.ReadCloser) (*jsonAPIReply, error) {
	var r jsonAPIReply
	err := json.NewDecoder(rc).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
