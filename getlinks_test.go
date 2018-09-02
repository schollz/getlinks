package getlinks

import (
	"io/ioutil"
	"log"
	"net/http"
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
		GetLinks(bodyBytes, url, false)
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
	links, err := GetLinks(bodyBytes, url)
	assert.Nil(t, err)
	assert.Equal(t, 470, len(links))

	// get all links on the same domain
	links, err = GetLinks(bodyBytes, url, true)
	assert.Equal(t, 378, len(links))
}
