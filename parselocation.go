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
	Unixtime  int64 `json:"unixtime"`
	Latitude  int64 `json:"latitude"`
	Longitude int64 `json:"longitude"`
}

type DataInput struct {
	Locations []LocationInput `json:"locations"`
}

type LocationInput struct {
	Timestamp string `json:"timestampMs"`
	Latitude  int64  `json:"latitudeE7"`
	Longitude int64  `json:"longitudeE7"`
}

func (l Location) String() string {
	s := fmt.Sprintf(`*****
Unixtime: %d
Lat: %d
Lon: %d
*****`, l.Unixtime, l.Latitude, l.Longitude)
	return s
}

func LoadJSON(filePath string) (*DataInput, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data DataInput

	json.Unmarshal(bytes, &data)

	return &data, nil
}

func FormatInput(loc LocationInput) Location {
	t, _ := strconv.Atoi(loc.Timestamp[0 : len(loc.Timestamp)-4])
	return Location{
		Unixtime:  int64(t),
		Latitude:  loc.Latitude / 10000000,
		Longitude: loc.Longitude / 10000000,
	}
}

func InsertLocation(db *sql.DB, loc Location) error {
	_, err := db.Exec(fmt.Sprintf(`
	INSERT INTO "locationhistory" ("timestamp", "latitude", "longitude")
	VALUES ("%d", "%d", "%d");
	`, loc.Unixtime, loc.Latitude, loc.Longitude))
	if err != nil {
		return err
	}
	return nil
}

func DeleteLocation(db *sql.DB, loc Location) error {

	_, err := db.Exec(fmt.Sprintf(`
	DELETE FROM "locationhistory" WHERE
	"unixtime" = "%d";
	`, loc.Unixtime))
	if err != nil {
		return err
	}
	return nil
}

func parseLocationRows(rows *sql.Rows) ([]Location, error) {
	defer rows.Close()

	results := []Location{}
	for rows.Next() {
		var t int64
		var lat int64
		var lon int64
		if err := rows.Scan(&t, &lat, &lon); err != nil {
			return nil, err
		}
		results = append(results, Location{
			Unixtime:  t,
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
