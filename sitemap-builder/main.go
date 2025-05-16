package main

import (
	"encoding/xml"
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/rajesh-sv/gophercises/html-link-parser/link"
)

type sitemapUrl struct {
	XMLName string `xml:"url"`
	Loc     string `xml:"loc"`
}

func main() {
	websiteDomain := flag.String("domain", "https://gophercises.com", "website domain to build a sitemap for")
	flag.Parse()

	queue := []string{*websiteDomain}
	visited := map[string]bool{*websiteDomain: true}

	for len(queue) > 0 {
		length := len(queue)
		for i := 0; i < length; i++ {
			url := queue[0]
			queue = queue[1:]
			newUrls := getNewUrls(url)

			for _, newUrl := range newUrls {
				if newUrl[0] == '/' {
					newUrl = *websiteDomain + newUrl
				}
				if newUrl[len(newUrl)-1] == '/' {
					newUrl = newUrl[:len(newUrl)-1]
				}
				if strings.Index(newUrl, *websiteDomain) == 0 && !visited[newUrl] {
					visited[newUrl] = true
					queue = append(queue, newUrl)
				}
			}
		}
	}

	f, err := os.OpenFile("sitemap.xml", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	xmlHeader := []byte(`<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	if _, err := f.Write(xmlHeader); err != nil {
		panic(err)
	}

	var urlSet []sitemapUrl
	for url := range visited {
		urlSet = append(urlSet, sitemapUrl{Loc: url})
	}
	xmlEncoder := xml.NewEncoder(f)
	if err := xmlEncoder.Encode(urlSet); err != nil {
		panic(err)
	}

	xmlFooter := []byte(`</urlset>`)
	if _, err := f.Write(xmlFooter); err != nil {
		panic(err)
	}
}

func getNewUrls(url string) []string {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	links, err := link.GetLinks(res.Body)
	if err != nil {
		panic(err)
	}
	res.Body.Close()

	var newUrls []string
	for _, link := range links {
		newUrls = append(newUrls, link.Href)
	}

	return newUrls
}
