package ParseTakeout

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
	Locations []Location `json:"locations"`
}

type Location struct {
	Timestamp string `json:"timestampMs"`
	Latitude  int64  `json:"latitudeE7"`
	Longitude int64  `json:"longitudeE7"`
}

func (l Location) String() string {
	s := fmt.Sprintf(`*****
Timestamp: %s
Lat: %d
Lon: %d
*****`, l.Timestamp, l.Latitude, l.Longitude)
	return s
}

func LoadJSON(filePath string) (*Data, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data Data

	json.Unmarshal(bytes, &data)

	return &data, nil
}

func InsertLocation(db *sql.DB, loc Location) error {
	ts, err := strconv.Atoi(loc.Timestamp)
	if err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf(`
	INSERT INTO "locationhistory" ("timestamp", "latitude", "longitude")
	VALUES ("%d", "%d", "%d");
	`, ts, loc.Latitude, loc.Longitude))
	if err != nil {
		return err
	}
	return nil
}

func DeleteLocation(db *sql.DB, loc Location) error {
	ts, err := strconv.Atoi(loc.Timestamp)
	if err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf(`
	DELETE FROM "locationhistory" WHERE
	"timestamp" = "%d";
	`, ts))
	if err != nil {
		return err
	}
	return nil
}

func parseLocationRows(rows *sql.Rows) ([]Location, error) {
	defer rows.Close()

	results := []Location{}
	for rows.Next() {
		var ts int64
		var lat int64
		var lon int64
		if err := rows.Scan(&ts, &lat, &lon); err != nil {
			return nil, err
		}

		tsString := fmt.Sprintf("%d", ts)
		results = append(results, Location{
			Timestamp: tsString,
			Latitude:  lat,
			Longitude: lon,
		})
		// Check for errors from iterating over rows.
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return results, nil
}

func GetAllLocations(db *sql.DB) ([]Location, error) {
	rows, err := db.Query(`
	SELECT * FROM "locationhistory";
	`)
	if err != nil {
		return nil, err
	}

	results, err := parseLocationRows(rows)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func InsertJSON(db *sql.DB, data Data) error {
	for _, loc := range data.Locations {
		err := InsertLocation(db, loc)
		if err != nil {
			return err
		}
	}
	return nil
}
