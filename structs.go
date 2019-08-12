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
	Requested         int64
	Responded         int64
	Hostname          string
	Port              string
	DocumentPath      string
	DocumentLength    int
	ConcurrencyLevel  int64
	TimeTaken         time.Duration
	CompletedRequests int64
	FailedRequests    int64
	TotalTransferred  int64
	Rps               int64
	TimePerRequest    time.Duration
	TransferRate      int64
}