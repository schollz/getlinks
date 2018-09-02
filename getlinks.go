package getlinks

import (
	"bytes"
	"io/ioutil"
	"net/url"

	"golang.org/x/net/html"
)

// GetLinks returns all the links for the given bytes, using the urlString to fix relative domain names. If sameDomain is true, then it only returns links on the same domain.
func GetLinks(htmlBytes []byte, urlString string, sameDomain ...bool) (linkList []string, err error) {
	parent, err := url.Parse(urlString)
	if err != nil {
		return
	}

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
					rel, err := parent.Parse(attr.Val)
					if err != nil {
						break
					}
					if len(sameDomain) > 0 && sameDomain[0] && rel.Hostname() != parent.Hostname() {
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
