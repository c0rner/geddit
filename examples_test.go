package geddit

func ExampleSession_Listing() {
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

func ExampleSession_Login_auth() {
	session := geddit.NewSession("GedditBot/1.0")
	auth := geddit.Authconfig{
		User:     "username",
		Password: "password",
	}
	err = session.Login(&auth)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleSession_Login_cookie() {
	session := geddit.NewSession("GedditBot/1.0")
	session.Cookie = "reddit_session=4826...d8eca"
}
