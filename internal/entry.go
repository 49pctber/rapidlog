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

var DefaultTime time.Time

type Entry struct {
	Id        string
	Timestamp time.Time
	Type      string
	Entry     string
	Modified  time.Time
}

func init() {
	reader = bufio.NewReader(os.Stdin)

	EntryTypes["."] = "Todo"
	EntryTypes["-"] = "Note"
	EntryTypes["o"] = "Event"
	EntryTypes["?"] = "Question"
	EntryTypes["="] = "Emotion"

	Re_quit = regexp.MustCompile("^quit|^exit")

	var err error
	DefaultTime, err = time.Parse("2006 01 02", "2000 01 01")
	if err != nil {
		panic(err)
	}
}

func (e *Entry) Log(use_now bool) error {

	insert_statement, err := Db.Prepare(`INSERT OR REPLACE INTO entries (id, timestamp, type, content, modified) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	if use_now {
		_, err = insert_statement.Exec(e.Id, e.Timestamp, e.Type, e.Entry, time.Now().UTC())
	} else {
		_, err = insert_statement.Exec(e.Id, e.Timestamp, e.Type, e.Entry, e.Modified)
	}

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
		fmt.Printf("%s %s | %s\n%v\n\n", e.Type, e.Timestamp.UTC().Local().Format(TimestampFormat), e.Id, e.Entry)
	} else {
		fmt.Printf("%v\n", e)
	}
}

func (e *Entry) MarkCompleted() error {

	if e.Type != "." {
		return fmt.Errorf("entry not a todo item (type %s, not .)", e.Type)
	}

	e.Type = "x"
	return e.Log(true)
}

func (e *Entry) MarkIncomplete() error {

	if e.Type != "x" {
		return fmt.Errorf("entry not a completed todo item (type %s, not x)", e.Type)
	}

	e.Type = "."
	return e.Log(true)
}

func GetEntries(entry_type string, time_constraint string) ([]Entry, error) {
	if time_constraint == "" {
		time_constraint = "100 years"
	}

	time_constraint = "-" + time_constraint

	var rows *sql.Rows
	var err error
	if entry_type != "" {
		rows, err = Db.Query(`SELECT id, timestamp, type, content, modified FROM entries WHERE type=? AND timestamp >= datetime('now', ?) ORDER BY timestamp ASC`, entry_type, time_constraint)
	} else {
		rows, err = Db.Query(`SELECT id, timestamp, type, content, modified FROM entries WHERE timestamp >= datetime('now', ?) ORDER BY timestamp ASC`, time_constraint)
	}
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0)

	for rows.Next() {
		var entry Entry
		rows.Scan(&entry.Id, &entry.Timestamp, &entry.Type, &entry.Entry, &entry.Modified)
		entries = append(entries, entry)
	}

	return entries, nil
}

func GetEntry(id string) (*Entry, error) {

	entry := Entry{}

	row := Db.QueryRow(`SELECT id, timestamp, type, content, modified FROM entries WHERE id=?`, id)

	err := row.Scan(&entry.Id, &entry.Timestamp, &entry.Type, &entry.Entry, &entry.Modified)
	if err != nil {
		row = Db.QueryRow(`SELECT id, timestamp, type, content FROM entries WHERE id=?`, id)
		err = row.Scan(&entry.Id, &entry.Timestamp, &entry.Type, &entry.Entry)
		if err != nil {
			return nil, err
		}
		entry.Modified = DefaultTime
	}

	return &entry, nil
}

func PrintEntries(entries []Entry, verbose bool, withdate bool) {
	if withdate {
		prev_date := ""
		for _, entry := range entries {
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
		}
	} else {
		for _, entry := range entries {
			entry.Print(verbose)
		}
	}

}

func ListEntries(entry_type, time_constraint string, verbose bool) error {

	entry_type = strings.TrimSpace(entry_type)
	if _, in_set := EntryTypes[entry_type]; !in_set {
		entry_type = ""
	}

	time_constraint = strings.TrimSpace(time_constraint)

	entries, err := GetEntries(entry_type, time_constraint)
	if err != nil {
		return err
	}

	PrintEntries(entries, verbose, true)

	if verbose {
		fmt.Printf("[Listed entries less than %s old. Change with -t flag.]\n", time_constraint)
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
