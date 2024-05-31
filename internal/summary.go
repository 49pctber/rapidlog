package rapidlog

import (
	"html/template"
	"os"
	"path/filepath"
	"time"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed summary.tmpl
var tmpl_str string

func RenderSummary() (*string, error) {

	rows, err := GetEntries("", "100 years")
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("summary").Parse(tmpl_str)
	if err != nil {
		return nil, err
	}

	var entries []Entry
	for rows.Next() {
		var entry Entry
		rows.Scan(&entry.Id, &entry.Timestamp, &entry.Type, &entry.Entry)
		entries = append(entries, entry)
	}

	// reverse the entries for the template
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	// save latest copy
	fname := filepath.Join(GetCacheDir(), time.Now().UTC().Local().Format("summary.html"))
	f, err := os.Create(fname)
	if err != nil {
		return nil, err
	}

	err = tmpl.Execute(f, struct {
		Entries []Entry
		Time    time.Time
	}{entries, time.Now()})
	if err != nil {
		f.Close()
		return nil, err
	}
	f.Close()

	return &fname, nil
}
