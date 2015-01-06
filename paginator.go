package geddit

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	MaxLimit = 100 // Listing max number of returned Link items
)

// Paginator represents a Listing endpoint
type Paginator struct {
	s     *Session
	url   string
	name  string // Fullname of reference Thing
	count int    // Current item offset *NOT USED*
	limit int    // Limit of items returned
}

// Next returns a new set of links directly following a previous request.
func (p *Paginator) Next() ([]Link, error) {
	resp, err := p.list(true)
	if err != nil {
		return nil, err
	}
	if len(resp) > 0 {
		p.name = resp[len(resp)-1].Name
	}
	return resp, err
}

// Previous returns a new set of links directly preceeding a previous request.
func (p *Paginator) Previous() ([]Link, error) {
	resp, err := p.list(false)
	if err != nil {
		return nil, err
	}
	if len(resp) > 0 {
		p.name = resp[0].Name
	}
	return resp, err
}

// SetLimit sets the max number of links returned from calls to Previous and Next
func (p *Paginator) SetLimit(l int) {
	if l < 0 {
		l = 0
	}
	if l > MaxLimit {
		l = MaxLimit
	}
	p.limit = l
}

func (p *Paginator) list(after bool) ([]Link, error) {
	v := p.values(after)
	resp, err := p.s.get(p.url, v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	listing := Listing{}
	err = json.NewDecoder(resp.Body).Decode(&listing)
	if err != nil {
		return nil, err
	}

	if listing.Kind != TypeListing {
		// TODO.. handle error
	}
	var links []Link
	for _, c := range listing.Data.Children {
		link := Link{}
		err = json.Unmarshal(c.Data, &link)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

func (p *Paginator) values(after bool) url.Values {
	v := url.Values{}
	if p.name != "" {
		if after {
			v.Set("after", p.name)
		} else {
			v.Set("before", p.name)
		}
	}
	if p.count > 0 {
		v.Set("count", fmt.Sprintf("%d", p.count))
	}
	if p.limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", p.limit))
	}

	return v
}
