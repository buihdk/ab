package main

import (
	"time"
)

type responseInfo struct {
	status   int
	bytes    int64
	duration time.Duration
}

type summaryInfo struct {
	Requested         int64 `json:"-"`
	Responded         int64 `json:"-"`
	Failed            int64 `json:"-"`
	Hostname          string
	Port              string
	DocumentPath      string
	DocumentLength    int
	ConcurrencyLevel  int64
	TimeTaken         time.Duration
	CompletedRequests int64
	FailedRequests    int64
	Non2xxResponses   int64
	TotalTransferred  int64
	Rps               float64
	TimePerRequest    time.Duration
	TransferRate      float64
}