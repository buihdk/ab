package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func checkLink(link string, c chan responseInfo) {
	start := time.Now()
	res, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	read, _ := io.Copy(ioutil.Discard, res.Body)
	c <- responseInfo{
		status:   res.StatusCode,
		bytes:    read,
		duration: time.Now().Sub(start),
	}
}

func createSummary(link string, timeTaken time.Duration, totalTransferred int64, totalTimeAllRequests time.Duration, concurrency int64, summary summaryInfo) string {
	u, _ := url.Parse(link)

	seconds := timeTaken.Seconds()

	summary.DocumentPath = link
	summary.DocumentLength = len(link)
	summary.Hostname = u.Hostname()
	summary.Port = u.Port()
	summary.ConcurrencyLevel = concurrency
	summary.TimeTaken = timeTaken
	summary.CompletedRequests = summary.Responded
	summary.FailedRequests = summary.Requested - summary.Responded
	summary.TotalTransferred = totalTransferred
	summary.TimePerRequest = time.Duration(int64(totalTimeAllRequests) / summary.Responded)
	summary.Rps = int64(float64(summary.Responded) / seconds)
	summary.TransferRate = int64(float64(totalTransferred) / seconds / 1024)

	sortedSummary, _ := json.MarshalIndent(summary, "", "\t")
	formattedSummary := string(sortedSummary)

	return formattedSummary
}
