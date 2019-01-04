package main

import "time"

type mktoAuth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type mktoExportRequest struct {
	Filter struct {
		CreatedAt struct {
			StartAt string `json:"startAt"`
			EndAt   string `json:"endAt"`
		} `json:"createdAt"`
	} `json:"filter"`
}

type mktoExportResult struct {
	RequestID string `json:"requestId"`
	Result    []struct {
		ExportID  string    `json:"exportId"`
		Format    string    `json:"format"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"createdAt"`
	} `json:"result"`
	Success bool `json:"success"`
}

type mktoEnqueueResult struct {
	RequestID string `json:"requestId"`
	Result    []struct {
		ExportID  string    `json:"exportId"`
		Format    string    `json:"format"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"createdAt"`
		QueuedAt  time.Time `json:"queuedAt"`
	} `json:"result"`
	Success bool `json:"success"`
}

type mktoJobStatus struct {
	RequestID string `json:"requestId"`
	Result    []struct {
		ExportID        string    `json:"exportId"`
		Format          string    `json:"format"`
		Status          string    `json:"status"`
		CreatedAt       time.Time `json:"createdAt"`
		QueuedAt        time.Time `json:"queuedAt"`
		StartedAt       time.Time `json:"startedAt"`
		FinishedAt      time.Time `json:"finishedAt"`
		NumberOfRecords int       `json:"numberOfRecords"`
		FileSize        int       `json:"fileSize"`
	} `json:"result"`
	Success bool `json:"success"`
}
