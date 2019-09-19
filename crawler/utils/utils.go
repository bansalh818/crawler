package utils

import "net/url"

// Gethost returns the hostname of the url
func Gethost(rawurl string) (string, error) {
	targetURL, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	} else {
		return targetURL.Host, nil
	}
}

func Parse(rawurl string) (*url.URL, error) {
	targetURL, err := url.Parse(rawurl)
	return targetURL, err
}

// ResolveUrl resolves the rawurl
func ResolveUrl(rawURL, siteURL string) (string, error) {
	targetURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	tarURL, err := url.Parse(siteURL)
	if err != nil {
		return "", err
	}
	return tarURL.ResolveReference(targetURL).String(), nil
}
