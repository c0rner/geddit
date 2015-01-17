package rego

import (
	"fmt"
	"time"
)

// APIError represents a Reddit API error
type APIError struct {
	id   string
	desc string
	wait time.Time
}

func newAPIError(e *jsonAPIReply) error {
	if len(e.Errors) == 0 {
		return nil
	}
	// FIXME: This will explode if Errors is not a two-dimensional array
	err := APIError{
		id:   e.Errors[0][0],
		desc: e.Errors[0][1],
		wait: time.Now(),
	}
	if e.Ratelimit > 0 {
		err.wait = err.wait.Add(time.Duration(e.Ratelimit) * time.Second)
	}
	return err
}

// Error returns a descriptive string of the error
func (e APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.id, e.desc)
}

// IsRatelimited returns true if a ratelimit is in effect for the error
func (e APIError) IsRatelimited() bool {
	return e.wait.After(time.Now())
}

// Duration returns the time remaining of active ratelimit for the error
func (e APIError) Duration() time.Duration {
	return e.wait.Sub(time.Now())
}
