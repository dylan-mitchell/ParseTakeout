package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dylan-mitchell/ParseTakeout"
	"github.com/dylan-mitchell/ParseTakeout/models"
)

func main() {
	html := flag.String("html", "", "Path to HTML from Google Takeout Data to parse")
	dbPath := flag.String("db", "", "Path to SQLITE3 DB")
	flag.Parse()

	if len(*html) == 0 {
		log.Fatal("Please specify a html file to parse")
	}

	if len(*dbPath) == 0 {
		log.Fatal("Please specify a db file")
	}

	results, err := ParseTakeout.ParseHTML(*html)
	if err != nil {
		log.Fatal(err)
	}

	db, err := models.OpenDB(*dbPath)
	if err != nil {
		log.Fatal(err)
	}

	//Do something with the results
	for _, result := range results {
		err := result.Validate()
		if err == nil {
			err := models.InsertItem(db, result)
			if err != nil {
				fmt.Println(result)
				fmt.Println(err)
			}
		}

	}
}
