package ParseTakeout

import (
	"fmt"
	"testing"
)

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

func TestGetAllItems(t *testing.T) {
	db, err := OpenDB(testHome + "takeout.db")
	if err != nil {
		t.Fatal(err)
	}

	results, err := GetAllItems(db)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(len(results))
}

func TestGetAllItemsFromYear(t *testing.T) {
	db, err := OpenDB(testHome + "takeout.db")
	if err != nil {
		t.Fatal(err)
	}

	results, err := GetAllItemsFromYear(db, 2015)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("2015: %d", len(results)))

	results, err = GetAllItemsFromYear(db, 2016)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("2016: %d", len(results)))

	results, err = GetAllItemsFromYear(db, 2017)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("2017: %d", len(results)))

	results, err = GetAllItemsFromYear(db, 2018)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("2018: %d", len(results)))

	results, err = GetAllItemsFromYear(db, 2019)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("2019: %d", len(results)))
}
