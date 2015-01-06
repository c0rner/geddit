package geddit

import (
	"fmt"
	"log"
	"time"
)

func ExampleAPIError() {
	session := NewSession("GedditBot/1.0")
retry:
	_, err := session.Comment("t3_xxxxx", "This is a demo")
	if err != nil {
		if apierr, ok := err.(APIError); ok {
			if apierr.IsRatelimited() {
				fmt.Printf("We are being ratelimited for %.0f minutes\n", apierr.Duration().Minutes())
				time.Sleep(apierr.Duration())
				goto retry
			}
		}
	}
}

func ExampleSession_Listing() {
	session := NewSession("GedditBot/1.0")
	listing := session.Listing("worldnews/new")
	listing.SetLimit(5)
	links, err := listing.Next()
	if err != nil {
		log.Print(err)
	}
	for _, l := range links {
		if l.Selfpost {
			continue
		}
		fmt.Printf("Title: %30s Url: %s\n", l.Title, l.Url)
	}
}

func ExampleSession_Login() {
	session := NewSession("GedditBot/1.0")
	auth := Authconfig{
		User:     "username",
		Password: "password",
	}
	err := session.Login(&auth)
	if err != nil {
		log.Fatal(err)
	}
}
