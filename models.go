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

type TotalSummary struct {
	MostCommon    []ItemFreq      `json:"mostcommon"`
	YoutubeTotal  int             `json:"youtubetotal"`
	ChannelCommon []ChannelFreq   `json:"channelcommon"`
	Total         int             `json:"total"`
	Yearly        []YearlySummary `json:"yearly"`
}

type YearlySummary struct {
	Year          int            `json:"year"`
	MostCommon    []ItemFreq     `json:"mostcommon"`
	YoutubeTotal  int            `json:"youtubetotal"`
	ChannelCommon []ChannelFreq  `json:"channelcommon"`
	Total         int            `json:"total"`
	Monthly       []MonthSummary `json:"monthly"`
}

type MonthSummary struct {
	Name  string `json:"name"`
	Begin int64  `json:"begin"`
	End   int64  `json:"end"`
	Total int    `json:"total"`
}

type ItemFreq struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type ChannelFreq struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
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
	if len(r.Item) > 250 {
		return errors.New("Item too long")
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
		PRIMARY KEY("action","unixtime","item")
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

func GetItemsFromUnixtime(db *sql.DB, begin, end int64) ([]Result, error) {
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

func constructMonthlySummary(db *sql.DB, year int) ([]MonthSummary, error) {
	months := []string{
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	}

	var mSum []MonthSummary

	var monthCount time.Month = 1
	for _, month := range months {
		begin := time.Date(year, monthCount, 1, 0, 0, 0, 0, time.UTC).Unix()
		end := time.Date(year, monthCount+1, 0, 23, 59, 59, 999999999, time.UTC).Unix()

		sum, err := db.Query(fmt.Sprintf(`
		SELECT COUNT(*) FROM "items"
		WHERE "unixtime" > %d AND "unixtime" < %d;
		`, begin, end))
		if err != nil {
			return nil, err
		}
		defer sum.Close()

		var sumCount int
		for sum.Next() {
			if err := sum.Scan(&sumCount); err != nil {
				return nil, err
			}
		}
		// Check for errors from iterating over rows.
		if err := sum.Err(); err != nil {
			return nil, err
		}

		mSum = append(mSum, MonthSummary{
			Name:  month,
			Begin: begin,
			End:   end,
			Total: sumCount,
		})

		monthCount++
	}

	return mSum, nil
}

func getMostCommonForYear(db *sql.DB, year int) ([]ItemFreq, error) {
	begin, end := calculateUnixRangeOfYear(year)

	freqs, err := db.Query(fmt.Sprintf(`
	SELECT "item", COUNT(*) AS FREQ
	FROM "items"
	WHERE "unixtime" > %d AND "unixtime" < %d
	GROUP BY "item"
	ORDER BY COUNT(*) DESC
	LIMIT 10;
	`, begin, end))
	if err != nil {
		return nil, err
	}
	defer freqs.Close()

	var itemFreqs []ItemFreq
	for freqs.Next() {
		var freqTotal int
		var item string
		if err := freqs.Scan(&item, &freqTotal); err != nil {
			return nil, err
		}
		item, _ = url.QueryUnescape(item)
		itemFreqs = append(itemFreqs, ItemFreq{
			Name:  item,
			Count: freqTotal,
		})
	}
	// Check for errors from iterating over rows.
	if err := freqs.Err(); err != nil {
		return nil, err
	}

	return itemFreqs, nil
}

func getCountForYear(db *sql.DB, year int) (int, error) {
	begin, end := calculateUnixRangeOfYear(year)

	count, err := db.Query(fmt.Sprintf(`
	SELECT COUNT(*)
	FROM "items"
	WHERE "unixtime" > %d AND "unixtime" < %d;
	`, begin, end))
	if err != nil {
		return 0, err
	}
	defer count.Close()

	var sum int
	for count.Next() {
		if err := count.Scan(&sum); err != nil {
			return 0, err
		}
	}
	// Check for errors from iterating over rows.
	if err := count.Err(); err != nil {
		return 0, err
	}

	return sum, nil
}

func getMostCommonChannelForYear(db *sql.DB, year int) ([]ChannelFreq, error) {
	begin, end := calculateUnixRangeOfYear(year)

	freqs, err := db.Query(fmt.Sprintf(`
	SELECT "channel", COUNT(*) AS FREQ
	FROM "items"
	WHERE "unixtime" > %d AND "unixtime" < %d AND "channel" != ""
	GROUP BY "channel"
	ORDER BY COUNT(*) DESC
	LIMIT 10;
	`, begin, end))
	if err != nil {
		return nil, err
	}
	defer freqs.Close()

	var channelFreqs []ChannelFreq
	for freqs.Next() {
		var freqTotal int
		var channel string
		if err := freqs.Scan(&channel, &freqTotal); err != nil {
			return nil, err
		}
		channel, _ = url.QueryUnescape(channel)
		channelFreqs = append(channelFreqs, ChannelFreq{
			Name:  channel,
			Count: freqTotal,
		})
	}
	// Check for errors from iterating over rows.
	if err := freqs.Err(); err != nil {
		return nil, err
	}

	return channelFreqs, nil
}

func getYoutubeForYear(db *sql.DB, year int) (int, error) {
	begin, end := calculateUnixRangeOfYear(year)

	count, err := db.Query(fmt.Sprintf(`
	SELECT COUNT(*)
	FROM "items"
	WHERE "unixtime" > %d AND "unixtime" < %d AND "channel" != "";
	`, begin, end))
	if err != nil {
		return 0, err
	}
	defer count.Close()

	var sum int
	for count.Next() {
		if err := count.Scan(&sum); err != nil {
			return 0, err
		}
	}
	// Check for errors from iterating over rows.
	if err := count.Err(); err != nil {
		return 0, err
	}

	return sum, nil
}

func GetSummaryofYear(db *sql.DB, year int) (*YearlySummary, error) {

	monthly, err := constructMonthlySummary(db, year)
	if err != nil {
		return nil, err
	}
	common, err := getMostCommonForYear(db, year)
	if err != nil {
		return nil, err
	}
	channelCommon, err := getMostCommonChannelForYear(db, year)
	if err != nil {
		return nil, err
	}
	total, err := getCountForYear(db, year)
	if err != nil {
		return nil, err
	}
	youtubeTotal, err := getYoutubeForYear(db, year)
	if err != nil {
		return nil, err
	}

	yearlySum := YearlySummary{
		Year:          year,
		Monthly:       monthly,
		MostCommon:    common,
		ChannelCommon: channelCommon,
		Total:         total,
		YoutubeTotal:  youtubeTotal,
	}

	return &yearlySum, nil
}

func getMostCommonItem(db *sql.DB) ([]ItemFreq, error) {
	freqs, err := db.Query(`
	SELECT "item", COUNT(*) AS FREQ
	FROM "items"
	GROUP BY "item"
	ORDER BY COUNT(*) DESC
	LIMIT 10;
	`)
	if err != nil {
		return nil, err
	}
	defer freqs.Close()

	var itemFreqs []ItemFreq
	for freqs.Next() {
		var freqTotal int
		var item string
		if err := freqs.Scan(&item, &freqTotal); err != nil {
			return nil, err
		}
		item, _ = url.QueryUnescape(item)
		itemFreqs = append(itemFreqs, ItemFreq{
			Name:  item,
			Count: freqTotal,
		})
	}
	// Check for errors from iterating over rows.
	if err := freqs.Err(); err != nil {
		return nil, err
	}

	return itemFreqs, nil
}

func getCountTotal(db *sql.DB) (int, error) {

	count, err := db.Query(`
	SELECT COUNT(*)
	FROM "items";
	`)
	if err != nil {
		return 0, err
	}
	defer count.Close()

	var sum int
	for count.Next() {
		if err := count.Scan(&sum); err != nil {
			return 0, err
		}
	}
	// Check for errors from iterating over rows.
	if err := count.Err(); err != nil {
		return 0, err
	}

	return sum, nil
}

func getYoutubeTotal(db *sql.DB) (int, error) {

	count, err := db.Query(`
	SELECT COUNT(*)
	FROM "items"
	WHERE "channel" != "";
	`)
	if err != nil {
		return 0, err
	}
	defer count.Close()

	var sum int
	for count.Next() {
		if err := count.Scan(&sum); err != nil {
			return 0, err
		}
	}
	// Check for errors from iterating over rows.
	if err := count.Err(); err != nil {
		return 0, err
	}

	return sum, nil
}

func getMostCommonChannel(db *sql.DB) ([]ChannelFreq, error) {
	freqs, err := db.Query(`
	SELECT "channel", COUNT(*) AS FREQ
	FROM "items"
	WHERE "channel" != ""
	GROUP BY "channel"
	ORDER BY COUNT(*) DESC
	LIMIT 10;
	`)
	if err != nil {
		return nil, err
	}
	defer freqs.Close()

	var channelFreqs []ChannelFreq
	for freqs.Next() {
		var freqTotal int
		var channel string
		if err := freqs.Scan(&channel, &freqTotal); err != nil {
			return nil, err
		}
		channel, _ = url.QueryUnescape(channel)
		channelFreqs = append(channelFreqs, ChannelFreq{
			Name:  channel,
			Count: freqTotal,
		})
	}
	// Check for errors from iterating over rows.
	if err := freqs.Err(); err != nil {
		return nil, err
	}

	return channelFreqs, nil
}

// type TotalSummary struct {
// 	MostCommon    []ItemFreq      `json:"mostcommon"`
// 	YoutubeTotal  int             `json:"youtubetotal"`
// 	ChannelCommon []ChannelFreq   `json:"channelcommon"`
// 	Total         int             `json:"total"`
// 	Yearly        []YearlySummary `json:"yearly"`
// }

func GetTotalSummary(db *sql.DB) (*TotalSummary, error) {

	years, err := GetYears(db)
	if err != nil {
		return nil, err
	}

	var yearSums []YearlySummary
	for _, year := range years {
		yearSum, err := GetSummaryofYear(db, year)
		if err != nil {
			return nil, err
		}
		yearSums = append(yearSums, *yearSum)
	}

	mostCommon, err := getMostCommonItem(db)
	if err != nil {
		return nil, err
	}

	total, err := getCountTotal(db)
	if err != nil {
		return nil, err
	}

	youtubeTotal, err := getYoutubeTotal(db)
	if err != nil {
		return nil, err
	}

	channelCommon, err := getMostCommonChannel(db)
	if err != nil {
		return nil, err
	}

	totalSum := TotalSummary{
		MostCommon:    mostCommon,
		YoutubeTotal:  youtubeTotal,
		ChannelCommon: channelCommon,
		Total:         total,
		Yearly:        yearSums,
	}

	return &totalSum, nil
}

func GetYears(db *sql.DB) ([]int, error) {
	rows, err := db.Query(`
	SELECT MIN("unixtime"), MAX("unixtime") FROM "items";
	`)
	if err != nil {
		return nil, err
	}
	var min, max int64

	for rows.Next() {
		if err := rows.Scan(&min, &max); err != nil {
			return nil, err
		}

	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	begin := time.Unix(min, 0).Year()
	end := time.Unix(max, 0).Year()
	years := []int{}

	for i := begin; i <= end; i++ {
		years = append(years, i)
	}

	return years, nil
}

func SearchItems(db *sql.DB, searchString string) ([]Result, error) {
	rows, err := db.Query(`
	SELECT * FROM "items"
	WHERE "item" LIKE '%` + searchString + `%' ORDER BY "unixtime" ASC;
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
