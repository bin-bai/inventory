# inventory
A simple document database interface and a thread safe flat json implementation.

Story
-----
I have been looking for a simple database to store strings by groups, which the number is less than 1K.
I found "https://github.com/HouzuoGuo/tiedot", which is pretty well designed and implemented.
The only concern is that's a liiittle bit heavy for my requirement, I just need a lite version of Tiedot.
So I wrote a simple document database which is inspired by Tiedot, and simplifize the get function with reflection.

Get source code
---------------
```
    go get github.com/bin-bai/inventory
```

How to use
----------
Here is example code from flatjson_test:
```
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
```
