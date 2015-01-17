# Rego ![Gopher](https://github.com/c0rner/c0rner.github.io/blob/master/images/redditgopher_small.png)
Reddit API bindings for Go. The primary focus of this package is to be a tool for writing bots and not a full API implementation.

## Documentation
[![GoDoc](https://godoc.org/github.com/c0rner/rego?status.svg)](https://godoc.org/github.com/c0rner/rego)

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
