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

func TestGetItemsFromYear(t *testing.T) {
	db, err := OpenDB(testHome + "takeout.db")
	if err != nil {
		t.Fatal(err)
	}

	years, err := GetYears(db)
	if err != nil {
		t.Fatal(err)
	}

	for _, year := range years {
		results, err := GetItemsFromYear(db, year)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(fmt.Sprintf("%d: %d", year, len(results)))
	}
}

func TestGetItems(t *testing.T) {
	db, err := OpenDB(testHome + "takeout.db")
	if err != nil {
		t.Fatal(err)
	}

	begin := int64(1557833405)
	end := int64(1581077662)

	results, err := GetItemsFromUnixtime(db, begin, end)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("%d-%d: %d results", begin, end, len(results)))

}

func TestGetYears(t *testing.T) {
	db, err := OpenDB(testHome + "takeout.db")
	if err != nil {
		t.Fatal(err)
	}

	years, err := GetYears(db)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(years)
}

func TestSearchItems(t *testing.T) {
	searchString := "google"
	db, err := OpenDB(testHome + "takeout.db")
	if err != nil {
		t.Fatal(err)
	}

	results, err := SearchItems(db, searchString)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Search for '%s' returned %d results", searchString, len(results)))
}
