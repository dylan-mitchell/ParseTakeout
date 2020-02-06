package ParseTakeout

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

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

func parseRows(rows *sql.Rows) ([]Result, error) {
	defer rows.Close()

	results := []Result{}
	for rows.Next() {
		var title string
		var action string
		var item string
		var channel string
		var date string
		var unixtime int64
		if err := rows.Scan(&title, &action, &item, &channel, &date, &unixtime); err != nil {
			return nil, err
		}
		title, _ = url.QueryUnescape(title)
		action, _ = url.QueryUnescape(action)
		item, _ = url.QueryUnescape(item)
		channel, _ = url.QueryUnescape(channel)
		date, _ = url.QueryUnescape(date)

		results = append(results, Result{
			Title:    title,
			Action:   action,
			Item:     item,
			Channel:  channel,
			Date:     date,
			UnixTime: unixtime,
		})
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func GetAllItems(db *sql.DB) ([]Result, error) {
	rows, err := db.Query(`
	SELECT * FROM "items";
	`)
	if err != nil {
		return nil, err
	}

	results, err := parseRows(rows)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func calculateUnixRangeOfYear(year int) (int64, int64) {
	begin := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	end := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC).Unix()

	return begin, end
}

func GetAllItemsFromYear(db *sql.DB, year int) ([]Result, error) {
	begin, end := calculateUnixRangeOfYear(year)

	rows, err := db.Query(fmt.Sprintf(`
	SELECT * FROM "items"
	WHERE "unixtime" > %d AND "unixtime" < %d;
	`, begin, end))
	if err != nil {
		return nil, err
	}

	results, err := parseRows(rows)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func GetItemsFromYear(db *sql.DB, year int) ([]Result, error) {
	begin, end := calculateUnixRangeOfYear(year)

	rows, err := db.Query(fmt.Sprintf(`
	SELECT * FROM "items"
	WHERE "unixtime" > %d AND "unixtime" < %d;
	`, begin, end))
	if err != nil {
		return nil, err
	}

	results, err := parseRows(rows)
	if err != nil {
		return nil, err
	}

	return results, nil
}
