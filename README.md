# Go MYSQL Dump
Create MYSQL dumps in Go without the `mysqldump` CLI as a dependancy.

### Simple Example
```go
package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/sparkedhost/go-mysqldump"
)

func main() {
	// Open connection to database
	config := mysql.NewConfig()
	config.User = "your-user"
	config.Passwd = "your-pw"
	config.DBName = "your-db"
	config.Net = "tcp"
	config.Addr = "your-hostname:your-port"

	dumpDir := "dumps"
	timeFormat := fmt.Sprintf(time.Now().Format(time.RFC3339)) // accepts time layout string and add .sql at the end of file
	fileName := fmt.Sprintf("%s", config.DBName)

	if err := os.MkdirAll(dumpDir, 0755); err != nil {
		fmt.Println("Error mkdir:", err)
		return
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}

	// Register database with mysqldump
	dumper, err := mysqldump.Register(db, dumpDir, timeFormat, fileName)
	if err != nil {
		fmt.Println("Error registering databse:", err)
		return
	}

	// Dump database to file
	if err := dumper.Dump(); err != nil {
		fmt.Println("Error dumping:", err)
		return
	}
	if file, ok := dumper.Out.(*os.File); ok {
		fmt.Println("File is saved to", file.Name())
	} else {
		fmt.Println("It's not part of *os.File, but dump is done")
	}

	// Close dumper, connected database and file stream.
	dumper.Close()
}
```

[![GoDoc](https://godoc.org/github.com/jamf/go-mysqldump?status.svg)](https://godoc.org/github.com/jamf/go-mysqldump)
[![Build Status](https://travis-ci.org/jamf/go-mysqldump.svg?branch=master)](https://travis-ci.org/jamf/go-mysqldump)
