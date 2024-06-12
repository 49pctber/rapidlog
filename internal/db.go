package rapidlog

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// database settings
var Dbname = "rapid_log.sqlite3"
var Dbpath string
var Db *sql.DB

func GetCacheDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(dir, "RapidLog")
}

func init() {

	var err error

	err = os.MkdirAll(GetCacheDir(), os.ModePerm)
	if err != nil {
		panic(err)
	}
	Dbpath = filepath.Join(GetCacheDir(), Dbname)

	Db, err = sql.Open("sqlite3", Dbpath)
	if err != nil {
		panic(err)
	}
	// defer Db.Close()

	_, err = Db.Exec(`CREATE TABLE IF NOT EXISTS entries (id TEXT NOT NULL PRIMARY KEY, timestamp DATETIME, type TEXT, content TEXT, modified DATETIME)`)
	if err != nil {
		panic(err)
	}

}

func ImportDatabase(fname string) error {

	// open active database
	activedb, err := sql.Open("sqlite3", Dbpath)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer activedb.Close()

	// open imported database
	db, err := sql.Open("sqlite3", fname)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	modified_included := true
	defer db.Close()
	rows, err := db.Query(`SELECT id, timestamp, type, content, modified FROM entries`)
	if err != nil {
		modified_included = false
		rows, err = db.Query(`SELECT id, timestamp, type, content FROM entries`)
		if err != nil {
			return fmt.Errorf("error reading database: %v", err)
		}
	}
	defer rows.Close()

	// add data to active database
	for rows.Next() {
		var entry Entry
		if modified_included {
			rows.Scan(&entry.Id, &entry.Timestamp, &entry.Type, &entry.Entry, &entry.Modified)
		} else {
			rows.Scan(&entry.Id, &entry.Timestamp, &entry.Type, &entry.Entry)
			entry.Modified = DefaultTime
		}

		// check if entry already exists
		row := activedb.QueryRow(`SELECT COUNT(*) FROM entries WHERE entries.id=? OR entries.timestamp=?`, entry.Id, entry.Timestamp)
		var count int
		row.Scan(&count)
		if count == 1 {

			existing_entry, err := GetEntry(entry.Id)
			if err != nil {
				panic(err)
			}

			if existing_entry.Modified.Before(entry.Modified) {
				err = entry.Log(false)
				if err != nil {
					return fmt.Errorf("error updating entry: %v", err)
				}
			}

		} else {
			// add entry to database
			err = entry.Log(false)
			if err != nil {
				return fmt.Errorf("error adding entry to database: %v", err)
			}

			// show new entry to user
			entry.Print(false)
		}

	}

	return nil
}
