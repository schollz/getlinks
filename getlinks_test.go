package getlinks

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkGetLinks(b *testing.B) {
	url := "https://en.wikipedia.org/w/index.php?title=Pauli_exclusion_principle&oldid=854810355"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetLinks(bodyBytes, url, false, false, false)
	}
}

func TestGetLinks(t *testing.T) {
	url := "https://en.wikipedia.org/w/index.php?title=Pauli_exclusion_principle&oldid=854810355"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// get all links
	links, err := GetLinks(bodyBytes, url, false, true, true)
	assert.Nil(t, err)
	assert.Equal(t, 470, len(links))

	// get all links on the same domain
	links, err = GetLinks(bodyBytes, url, true, true, true)
	assert.Equal(t, 378, len(links))
}

func TestLinksHash(t *testing.T) {
	urlString := "https://schollz.github.io/watercolor"
	resp, err := http.Get(urlString)
	if err != nil {
		log.Fatal(err)
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// get all links
	links, err := GetLinks(bodyBytes, urlString, true, true, false)
	assert.Nil(t, err)
	fmt.Println(links)
}

func TestGetLinksAll(t *testing.T) {
	b := []byte(`
<html>
<head>
</head>
<body>
<a href="./test.html">test1</a>
<a href="/test.html">test1</a>
<a href="test.html">test1</a>
<a href="./test/">test1</a>
<a href="https://somethingelse.com">test1</a>
<a href="https://test.com/test2/test?q=hi#world">test1</a>
</body>
</html>`)
	// get all links
	links, err := GetLinks(b, "https://test.com/test", false, true, true)
	assert.Nil(t, err)
	assert.Equal(t,
		`https://test.com/test/test.html|https://test.com/test.html|https://test.com/test/test|https://somethingelse.com|https://test.com/test2/test?q=hi#world`,
		strings.Join(links, "|"),
	)

	links, err = GetLinks(b, "https://test.com/test", true, false, false)
	assert.Nil(t, err)
	assert.Equal(t,
		`https://test.com/test/test.html|https://test.com/test.html|https://test.com/test/test|https://test.com/test2/test`,
		strings.Join(links, "|"),
	)
}
