package link_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/rajesh-sv/gophercises/html-link-parser/link"
)

func TestGetLinks(t *testing.T) {
	htmlFiles := []string{"ex1.html", "ex2.html", "ex3.html", "ex4.html"}
	expectedLinksSet := [][]link.Link{
		{
			link.Link{"/other-page", "A link to another page"},
		},
		{
			link.Link{"https://www.twitter.com/joncalhoun", "Check me out on twitter"},
			link.Link{"https://github.com/gophercises", "Gophercises is on Github!"},
		},
		{
			link.Link{"#", "Login"},
			link.Link{"/lost", "Lost? Need help?"},
			link.Link{"https://twitter.com/marcusolsson", "@marcusolsson"},
		},
		{
			link.Link{"/dog-cat", "dog cat"},
		},
	}

	for i, htmlFile := range htmlFiles {
		t.Run(fmt.Sprintf("testing html file: %s", htmlFile), func(t *testing.T) {
			f, err := os.Open(htmlFile)
			if err != nil {
				t.Error(err)
			}
			defer f.Close()

			links, err := link.GetLinks(f)
			if err != nil {
				t.Error(err)
			}

			for j, link := range links {
				if link != expectedLinksSet[i][j] {
					fmt.Println(link)
					t.Fail()
				}
			}
		})
	}
}
