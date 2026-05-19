package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

func checkLink(ctx context.Context, client *http.Client, link string, c chan responseInfo) {
	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		c <- responseInfo{status: -1, duration: time.Since(start)}
		return
	}
	res, err := client.Do(req)
	if err != nil {
		c <- responseInfo{status: -1, duration: time.Since(start)}
		return
	}
	defer res.Body.Close()
	read, _ := io.Copy(io.Discard, res.Body)
	c <- responseInfo{
		status:   res.StatusCode,
		bytes:    read,
		duration: time.Since(start),
	}
}

func createSummary(link string, timeTaken time.Duration, totalTransferred int64, totalTimeAllRequests time.Duration, concurrency int64, documentLength int64, summary summaryInfo) string {
	u, _ := url.Parse(link)

	seconds := timeTaken.Seconds()

	summary.DocumentPath = link
	summary.DocumentLength = int(documentLength)
	summary.Hostname = u.Hostname()
	summary.Port = u.Port()
	summary.ConcurrencyLevel = concurrency
	summary.TimeTaken = timeTaken
	summary.CompletedRequests = summary.Responded
	summary.FailedRequests = summary.Requested - summary.Responded
	summary.TotalTransferred = totalTransferred
	if summary.Responded > 0 {
		summary.TimePerRequest = time.Duration(int64(totalTimeAllRequests) / summary.Responded)
	}
	summary.Rps = int64(float64(summary.Responded) / seconds)
	summary.TransferRate = int64(float64(totalTransferred) / seconds / 1024)

	sortedSummary, _ := json.MarshalIndent(summary, "", "\t")
	return string(sortedSummary)
}
