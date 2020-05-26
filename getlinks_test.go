package getlinks

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func ExampleGetLinks() {
	urlString := "https://en.wikipedia.org/w/index.php?title=Pauli_exclusion_principle&oldid=854810355"
	resp, _ := http.Get(urlString)
	htmlBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// get all links
	links, _ := GetLinks(htmlBytes, urlString)
	fmt.Println(strings.Join(links[:3], ", "))
	// Output: https://en.wikipedia.org/wiki/Help:Page_history, https://en.wikipedia.org/wiki/User:Turgidson, https://en.wikipedia.org/wiki/User_talk:Turgidson
}

func TestGetLinks(t *testing.T) {
	var flagtests = []struct {
		u    string
		href string
		out  string
	}{
		{"https://a.com", "/1", "https://a.com/1"},
		{"https://a.com", "/1?q=1#2", "https://a.com/1?q=1#2"},
		{"https://a.com/test/", "/1", "https://a.com/1"},
		{"https://a.com/2.html", "/1", "https://a.com/1"},
		{"https://a.com/test/something", "../1", "https://a.com/test/1"},
		{"https://a.com/test/something", "../1/", "https://a.com/test/1/"},
		{"https://a.com/test", "1", "https://a.com/test/1"},
		{"https://a.com/test/something", "https://a.com/2", "https://a.com/2"},
		{"https://a.com/2", "https://somethingelse.com/1", "https://somethingelse.com/1"},
	}

	for _, tt := range flagtests {
		links, err := GetLinks([]byte(fmt.Sprintf("<a href='%s'>test</a>", tt.href)), tt.u)
		if err != nil {
			t.Errorf("%+v: %s", tt, err)
		}
		if links[0] != tt.out {
			t.Errorf("%+v: actually %s", tt, links[0])
		}
	}
}
func TestSameOrigin(t *testing.T) {
	var flagtests = []struct {
		u    string
		href string
	}{
		{"https://a.com/2", "https://somethingelse.com/1"},
	}

	for _, tt := range flagtests {
		links, err := GetLinks([]byte(fmt.Sprintf("<a href='%s'>test</a>", tt.href)), tt.u, OptionSameDomain(true))
		if err != nil {
			t.Errorf("%+v: %s", tt, err)
		}
		if len(links) > 0 {
			t.Errorf("got links when should have not: %+v", links)
		}
	}
}

func TestGetLinksNoFragment(t *testing.T) {
	var flagtests = []struct {
		u    string
		href string
		out  string
	}{
		{"https://a.com", "/1?q=1#2", "https://a.com/1?q=1"},
	}

	for _, tt := range flagtests {
		links, err := GetLinks([]byte(fmt.Sprintf("<a href='%s'>test</a>", tt.href)), tt.u, OptionDisallowFragment(true))
		if err != nil {
			t.Errorf("%+v: %s", tt, err)
		}
		if links[0] != tt.out {
			t.Errorf("%+v: actually %s", tt, links[0])
		}
	}

}

func TestGetLinksNoQuery(t *testing.T) {
	var flagtests = []struct {
		u    string
		href string
		out  string
	}{
		{"https://a.com", "/1?q=1#2", "https://a.com/1#2"},
	}

	for _, tt := range flagtests {
		links, err := GetLinks([]byte(fmt.Sprintf("<a href='%s'>test</a>", tt.href)), tt.u, OptionDisallowQuery(true))
		if err != nil {
			t.Errorf("%+v: %s", tt, err)
		}
		if links[0] != tt.out {
			t.Errorf("%+v: actually %s", tt, links[0])
		}
	}
}

func TestGetLinksNoQueryNoFragment(t *testing.T) {
	var flagtests = []struct {
		u    string
		href string
		out  string
	}{
		{"https://a.com", "/1/?q=1#2", "https://a.com/1/"},
	}

	for _, tt := range flagtests {
		links, err := GetLinks([]byte(fmt.Sprintf("<a href='%s'>test</a>", tt.href)), tt.u, OptionDisallowQuery(true), OptionDisallowFragment(true))
		if err != nil {
			t.Errorf("%+v: %s", tt, err)
		}
		if links[0] != tt.out {
			t.Errorf("%+v: actually %s", tt, links[0])
		}
	}
}
