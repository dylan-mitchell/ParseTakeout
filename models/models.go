package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"

	_ "github.com/mattn/go-sqlite3"
)

type Result struct {
	Title    string `json:"title"`
	Action   string `json:"action"`
	Item     string `json:"item"`
	Channel  string `json:"channel"`
	Date     string `json:"date"`
	UnixTime int64  `json:"unixtime"`
}

func (r Result) String() string {
	if len(r.Channel) == 0 {
		s := fmt.Sprintf(`*****
	Title: %s
	Action: %s
	Item: %s
	Date: %s
	UnixTime: %d
	*****`, r.Title, r.Action, r.Item, r.Date, r.UnixTime)
		return s
	} else {
		s := fmt.Sprintf(`*****
	Title: %s
	Action: %s
	Item: %s
	Channel: %s
	Date: %s
	UnixTime: %d
	*****`, r.Title, r.Action, r.Item, r.Channel, r.Date, r.UnixTime)
		return s
	}
}

func (r Result) Validate() error {
	if r.Title == "" {
		return errors.New("Empty title")
	}
	if r.Action == "" {
		return errors.New("Empty action")
	}
	if r.Item == "" {
		return errors.New("Empty item")
	}
	if r.UnixTime == 0 {
		return errors.New("Invalid unix time")
	}
	if r.Date == "" {
		return errors.New("Empty date")
	}
	return nil
}

func OpenDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Error opening")
		return nil, err
	}

	sqlStmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "items" (
		"title"	TEXT,
		"action"	TEXT,
		"item"	TEXT,
		"channel"	TEXT,
		"date"	TEXT,
		"unixtime"	INTEGER,
		PRIMARY KEY("unixtime","item")
	);
	`)
	if err != nil {
		return nil, err
	}

	_, err = sqlStmt.Exec()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InsertItem(db *sql.DB, res Result) error {
	_, err := db.Exec(fmt.Sprintf(`
	INSERT INTO "items" ("title", "action", "item", "channel", "date", "unixtime")
	VALUES ("%s", "%s", "%s", "%s", "%s", "%d");
	`, url.QueryEscape(res.Title), url.QueryEscape(res.Action), url.QueryEscape(res.Item), url.QueryEscape(res.Channel), url.QueryEscape(res.Date), res.UnixTime))
	if err != nil {
		return err
	}
	return nil
}

func DeleteItem(db *sql.DB, res Result) error {
	_, err := db.Exec(fmt.Sprintf(`
	DELETE FROM "items" WHERE
	"unixtime" = "%d";
	`, res.UnixTime))
	if err != nil {
		return err
	}
	return nil
}

// func InsertMultiple(dbPath string, results []Result) error {
// 	fmt.Println(len(results))
// 	for _, res := range results {
// 		db, err := OpenDB(dbPath)
// 		if err != nil {
// 			return err
// 		}
// 		err = InsertItem(db, res)
// 		if err != nil {
// 			return err
// 		}
// 		db.Close()
// 	}
// 	return nil
// }
