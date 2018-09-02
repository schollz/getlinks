# getlinks

[![travis](https://travis-ci.org/schollz/getlinks.svg?branch=master)](https://travis-ci.org/schollz/getlinks) 
[![go report card](https://goreportcard.com/badge/github.com/schollz/getlinks)](https://goreportcard.com/report/github.com/schollz/getlinks) 
[![coverage](https://img.shields.io/badge/coverage-87%25-brightgreen.svg)](https://gocover.io/github.com/schollz/getlinks)
[![godocs](https://godoc.org/github.com/schollz/getlinks?status.svg)](https://godoc.org/github.com/schollz/getlinks) 

A very simple way to get links from web page

## Install

```
go get -u github.com/schollz/getlinks
```

## Usage 

```golang
url := "https://en.wikipedia.org/w/index.php?title=Pauli_exclusion_principle&oldid=854810355"
resp, _ := http.Get(url)
htmlBytes, _ := ioutil.ReadAll(resp.Body)
resp.Body.Close()

// get all links
linkgetter, err := getlinks.New(url)
if err != nil  {
  panic(err)
}
links := linkgetter.GetLinks(htmlBytes)
assert.Equal(t, 470, len(links))
```


## Contributing

Pull requests are welcome. Feel free to...

- Revise documentation
- Add new features
- Fix bugs
- Suggest improvements

## License

MIT
