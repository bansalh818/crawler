// file parser.go
package parser

import (
	"crawler/types"
	"io"
	"net/http"
	"strings"
	"sync"

	"crawler/utils"

	"golang.org/x/net/html"
	//"github.com/sirupsen/logrus"
)

var mutex sync.Mutex

func GetLinks(target *types.Site, visitedURL *types.VisitedURL) {
	defer types.Waitlist.Done()

	// Get page html
	resp, err := http.Get((*target).URL.String())
	if err != nil {
		//Log error
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Log
		return
	}

	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		// Log
		return
	}

	// Create channel to add links to Site
	links := make(chan *types.Site)
	var linkswg sync.WaitGroup

	linkswg.Add(1)
	defer linkswg.Done()

	types.Waitlist.Add(1)
	// Close link channels when parsing finishes
	go func() {
		defer types.Waitlist.Done()
		linkswg.Wait()
		close(links)
	}()

	types.Waitlist.Add(1)
	// Add link to site
	go func() {
		defer types.Waitlist.Done()
		for link := range links {
			(*target).Links = append((*target).Links, link)
		}
	}()

	mutex.Lock()
	seenRefs := make(map[string]struct{})
	mutex.Unlock()

	tokenizer := html.NewTokenizer(resp.Body)
	for {
		// get the next token type
		tokenType := tokenizer.Next()

		// if it's an error token, we either reached
		// the end of the file, or the HTML was malformed
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				//end of the file, break out of the loop
				break
			}
			// Log error
		}

		if tokenType == html.StartTagToken {
			// get the token
			token := tokenizer.Token()

			switch token.DataAtom.String() {
			case "a", "link": //link tags
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						_, ok := seenRefs[attr.Val]
						if !ok {
							seenRefs[attr.Val] = struct{}{}
							linkswg.Add(1)
							go ParseLink(attr.Val, target, links, visitedURL, &linkswg)
						}
					}
				}
			}
		}
	}
	return
}

func ParseLink(href string, target *types.Site, result chan *types.Site, visitedURL *types.VisitedURL, linkswg *sync.WaitGroup) {
	defer (*linkswg).Done()

	relURL, err := utils.Parse(href)
	if err != nil {
		// Log error
		return
	}

	newURL := (*target).URL.ResolveReference(relURL)
	if newURL.Host != target.URL.Host {
		return
	}

	// Remove the fragment
	newURL.Fragment = ""

	// Check whether we have seen the URL before
	mutex.Lock()
	_, ok := visitedURL.List[newURL.String()]
	if ok {
		mutex.Unlock()
		return
	} else {
		// Mark the URL as visited
		visitedURL.List[newURL.String()] = struct{}{}
		mutex.Unlock()
	}

	// Add new url as site
	newPage := types.Site{URL: newURL}
	types.Waitlist.Add(1)
	go GetLinks(&newPage, visitedURL)

	result <- &newPage
	return
}
