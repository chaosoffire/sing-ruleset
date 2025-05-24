package network

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFileFromURL(url string, filePath string) error {

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Create the HTTP client
	reponse, err := http.Get(url)
	if err != nil {
		return err
	}
	defer reponse.Body.Close()
	// Check for HTTP Status
	if reponse.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", reponse.Status)
	}
	// Write the body to file
	_, err = io.Copy(out, reponse.Body)
	if err != nil {
		return err
	}
	return nil
}
