# Geddit
Reddit API bindings for Go. Things are currently in a state of flux as I try to figure out what parts of the API I need and howto implement them.

## Basic usage
```
session := geddit.NewSession("GedditBot/1.0")
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
```
