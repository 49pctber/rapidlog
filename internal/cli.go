package rapidlog

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
)

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
