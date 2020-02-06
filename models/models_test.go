package models

import (
	"testing"
)

const testHome = "../test/"

func TestCreateDB(t *testing.T) {
	_, err := OpenDB(testHome + "takeout.db")
	if err != nil {
		t.Fatal(err)
	}
}

func TestInsert(t *testing.T) {
	db, err := OpenDB(testHome + "takeout.db")
	if err != nil {
		t.Fatal(err)
	}

	res := Result{
		Title:   "test",
		Action:  "test",
		Item:    "test",
		Channel: "test",
		Date:    "test",
	}
	err = InsertItem(db, res)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	db, err := OpenDB(testHome + "takeout.db")
	if err != nil {
		t.Fatal(err)
	}

	res := Result{
		Title:   "test",
		Action:  "test",
		Item:    "test",
		Channel: "test",
		Date:    "test",
	}
	err = DeleteItem(db, res)
	if err != nil {
		t.Fatal(err)
	}
}
