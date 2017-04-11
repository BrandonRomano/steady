package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type RemoteFiles struct {
	Files []RemoteFile `json:"files"`
}

type RemoteFile struct {
	Url  string `json:"url"`
	Type string `json:"type"`
}

const (
	outputDir       = "public"
	steadyFilesPath = "/steady-files"
	numWorkers      = 5
)

func main() {
	// Get base URL
	if len(os.Args) < 2 {
		log.Fatal(errors.New("`steady` must be called passing the base URL of the site to convert"))
	}
	baseUrl := os.Args[1]

	// Load steady files
	resp, err := http.Get(baseUrl + steadyFilesPath)
	if err != nil {
		log.Fatal(err)
	}
	steadyFileBytes, err := ioutil.ReadAll(resp.Body)

	// Unmarshal Files
	var remoteFiles RemoteFiles
	err = json.Unmarshal(steadyFileBytes, &remoteFiles)
	if err != nil {
		log.Fatal(err)
	}

	// Download the files
	err = downloadFiles(&remoteFiles, baseUrl)
	if err != nil {
		log.Fatal(err)
	}
}

// downloadFiles downloads all remote files from the specified remote
func downloadFiles(remoteFiles *RemoteFiles, baseUrl string) error {
	// Create channels
	numJobs := len(remoteFiles.Files)
	jobs := make(chan RemoteFile, numJobs)
	results := make(chan error, numJobs)

	// Spin up some workers
	for workerId := 0; workerId < numWorkers; workerId++ {
		go downloadFilesWorker(baseUrl, jobs, results)
	}

	// Create the jobs
	for _, remoteFile := range remoteFiles.Files {
		jobs <- remoteFile
	}
	close(jobs)

	// Block until we get all results
	for result := 0; result < numJobs; result++ {
		err := <-results
		if err != nil {
			return err
		}
	}

	// OK
	return nil
}

// downloadFilesWorker is a single worker that will call downloadFile on a job
func downloadFilesWorker(baseUrl string, jobs <-chan RemoteFile, results chan<- error) {
	for remoteFile := range jobs {
		results <- downloadFile(&remoteFile, baseUrl)
	}
}

// downloadFile downloads a single remote file
func downloadFile(remoteFile *RemoteFile, baseUrl string) error {
	// Download
	downloadUrl := baseUrl + remoteFile.Url
	fmt.Println("======> Downloading " + downloadUrl)
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}

	// Get local file URL
	localFileUrl := remoteFile.Url
	if remoteFile.Type != "" {
		localFileUrl += "." + remoteFile.Type
	}
	localFileUrl = outputDir + localFileUrl

	// Ensure all directories are made
	fileSegmentArray := strings.Split(localFileUrl, "/")
	dirPath := ""
	for _, segment := range fileSegmentArray[0 : len(fileSegmentArray)-1] {
		dirPath += segment + string(filepath.Separator)
	}
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}

	// Create local file
	localFile, err := os.Create(localFileUrl)
	if err != nil {
		return err
	}

	// Write the response body to the file
	_, err = io.Copy(localFile, resp.Body)
	return err
}
