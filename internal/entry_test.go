package rapidlog

import (
	"testing"
)

func TestEntryParseString(t *testing.T) {
	var e Entry

	valid_entry_strings := []string{
		"? this is a valid question",
		"- this is a valid note",
		".this is a valid to-do item",
		"=  this is a valid emotion     ",
		"o this is a valid o event!!!! ? - . =",
	}

	for _, entry_string := range valid_entry_strings {
		err := e.ParseString(entry_string)
		if err != nil {
			t.Errorf(`"%s" is a valid entry string`, entry_string)
		}
	}

	invalid_entry_strings := []string{
		"this is an invalid string",
		"  ? this is an invalid question",
		"* this is an invalid string",
		"     ",
		"& ? - . =",
	}

	for _, entry_string := range invalid_entry_strings {
		err := e.ParseString(entry_string)
		if err == nil {
			t.Errorf(`"%s" is an invalid entry string`, entry_string)
		}
	}

}
