package getlinks

import (
	"bytes"
	"io/ioutil"
	"net/url"

	"golang.org/x/net/html"
)

type getLinks struct {
	parent         *url.URL
	sameDomainOnly bool
}

// New creates a new instance of the get links
func New(urlString string, sameDomain ...bool) (gl *getLinks, err error) {
	gl = new(getLinks)
	gl.parent, err = url.Parse(urlString)
	if len(sameDomain) > 0 {
		gl.sameDomainOnly = sameDomain[0]
	}
	return
}

// GetLinks returns all the links for the given bytes
func (gl *getLinks) GetLinks(htmlBytes []byte) (linkList []string) {
	links := make(map[string]struct{})
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
					rel, err := gl.parent.Parse(attr.Val)
					if err != nil {
						break
					}
					if gl.sameDomainOnly && rel.Hostname() != gl.parent.Hostname() {
						break
					}
					links[rel.String()] = struct{}{}
				}
			}
		}
	}

	linkList = make([]string, len(links))
	i := 0
	for link := range links {
		linkList[i] = link
		i++
	}
	return
}
