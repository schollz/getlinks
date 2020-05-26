package getlinks

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"path"
	"strings"

	"golang.org/x/net/html"
)

type config struct {
	sameDomain, disallowQuery, disallowFragment bool
}

// Option is the type all options need to adhere to
type Option func(c *config)

// OptionSameDomain determines whether to get from same domain
func OptionSameDomain(sameDomain bool) Option {
	return func(c *config) {
		c.sameDomain = sameDomain
	}
}

// OptionDisallowQuery will disallow queries
func OptionDisallowQuery(disallowQuery bool) Option {
	return func(c *config) {
		c.disallowQuery = disallowQuery
	}
}

// OptionDisallowFragment will disallow fragments
func OptionDisallowFragment(disallowFragment bool) Option {
	return func(c *config) {
		c.disallowFragment = disallowFragment
	}
}

// GetLinks returns all unique links, in order from the bytes, while setting
// urls relative to the provided URL in the case there is no hostname.
func GetLinks(htmlBytes []byte, urlString string, options ...Option) (linkList []string, err error) {
	c := new(config)
	for _, o := range options {
		o(c)
	}

	parent, err := url.Parse(urlString)
	if err != nil {
		return
	}

	links := make(map[string]int)
	z := html.NewTokenizer(ioutil.NopCloser(bytes.NewBuffer(htmlBytes)))
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.StartTagToken {
			token := z.Token()
			if "a" != token.Data {
				continue
			}
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					parse1, err := url.Parse(attr.Val)
					if err != nil {
						break
					}
					var rel *url.URL
					if parse1.Hostname() == "" {
						if strings.HasPrefix(parse1.Path, "/") {
							rel, err = parent.Parse(attr.Val)
						} else {
							rel, err = parent.Parse(path.Join(parent.Path, attr.Val))
						}
					} else {
						rel = parse1
					}
					if err != nil {
						break
					}

					if rel.Scheme == "" {
						rel.Scheme = "https"
					}
					cleanedUrl := fmt.Sprintf("%s://%s%s", rel.Scheme, rel.Host, rel.Path)
					if strings.HasSuffix(attr.Val, "/") && !strings.HasSuffix(cleanedUrl, "/") {
						cleanedUrl += "/"
					}
					if rel.RawQuery != "" && !c.disallowQuery {
						cleanedUrl += "?" + rel.RawQuery
					}
					if rel.Fragment != "" && !c.disallowFragment {
						cleanedUrl += "#" + rel.Fragment
					}
					if c.sameDomain && rel.Hostname() != parent.Hostname() {
						break
					}
					if _, ok := links[cleanedUrl]; !ok {
						links[cleanedUrl] = len(links)
					}
				}
			}
		}
	}

	if len(links) == 0 {
		return
	}

	linkList = make([]string, len(links))
	for link := range links {
		linkList[links[link]] = link
	}
	return
}
