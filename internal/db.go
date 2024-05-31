package rapidlog

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/segmentio/ksuid"
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

	_, err = Db.Exec(`CREATE TABLE IF NOT EXISTS entries (id TEXT NOT NULL PRIMARY KEY, timestamp DATETIME, type TEXT, content TEXT)`)
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
	insert_statement, err := activedb.Prepare(`INSERT INTO entries (id, timestamp, type, content) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("error reading database: %v", err)
	}

	// open imported database
	db, err := sql.Open("sqlite3", fname)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()
	rows, err := db.Query(`SELECT id, timestamp, type, content FROM entries`)
	if err != nil {
		return fmt.Errorf("error reading database: %v", err)
	}
	defer rows.Close()

	// add data to active database
	for rows.Next() {
		var entry Entry
		rows.Scan(&entry.Id, &entry.Timestamp, &entry.Type, &entry.Entry)

		// check if entry already exists
		row := activedb.QueryRow(`SELECT COUNT(*) FROM entries WHERE entries.id=? OR entries.timestamp=?`, entry.Id, entry.Timestamp)
		var count int
		row.Scan(&count)
		if count != 0 {
			continue
		}
		entry.Print(false)

		// add to active database
		k, err := ksuid.Parse(entry.Id)
		if err != nil {
			new_id, err := ksuid.NewRandomWithTime(entry.Timestamp)
			if err != nil {
				return fmt.Errorf("error producing new ID: %v", err)
			}

			entry.Id = new_id.String()
		} else {
			entry.Id = k.String()
		}
		_, err = insert_statement.Exec(entry.Id, entry.Timestamp, entry.Type, entry.Entry)
		if err != nil {
			return fmt.Errorf("error adding entry to database: %v", err)
		}
	}

	return nil
}
