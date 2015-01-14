package main

import (
	"fmt"
	"log"
	"time"

	"github.com/c0rner/rego"
)

func main() {
	session := rego.NewSession("RegoBot/1.0")
	listing := session.Listing("worldnews/new")
	listing.SetLimit(10)
	links, err := listing.Next()
	if err != nil {
		log.Printf("Error: %s", err)
	}
	now := time.Now()
	for _, l := range links {
		age := now.Sub(l.Created.Time())
		fmt.Printf("Created %d minutes ago, Title: %s\n", int(age.Minutes()), l.Title)
	}
}
