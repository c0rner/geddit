# Rego
Reddit API bindings for Go. Things are currently in a state of flux as I try to figure out what parts of the API I need and how to implement them.

## Documentation
http://godoc.org/github.com/c0rner/rego

### Example: Listing 5 newest /r/worldnews link posts
```go
session := rego.NewSession("RegoBot/1.0")
page := session.Listing("r/worldnews/new")
page.SetLimit(5)
list, err := page.Next()
if err != nil {
        log.Print(err)
}
for _, l := range list.Links() {
        if l.Selfpost {
                continue
        }
        fmt.Printf("Title: %30s, Url: %s\n", l.Title, l.Url)
}
```

## TODO
- [ ] Unauthenticated sessions should use http by default to take advantage of Reddit caches
  - Hitting Reddit caches are 'free' requests and do not count against rate limits
  - This will be configurable for cases where TLS is desired
