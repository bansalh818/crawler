// file: crawl.go
package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"crawler/parser"
	"crawler/types"
	"crawler/utils"
)

func main() {
	flag.Parse()
	args := flag.Args()
	siteURL := ""
	var depth int
	var err error
	var visitedURL types.VisitedURL

	// Fetch the URL and depth
	// default URL : monzo.com and depth : 1
	switch len(args) {
	case 1:
		siteURL = args[0]
		depth = 1
	case 2:
		siteURL = args[0]
		depth, err = strconv.Atoi(args[1])
		if err != nil {
			//log
			return
		}

	default:
		siteURL = "https://monzo.com"
		depth = 1
	}

	// Parse the site URL
	URL, err := utils.Parse(siteURL)
	if err != nil {
		// Log failed to parse site url
		return
	}

	//Keeping track of time taken by the progam.
	start := time.Now()

	site := types.Site{URL: URL}
	visitedURL = types.VisitedURL{List: make(map[string]struct{})}

	//
	types.Waitlist.Add(1)
	go parser.GetLinks(&site, &visitedURL)
	types.Waitlist.Wait()

	printPage(&site, 0)

	elapsed := time.Since(start)
	//Log the output
	fmt.Println(elapsed)
	fmt.Println(depth)
}

var count = 1

func printPage(page *types.Site, indent int) {
	count = count + 1
	//a := strings.Join([]string{strings.Repeat("    ", indent), (*page).URL.String()}, "")
	//fmt.Println(a)
	if len((*page).Links) > 0 {
		//d := strings.Join([]string{strings.Repeat("    ", indent+1), "Links:"}, "")
		//fmt.Println(d)
		for _, subpage := range (*page).Links {
			printPage(subpage, indent+2)
		}
	}
}
