/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	rapidlog "github.com/49pctber/rapidlog/internal"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func getClient(config *oauth2.Config) *http.Client {
	tokFile := filepath.Join(rapidlog.GetCacheDir(), "token.json")
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// downloadFile downloads a file from Google Drive
func syncFile(srv *drive.Service, fileId string) error {
	// Get the file metadata
	file, err := srv.Files.Get(fileId).Fields("name, mimeType").Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve file metadata: %v", err)
	}

	fmt.Printf("Downloading file: %s (%s)\n", file.Name, file.MimeType)

	// Get the file content
	response, err := srv.Files.Get(fileId).Download()
	if err != nil {
		return fmt.Errorf("unable to download file: %v", err)
	}
	defer response.Body.Close()

	// Create a file locally
	outFile, err := os.CreateTemp("", "rapidlog")
	if err != nil {
		return fmt.Errorf("unable to create local file: %v", err)
	}
	defer outFile.Close()

	// Copy the content to the local file
	_, err = io.Copy(outFile, response.Body)
	if err != nil {
		return fmt.Errorf("unable to save file content: %v", err)
	}

	// Import that database
	fmt.Println(outFile.Name())
	err = rapidlog.ImportDatabase(outFile.Name())
	if err != nil {
		return fmt.Errorf("unable to save file content: %v", err)
	}

	fmt.Printf("File synced successfully.\n")
	return nil
}

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		b, err := os.ReadFile(filepath.Join(rapidlog.GetCacheDir(), "credentials.json"))
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
		}

		config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}
		client := getClient(config)

		srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
		if err != nil {
			log.Fatalf("Unable to retrieve Drive client: %v", err)
		}

		// Sync files
		var df *drive.File = nil
		r, err := srv.Files.List().PageSize(10).Fields("nextPageToken, files(id, name)").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve files: %v", err)
		}
		fmt.Println("Files:")
		if len(r.Files) == 0 {
			fmt.Println("No files found.")
		} else {
			for _, i := range r.Files {
				fmt.Printf("%s (%s)\n", i.Name, i.Id)
				err = syncFile(srv, i.Id)
				if err != nil {
					log.Fatalln(err)
				}
				df = i
			}
		}

		f, err := os.Open(rapidlog.Dbpath)
		if err != nil {
			log.Fatalf("Unable to open file: %v", err)
		}
		defer f.Close()

		// upload updated file
		if df == nil {
			file := &drive.File{Name: "rapidlog.sqlite3"}
			driveFile, err := srv.Files.Create(file).Media(f).Do()
			if err != nil {
				log.Fatalf("Unable to create file: %v", err)
			}
			fmt.Printf("Created new file %s in Google drive.\n", driveFile.Name)
		} else {
			file := &drive.File{}
			driveFile, err := srv.Files.Update(df.Id, file).Media(f).Do()
			if err != nil {
				log.Fatalf("Unable to update file: %v", err)
			}
			fmt.Printf("Updated file %s in Google drive.\n", driveFile.Name)
		}

		fmt.Println("Sync was successful!")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
