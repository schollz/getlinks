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

// GetLinks returns all the links for the given bytes, using the urlString to fix relative domain names. If sameDomain is true, then it only returns links on the same domain.
func GetLinks(htmlBytes []byte, urlString string, sameDomain bool, allowQuery bool, allowFragment bool) (linkList []string, err error) {
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

					cleanedUrl := fmt.Sprintf("%s://%s%s", rel.Scheme, rel.Host, rel.Path)
					if rel.RawQuery != "" && allowQuery {
						cleanedUrl += "?" + rel.RawQuery
					}
					if rel.Fragment != "" && allowFragment {
						cleanedUrl += "#" + rel.Fragment
					}
					if sameDomain && rel.Hostname() != parent.Hostname() {
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
