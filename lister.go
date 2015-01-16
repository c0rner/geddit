package rego

import (
	"encoding/json"
)

// The Lister interface wraps methods used to extract Thing items
// from a Reddit Listing class.
//
// None of the methods return errors, as the data contained has already
// been unmarshalled once. It is assumed that the data is syntactically
// correct JSON and errors can safely be ignored.
//
// Any malformed/unrecognised Thing item will be silently dropped.
type Lister interface {
	Comments() []Comment
	Items() []interface{}
	Links() []Link
}

// Listing represents the Reddit Listing class documented
// at https://github.com/reddit/reddit/wiki/JSON#listing.
type Listing struct {
	Data struct {
		After    string  `json:"after"`
		Before   string  `json:"before"`
		Children []Thing `json:"children"`
		Modhash  string  `json:"modhash"`
	} `json:"data"`
	Kind string `json:"kind"` // Always "Listing"
}

// Comments return a slice of Comment types
func (l *Listing) Comments() []Comment {
	var items []Comment
	for _, c := range l.Data.Children {
		if c.Kind == TypeComment {
			item, _ := unmarshalComment(c.Data)
			items = append(items, *item)
		}
	}
	return items
}

// Items return a slice of interface{} items for cases where
// the caller want to do type assertions.
func (l *Listing) Items() []interface{} {
	var items []interface{}
	for _, c := range l.Data.Children {
		switch c.Kind {
		case TypeComment:
			item, _ := unmarshalComment(c.Data)
			items = append(items, *item)
		case TypeLink:
			item, _ := unmarshalLink(c.Data)
			items = append(items, *item)
		}
	}
	return items
}

// Links return a slice of Link types
func (l *Listing) Links() []Link {
	var items []Link
	for _, c := range l.Data.Children {
		if c.Kind == TypeLink {
			item, _ := unmarshalLink(c.Data)
			items = append(items, *item)
		}
	}
	return items
}

func unmarshalComment(j json.RawMessage) (*Comment, error) {
	item := Comment{}
	err := json.Unmarshal(j, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func unmarshalLink(j json.RawMessage) (*Link, error) {
	item := Link{}
	err := json.Unmarshal(j, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
