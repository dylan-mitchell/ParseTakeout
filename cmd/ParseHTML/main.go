package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dylan-mitchell/ParseTakeout"
)

func main() {
	html := flag.String("html", "", "Path to HTML from Google Takeout Data to parse")
	flag.Parse()

	if len(*html) == 0 {
		log.Fatal("Please specify a html file to parse")
	}

	results, err := ParseTakeout.ParseHTML(*html)
	if err != nil {
		log.Fatal(err)
	}

	//Do something with the results
	for _, result := range results {
		fmt.Println(result)
	}
}
