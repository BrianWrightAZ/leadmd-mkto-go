package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var baseURL *string
var clientID *string
var clientSecret *string
var resource *string
var startDate *string
var endDate *string
var authToken = ""

func main() {
	sleeper := func() {
		counter := 0
		for counter < 30 {
			time.Sleep(1 * time.Second)
			fmt.Print(".")
			counter++
		}
	}

	fmt.Println("Welcome to the LeadMD Marketo Bulk Extract Utility")

	setup()

	//Authenticate to Marketo
	token, authErr := authenticate()
	if authErr != nil {
		log.Panicf("Could not Authenticate. %v", authErr.Error())
	}
	fmt.Printf("Authentication success. token: %v expires in %v \n", token.AccessToken, token.ExpiresIn)
	authToken = token.AccessToken

	//Create the export
	export, exportErr := createExport()
	if exportErr != nil || export.Success == false {
		log.Panicf("Error creating export %v %v", exportErr.Error(), export.Result[0].Status)
	}
	exportID := export.Result[0].ExportID
	fmt.Printf("Export %v is ready to be queued \n", exportID)

	//Enqueue the job
	enqueue, enqueueErr := enqueueExport(exportID)
	if enqueue.Success == false || enqueueErr != nil {
		log.Panicf("Error enqueueing job: %v %v %v", exportID, enqueueErr.Error(), enqueue.Result[0].Status)
	}
	fmt.Printf("Job %v is has been queued \n", exportID)

	c := make(chan *mktoJobStatus)
	//Check job status
	go func(chan *mktoJobStatus) {
		jobReady := false
		for jobReady == false {
			fmt.Println("Checking job status in 30 seconds")
			sleeper()
			status, jobErr := checkJobStatus(export.Result[0].ExportID)
			if jobErr != nil {
				log.Panic(jobErr)
			}
			fmt.Printf("\n Job %v is %v \n", status.Result[0].ExportID, status.Result[0].Status)

			if status.Result[0].Status == "Completed" {
				jobReady = true
				c <- status
			}
		}
	}(c)

	jobStatus := <-c
	fmt.Printf("Job %v stats: records: %v size %vMb \n", jobStatus.Result[0].ExportID, jobStatus.Result[0].NumberOfRecords, jobStatus.Result[0].FileSize/1000000)

	//download the file
	fmt.Println("Starting download")
	download(export.Result[0].ExportID)
	fmt.Println("Download Complete")
}

func setup() {
	baseURL = flag.String("endpoint", "", "Marketo REST Endpoint")
	clientID = flag.String("client_id", "", "Marketo Client ID Endpoint")
	clientSecret = flag.String("client_secret", "", "Marketo Client Secret Endpoint")
	resource = flag.String("resource", "", "leads or actvities")
	startDate = flag.String("start", "", "Start Date")
	endDate = flag.String("end", "", "End Date")
	flag.Parse()

	if *baseURL == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *clientID == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *clientSecret == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *resource != "leads" && *resource != "activities" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *startDate == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *endDate == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}
