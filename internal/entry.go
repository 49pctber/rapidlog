package rapidlog

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

var EntryTypes = make(map[string]string)
var reader *bufio.Reader
var Re_quit, Re_list, Re_delete, Re_summary *regexp.Regexp

var ErrInvalidEntry error = errors.New("invalid entry type")

const TimestampFormat string = "2006 Jan 02 at 15:04"

type Entry struct {
	Id        string
	Timestamp time.Time
	Type      string
	Entry     string
}

func init() {
	reader = bufio.NewReader(os.Stdin)

	EntryTypes["."] = "Todo"
	EntryTypes["-"] = "Note"
	EntryTypes["o"] = "Event"
	EntryTypes["?"] = "Question"
	EntryTypes["="] = "Emotion"

	Re_quit = regexp.MustCompile("^quit|^exit")
}

func (e *Entry) Log() error {

	insert_statement, err := Db.Prepare(`INSERT OR REPLACE INTO entries (id, timestamp, type, content) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	_, err = insert_statement.Exec(e.Id, e.Timestamp, e.Type, e.Entry)
	if err != nil {
		return err
	}

	return nil
}

func (e Entry) String() string {
	return fmt.Sprintf("%s %s", e.Type, e.Entry)
}

func (e *Entry) Print(verbose bool) {
	if verbose {
		fmt.Printf("%s (%s) \n%v\n", e.Timestamp.UTC().Local().Format(TimestampFormat), e.Id, e)
	} else {
		fmt.Printf("%v\n", e)
	}
}

func GetEntries(entry_type string, time_constraint string) (*sql.Rows, error) {
	if time_constraint == "" {
		time_constraint = "100 years"
	}

	time_constraint = "-" + time_constraint

	if entry_type != "" {
		return Db.Query(`SELECT id, timestamp, type, content FROM entries WHERE type=? AND timestamp >= datetime('now', ?) ORDER BY timestamp ASC`, entry_type, time_constraint)
	} else {
		return Db.Query(`SELECT id, timestamp, type, content FROM entries WHERE timestamp >= datetime('now', ?) ORDER BY timestamp ASC`, time_constraint)
	}
}

func GetEntry(id string) (*Entry, error) {

	entry := Entry{}

	row := Db.QueryRow(`SELECT id, timestamp, type, content FROM entries WHERE id=?`, id)

	err := row.Scan(&entry.Id, &entry.Timestamp, &entry.Type, &entry.Entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func ListEntries(entry_type, time_constraint string, verbose bool) error {

	entry_type = strings.TrimSpace(entry_type)
	if _, in_set := EntryTypes[entry_type]; !in_set {
		entry_type = ""
	}

	time_constraint = strings.TrimSpace(time_constraint)

	rows, err := GetEntries(entry_type, time_constraint)
	if err != nil {
		return err
	}

	prev_date := ""
	for rows.Next() {
		var entry Entry
		rows.Scan(&entry.Id, &entry.Timestamp, &entry.Type, &entry.Entry)
		curr_date := entry.Timestamp.UTC().Local().Format("2006 Jan 02")
		if prev_date != curr_date {
			// print new date
			prev_date = curr_date
			fmt.Printf("\n---------------------------\n%s\n", curr_date)
			if verbose {
				fmt.Printf("---------------------------\n\n")
			}
		}
		entry.Print(verbose)
		if verbose {
			fmt.Printf("\n")
		}
	}

	if !verbose {
		fmt.Printf("\n")
	}

	return nil
}

func DeleteEntry(id string) error {

	row := Db.QueryRow(`SELECT COUNT(*) FROM entries WHERE id=?`, id)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("no entries found with that ID")
	}

	res, err := Db.Exec(`DELETE FROM entries WHERE id=?`, id)
	if err != nil {
		return err
	}

	if n, err := res.RowsAffected(); n == 0 || err != nil {

		if err != nil {
			return err
		}

		return errors.New("no rows deleted")
	}

	return nil

}
