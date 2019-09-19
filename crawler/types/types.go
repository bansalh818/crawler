package types

import (
	"net/url"
	"sync"
)

type Site struct {
	URL   *url.URL
	Links []*Site
}

type VisitedURL struct {
	List map[string]struct{}
}

var Waitlist sync.WaitGroup
