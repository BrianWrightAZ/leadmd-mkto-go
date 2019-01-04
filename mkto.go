package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func authenticate() (*mktoAuth, error) {

	url := *baseURL + "/identity/oauth/token?grant_type=client_credentials&client_id=" + *clientID + "&client_secret=" + *clientSecret

	req, _ := http.NewRequest("GET", url, nil)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	auth := mktoAuth{}
	var err error

	if res.StatusCode == http.StatusOK {

		json.NewDecoder(res.Body).Decode(&auth)
	} else {
		err = errors.New("Authentication error")
	}

	return &auth, err
}

func createExport() (*mktoExportResult, error) {

	url := *baseURL + "/bulk/v1/" + *resource + "/export/create.json"

	export := mktoExportRequest{}
	export.Filter.CreatedAt.StartAt = *startDate
	export.Filter.CreatedAt.EndAt = *endDate

	body := new(bytes.Buffer)

	json.NewEncoder(body).Encode(export)

	req, _ := http.NewRequest("POST", url, body)
	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	exportResult := mktoExportResult{}
	var err error

	if res.StatusCode == http.StatusOK {
		json.NewDecoder(res.Body).Decode(&exportResult)
	} else {
		err = errors.New("Export creation error")
	}

	return &exportResult, err
}

func enqueueExport(ExportID string) (*mktoEnqueueResult, error) {

	url := *baseURL + "/bulk/v1/" + *resource + "/export/" + ExportID + "/enqueue.json"

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	enqueueResult := mktoEnqueueResult{}
	var err error

	if res.StatusCode == http.StatusOK {
		json.NewDecoder(res.Body).Decode(&enqueueResult)
	} else {
		err = errors.New("Enqueuing  error")
	}

	return &enqueueResult, err
}

func checkJobStatus(ExportID string) (*mktoJobStatus, error) {
	url := *baseURL + "/bulk/v1/" + *resource + "/export/" + ExportID + "/status.json"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	jobStatus := mktoJobStatus{}
	var err error

	if res.StatusCode == http.StatusOK {
		json.NewDecoder(res.Body).Decode(&jobStatus)
	} else {
		err = errors.New("Job Status  error")
	}

	return &jobStatus, err
}

func download(ExportID string) error {
	url := *baseURL + "/bulk/v1/" + *resource + "/export/" + ExportID + "/file.json"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	if res.StatusCode == http.StatusOK {
		defer res.Body.Close()
		out, err := os.Create(*resource + "_" + ExportID + ".csv")
		if err != nil {
			return err
		}
		defer out.Close()
		io.Copy(out, res.Body)

	}

	return nil
}
