package flatjson_test

import (
	"fmt"
	"testing"

	"github.com/bin-bai/inventory"
	"github.com/bin-bai/inventory/flatjson"
)

func TestFlatJson(t *testing.T) {
	var db inventory.Database

	db = flatjson.NewFlatJson("/tmp/testdb.json")
	err := db.Load()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Test Collection begin
	col := db.Use("First", "")
	err = col.Set("1.1", "1-1")
	if err != nil {
		t.Fatal(err)
	}

	vi := col.Get("1.1")
	if vi == nil {
		t.Fatal("Can't get value of key 1.1")
	}

	v, ok := vi.(string)
	if !ok {
		t.Fatal("Can't convert Get result to string")
	}
	if v != "1-1" {
		t.Fatal("string value is wrong")
	}
	fmt.Println(v)
	// Test Collection end

	err = db.Save()
	if err != nil {
		t.Fatal(err)
	}
}
