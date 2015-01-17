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
	page.SetLimit(5)
	for i := 1; i <= 2; i++ {
		fmt.Printf("--- PAGE %.2d ---\n", i)
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

	fmt.Printf("\n--- USER LISTING ---\n")
	account, err := session.User("wil")
	if err != nil {
		log.Printf("Error: %s", err)
	} else {
		fmt.Printf("ID: %s\n", account.ID)
		fmt.Printf("User: %s\n", account.Name)
		fmt.Printf("Is a mod: %t\n", account.IsMod)
		fmt.Printf("Comment karma: %d\n", account.CommentKarma)
	}
}
