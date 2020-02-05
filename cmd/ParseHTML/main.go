package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

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

	resultsJson, err := json.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}

	filePath := "test.json"

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString(string(resultsJson))

	//Do something with the results
	for _, result := range results {
		fmt.Println(result)
	}
}
