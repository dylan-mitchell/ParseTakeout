package ParseTakeout

import (
	"fmt"
	"testing"
)

const testHome = "./test/"

func TestReadHTML(t *testing.T) {
	_, err := ReadHtml(testHome + "My-Activity-Developers.html")
	if err != nil {
		t.Fatal(err)
	}
	//	fmt.Println(s)
}

func TestParseHTML(t *testing.T) {
	results, err := ParseHTML(testHome + "My-Activity-Developers.html")
	if err != nil {
		t.Fatal(err)
	}

	for _, result := range results {
		fmt.Println(result)
	}
}
