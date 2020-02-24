package ParseTakeout

import (
	"fmt"
	"testing"
)

func TestLoadJSON(t *testing.T) {
	_, err := LoadJSON(testHome + "Location History.json")
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(s)
}

func TestInsertLocation(t *testing.T) {
	db, err := OpenDB(testHome + "takeout.db")
	if err != nil {
		t.Fatal(err)
	}

	loc := Location{
		Unixtime:  1234,
		Latitude:  1,
		Longitude: 2,
	}
	err = InsertLocation(db, loc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteLocation(t *testing.T) {
	db, err := OpenDB(testHome + "takeout.db")
	if err != nil {
		t.Fatal(err)
	}

	loc := Location{
		Unixtime:  1234,
		Latitude:  1,
		Longitude: 2,
	}
	err = DeleteLocation(db, loc)
	if err != nil {
		t.Fatal(err)
	}
}

// func TestInsertJSON(t *testing.T) {
// 	db, err := OpenDB(testHome + "takeout.db")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	data, err := LoadJSON(testHome + "Location History.json")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	err = InsertJSON(db, *data)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

func TestGetAllLocations(t *testing.T) {
	db, err := OpenDB(testHome + "takeout.db")
	if err != nil {
		t.Fatal(err)
	}

	results, err := GetAllLocations(db)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(len(results))
}
