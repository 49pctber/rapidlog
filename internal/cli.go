package rapidlog

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
)

func OpenEditor(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Use 'notepad' on Windows
		cmd = exec.Command("notepad", filePath)
	case "darwin":
		// Use 'open -e' which opens the file in the default text editor on macOS
		cmd = exec.Command("open", "-e", filePath)
	case "linux":
		// Use 'xdg-open' to open the file in the default text editor on Linux
		cmd = exec.Command("xdg-open", filePath)
	default:
		return fmt.Errorf("unsupported platform")
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func OpenFile(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", filePath)
	case "darwin":
		cmd = exec.Command("open", filePath)
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}

func GetUserInput() (string, error) {
	entry, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return "", nil
	}

	return strings.TrimSpace(entry), nil
}

func (e *Entry) ParseString(input string) error {
	if len(input) == 0 {
		return errors.New("zero-length string")
	}

	for key := range EntryTypes {
		if key == string(input[0]) {
			now := time.Now().UTC()
			id, err := ksuid.NewRandomWithTime(now)
			if err != nil {
				return err
			}
			e.Id = id.String()
			e.Timestamp = now
			e.Type = string(key)
			e.Entry = strings.TrimSpace(input[1:])
			return nil
		}
	}

	return ErrInvalidEntry
}

func CliInterface() error {

	defer Db.Close()
	fmt.Println("What's up?")

	for {
		// get user input
		input, err := GetUserInput()
		if err != nil {
			return err
		}

		// check if user wants to quit
		if Re_quit.MatchString(input) {
			// render summary on clean exit
			_, err := RenderSummary()
			return err
		}

		// log a new entry
		var entry Entry
		err = entry.ParseString(input)
		if err != nil {
			fmt.Println("Entries must begin with -, ., o, =, or ?. Use the up arrow to edit your previous input.")
			continue
		}
		err = entry.Log()
		if err != nil {
			return err
		}
	}
}
