package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
)

// Load database from source
func DB(dataSourceName string) (*sql.DB, error)  {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

/// Check and prepare table if not existing in source database
func Prepare(tableName string, withCreateQuery string, dataSourceName string) {
	db, err := DB(dataSourceName)
	checkErr(err)
	defer db.Close()

	checkQuery := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%v';",tableName)

	checkResult, err := db.Query(checkQuery)
	checkErr(err)

	var name string
	for checkResult.Next()  {
		checkResult.Scan(&name)
	}
	if len(name) == 0 {
		createTable(dataSourceName,tableName,withCreateQuery)
	}
}

func createTable(dataSourceName string, tableName string, withCreateQuery string)  {
	db, err := DB(dataSourceName)
	checkErr(err)
	defer db.Close()

	stmt, err := db.Prepare(withCreateQuery)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("\n Table %v  created....",tableName)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}