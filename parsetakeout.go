package ParseTakeout

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"golang.org/x/net/html"
)

type Result struct {
	Title   string `json:"title"`
	Action  string `json:"action"`
	Item    string `json:"item"`
	Channel string `json:"channel"`
	Date    string `json:"date"`
}

var actions = []string{"Listened to", "Searched for", "Visited", "Used", "Viewed", "Watched"}

func (r Result) String() string {
	if len(r.Channel) == 0 {
		s := fmt.Sprintf(`*****
	Title: %s
	Action: %s
	Item: %s
	Date: %s
	*****`, r.Title, r.Action, r.Item, r.Date)
		return s
	} else {
		s := fmt.Sprintf(`*****
	Title: %s
	Action: %s
	Item: %s
	Channel: %s
	Date: %s
	*****`, r.Title, r.Action, r.Item, r.Channel, r.Date)
		return s
	}
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
	nextItem := "title"
	getChannel := false

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
			switch nextItem {
			case "title":
				res.Title = t.Data
				nextItem = "action"
			case "action":
				exactMatch := false
				data := strings.TrimSpace(t.Data)
				//Check if action == Watched
				if data == "Watched" {
					getChannel = true
				}
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
					nextItem = "item"
				} else {
					nextItem = "date"
				}
			case "item":
				res.Item = t.Data
				if getChannel {
					nextItem = "channel"
					getChannel = false
				} else {
					nextItem = "date"
				}
			case "channel":
				res.Channel = t.Data
				nextItem = "date"
			case "date":
				layout, err := dateparse.ParseFormat(t.Data)
				if err != nil {
					res.Date = ""
				}
				t, err := time.Parse(layout, t.Data)
				if err != nil {
					res.Date = ""
				} else {
					res.Date = fmt.Sprintf("%02d-%02d-%02dT%02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
				}

				results = append(results, res)
				// Clear results
				res.Title = ""
				res.Action = ""
				res.Item = ""
				res.Date = ""
				nextItem = "unknown"
			}

		case tt == html.StartTagToken:
			if nextItem == "unknown" {
				//Check for new item denoted by class="mdl-typography--title"
				_, val, more := z.TagAttr()
				if string(val) == "mdl-typography--title" {
					// Next item will be title
					nextItem = "title"
				}
				for more {
					_, val, more = z.TagAttr()
					if string(val) == "mdl-typography--title" {
						// Next item will be title
						nextItem = "title"
					}
				}
			}
		}

	}
}
