package main

import (
	"strings"
	"golang.org/x/net/html"
	"fmt"
	"flag"
	"os"
)

const htm = `<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">A link to another page</a>
</body>
</html>`

type Link struct {
	Href string
	Text string
}

func extractLinks(n *html.Node) []Link {
	var Links []Link
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				txt := extractText(n)
				Links = append(Links, Link{a.Val, txt})
			}
		}
	}
	for m := n.FirstChild; m != nil; m = m.NextSibling {
		exLinks := extractLinks(m)
		Links = append(Links, exLinks...)
	}
	return Links
}

func extractText(n *html.Node) string {
	var text string
	if n.Type != html.ElementNode && n.Data != "a" && n.Type != html.CommentNode {
		text = n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c)
	}
	return strings.Trim(text, "\n")
}

func main() {
	filename := flag.String("file", "exam.html", "The HTML file to parse links from")
	flag.Parse()
	n, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	doc, _ := html.Parse(n)
	Links := extractLinks(doc)
	for _ , i := range Links {
		fmt.Println("href = ", i.Href)
		fmt.Println("Text = ", i.Text)
	}
}
