package main

import (
	"fmt"
	"log"
	"time"

	"github.com/c0rner/rego"
)

func main() {
	session := rego.NewSession("RegoBot/1.0")
	page := session.Listing("r/worldnews/new")
	page.SetLimit(10)
	for i := 0; i < 2; i++ {
		list, err := page.Previous()
		if err != nil {
			log.Printf("Error: %s", err)
		}
		now := time.Now()
		for _, l := range list.Links() {
			age := now.Sub(l.Created.Time())
			fmt.Printf("Created %d minutes ago, Name: %s, Title: %s\n", int(age.Minutes()), l.Name, l.Title)
		}
	}
}
