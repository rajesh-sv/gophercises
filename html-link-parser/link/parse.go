package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func GetLinks(r io.Reader) ([]Link, error) {
	htmlRootNode, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	links := dfsLinks(htmlRootNode)
	return links, nil
}

func dfsLinks(node *html.Node) []Link {
	if node.Data == "a" {
		link := Link{
			Href: getHrefValue(node.Attr),
			Text: getAnchorText(node),
		}
		return []Link{link}
	}
	var links []Link

	for childNode := range node.ChildNodes() {
		links = append(links, dfsLinks(childNode)...)
	}

	return links
}

func getHrefValue(attrs []html.Attribute) string {
	for _, attr := range attrs {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}

func getAnchorText(node *html.Node) string {
	text := ""
	for childNode := range node.ChildNodes() {
		if childNode.Type == html.CommentNode {
			continue
		} else if childNode.Type == html.TextNode {
			text += childNode.Data
		} else {
			text += getAnchorText(childNode)
		}
	}
	return strings.TrimSpace(text)
}
