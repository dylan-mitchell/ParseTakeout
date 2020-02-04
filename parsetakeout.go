package ParseTakeout

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/net/html"
)

type Result struct {
	Title  string `json:"title"`
	Action string `json:"action"`
	Item   string `json:"item"`
	Date   string `json:"date"`
}

var actions = []string{"Listened to", "Searched for", "Visited", "Used", "Viewed", "Watched"}

func (r Result) String() string {
	s := fmt.Sprintf(`*****
Title: %s
Action: %s
Item: %s
Date: %s
*****`, r.Title, r.Action, r.Item, r.Date)
	return s
}

func ReadHtml(filePath string) (string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func body(doc *html.Node) (*html.Node, error) {
	var body *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			body = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if body != nil {
		return body, nil
	}
	return nil, errors.New("Missing <body> in the node tree")
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func ParseHTML(filePath string) ([]Result, error) {
	results := []Result{}

	s, err := ReadHtml(filePath)
	if err != nil {
		return nil, err
	}
	//Strip down to the body
	doc, _ := html.Parse(strings.NewReader(s))
	bn, err := body(doc)
	if err != nil {
		return nil, err
	}
	body := renderNode(bn)

	z := html.NewTokenizer(strings.NewReader(body))

	var res Result
	count := 0

	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return results, nil
		case tt == html.TextToken:
			t := z.Token()
			// TODO: I think this can probably be transitioned to some kind of middleware-esque flow
			//This is 🤮
			switch count {
			case 0:
				res.Title = t.Data
				count++
			case 1:
				exactMatch := false
				data := strings.TrimSpace(t.Data)
				for _, action := range actions {
					if data == action {
						//If exact match then next node is item
						res.Action = data
						exactMatch = true
					}
					if !exactMatch {
						if strings.Contains(data, action) {
							res.Action = action
							res.Item = strings.TrimSpace(strings.TrimPrefix(data, action))
						}
					}
				}
				if exactMatch {
					count++
				} else {
					count = 3
				}
			case 2:
				res.Item = t.Data
				count++
			case 3:
				res.Date = t.Data
				results = append(results, res)
				// Clear results
				res.Title = ""
				res.Action = ""
				res.Item = ""
				res.Date = ""
				count++
			}

		case tt == html.StartTagToken:
			if count == 4 {
				//Check for new item denoted by class="mdl-typography--title"
				_, val, more := z.TagAttr()
				if string(val) == "mdl-typography--title" {
					// Next item will be title
					count = 0
				}
				for more {
					_, val, more = z.TagAttr()
					if string(val) == "mdl-typography--title" {
						// Next item will be title
						count = 0
					}
				}
			}
		}

	}
}
