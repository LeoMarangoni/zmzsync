package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	//fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Print("\033[u\033[K")
	fmt.Printf(" %s complete", humanize.Bytes(wc.Total))
}

func getQueryString(types, start, end string) (string, error) {
	//Get the download filters and transform in a http query string
	// types must be a combination any of those: message,conversation,contact,appointment,task
	// start and end must be a date in the following format YYYY-MM-DD
	expectedTypes := "message,conversation,contact,appointment,task"
	query := "?fmt=tgz"
	if types == "" {
		types = expectedTypes
	} else {
		//Check if the asked types to download are valid
		for _, x := range strings.Split(types, ",") {
			valid := false
			for _, y := range strings.Split(expectedTypes, ",") {
				if x == y {
					valid = true
				}
			}
			if valid != true {
				return "", fmt.Errorf("Invalid type: %s. Valid types are: %s", x, expectedTypes)
			}
		}
	}
	query += "&types=" + types
	// Parsing time date if it exists, ADD time 00:00 to startdate, and 23:59 to enddate
	// and convert to Unix Timestamp in Millisecond format
	if start != "" {
		start, err := time.Parse(time.RFC3339, start+"T00:00:00.000Z")
		if err != nil {
			return "", err
		} else {
			timestamp := fmt.Sprint(start.UnixMilli())
			query = fmt.Sprintf("%s&start=%s", query, timestamp)
		}

	}
	if end != "" {
		end, err := time.Parse(time.RFC3339, end+"T23:59:59.999Z")
		if err != nil {
			return "", err
		} else {
			timestamp := fmt.Sprint(end.Unix())
			query = fmt.Sprintf("%s&end=%s", query, timestamp)
		}
	}
	return query, nil
}

func downloadMailbox(filepath, admin, account, password, host, port, types, start, end string) error {
	var login string
	// Use admin account to log in if it's setted
	if admin == "" {
		login = account
	} else {
		login = admin
	}
	query, err := getQueryString(types, start, end)
	if err != nil {
		return err
	}
	uri := fmt.Sprintf("https://%s:%s/home/%s/%s", host, port, account, query)
	//Create a temporary file to store the downlod data
	file, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}
	defer file.Close()
	// Create the http client, ignoring insecure connections
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	// Authenticate and send the request
	req.SetBasicAuth(login, password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(file, io.TeeReader(resp.Body, counter)); err != nil {
		file.Close()
		return err
	}

	// Close the file without defer so it can happen before Rename()
	file.Close()

	//Break Line after closing the counter
	fmt.Println("")
	// Remove .tmp sulfix from file after the download is endded
	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}

func uploadMailbox(filepath, admin, account, password, host, port string) error {
	// Use admin account to log in if it's setted
	var login string
	if admin == "" {
		login = account
	} else {
		login = admin
	}
	uri := fmt.Sprintf("https://%s:%s/home/%s/?fmt=tgz", host, port, account)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	// Create the request using TeeReader() to send the status to the status counter
	counter := &WriteCounter{}
	req, err := http.NewRequest("POST", uri, io.TeeReader(file, counter))
	if err != nil {
		return err
	}
	req.SetBasicAuth(login, password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//Break Line after closing the counter
	fmt.Println("")
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	return nil
}
