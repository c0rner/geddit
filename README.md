# Rego
Reddit API bindings for Go. Things are currently in a state of flux as I try to figure out what parts of the API I need and how to implement them.

## Documentation
http://godoc.org/github.com/c0rner/rego

### Example: Listing 5 newest /r/worldnews posts
```go
session := rego.NewSession("RegoBot/1.0")
listing := session.Listing("worldnews/new")
listing.SetLimit(5)
links, err := listing.Next()
if err != nil {
        log.Print(err)
}
for _, l := range links {
        u, _ := l.Created.UTC.Float64()
        t := time.Unix(int64(u), 0)
        fmt.Printf("Created: %s, Title: %30s", t, l.Title)
        if !l.Selfpost {
                fmt.Printf("Url: %s", l.Url)
        }
        fmt.Printf("\n")
}
```
