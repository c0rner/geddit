package rego

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

const (
	// MaxLimit is the upper maximum limit for Session.Listing() requests
	MaxLimit = 100
)

// Paginator is the interface that wraps methods for pagination of the Listing type
type Paginator interface {
	Next() (Lister, error)
	Previous() (Lister, error)
	SetLimit(int)
}

type Page struct {
	s      *Session
	url    string
	after  string // Fullname of reference Thing
	before string // Fullname of reference Thing
	limit  int    // Limit of items returned
}

// Next returns a set of Thing items resulting from the requested API call.
//
// Consecutive calls will return subsequent items indefinitely.
func (p *Page) Next() (*Listing, error) {
	v := p.values()
	if len(p.before) > 0 {
		v.Set("before", p.before)
	}
	resp, err := p.list(v)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// Previous returns a set of Thing items resulting from the requested API call.
//
// Consecutive calls will return preceding items until source is exhausted.
func (p *Page) Previous() (*Listing, error) {
	v := p.values()
	if len(p.after) > 0 {
		v.Set("after", p.after)
	}
	resp, err := p.list(v)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// SetLimit sets the max number of links returned from calls to Previous and Next
func (p *Page) SetLimit(limit int) {
	if limit < 0 {
		limit = 0
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	p.limit = limit
}

func (p *Page) list(v url.Values) (*Listing, error) {
	resp, err := p.s.get(p.url, v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	list := Listing{}
	err = json.NewDecoder(resp.Body).Decode(&list)
	if err != nil {
		return nil, err
	}

	if list.Kind != TypeListing {
		// TODO.. handle error
		log.Panic("Bricks!")
	}

	if len(list.Data.Children) > 0 {
		item := struct {
			Name string
		}{}
		err = json.Unmarshal(list.Data.Children[0].Data, &item)
		p.before = item.Name
		err = json.Unmarshal(list.Data.Children[len(list.Data.Children)-1].Data, &item)
		p.after = item.Name
	}

	return &list, nil
}

func (p *Page) values() url.Values {
	v := url.Values{}
	if p.limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", p.limit))
	}
	return v
}
