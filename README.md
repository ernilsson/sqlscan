# sqlscan

A light-weight alternative for scanning SQL row data into structs using tags to identify column names. For when you want 
to avoid lengthy `rows.Scan(...)` statements without introducing extraneous functions into the project. The library does
not depend on any others and the core logic consist of less than 100 lines of code. The algorithm to match struct fields
to columns has a naive O(N^2) time complexity. This does however make the algorithm quite simple, thus an improved 
algorithm will not be implemented until a need for it rises. If you would like to see this implemented, or have a 
suggestion for it, then open an issue.

## Disclaimer

This library uses reflection and I myself is no expert in the topic, hence there are most probably improvements to be 
made on that part of the code. To contribute to this project, fork this repository and make your changes on a new branch 
which you then use to open a pull request into this repository.

The library currently does **not** support nested structs, if such a feature is desired then open a ticket and/or submit
a PR to include it.

## Installation

It is as easy as installing any other module, simply execute the command below whilst in your project directory.
```shell
$ go get github.com/ernilsson/sqlscan
```

## Use

Start by defining your entity as a struct and add `sql` tags with values corresponding to the names of the columns that 
should be used to populate the field. 

```go
package main

type Entity struct {
	ID int `sql:"id"`
	Name string `sql:"name"`
}
```

Scanning the row data into the struct is as simple as using the core SQL library. Create a new instance of the scanner 
and feed it a pointer to the struct that it should scan the data into.

```go
package main 

import (
	"fmt"
	"database/sql"
	"github.com/ernilsson/sqlscan"
)

type Entity struct {
    ID   int    `sql:"id"`
    Name string `sql:"name"`
}

func get(db *sql.DB) {
    rows, err := db.Query("SELECT id, name FROM entity")
    if err != nil {
        // Handle error
        return
    }   
    scanner := sqlscan.New(rows)
    for rows.Next() {
        var e Entity
        err = scanner.Scan(&e)
        if err != nil {
            // Handle error 
        }
        fmt.Printf("Scanned: %+v\n", e)
    }   
}
```

The library does not work directly towards the `*sql.Rows` implementation but rather an internal interface which is 
implemented by the row struct. This makes it somewhat simpler to test, although I would say that test-trust is
decreased. The scannable interface has a mock for internal use in the test file.