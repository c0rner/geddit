# Rego
Reddit API bindings for Go. Things are currently in a state of flux as I try to figure out what parts of the API I need and how to implement them.

## Documentation
http://godoc.org/github.com/c0rner/rego

### Example: Listing 5 newest /r/worldnews link posts
```go
session := rego.NewSession("RegoBot/1.0")
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
        fmt.Printf("Title: %30s, Url: %s\n", l.Title, l.Url)
}
```
