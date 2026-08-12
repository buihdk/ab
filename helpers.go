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
		c <- responseInfo{status: errorCode, duration: time.Since(start)}
		return
	}
	res, err := client.Do(req)
	if err != nil {
		c <- responseInfo{status: errorCode, duration: time.Since(start)}
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

	completed := summary.Responded - summary.Failed

	summary.DocumentPath = link
	summary.DocumentLength = int(documentLength)
	summary.Hostname = u.Hostname()
	summary.Port = u.Port()
	summary.ConcurrencyLevel = concurrency
	summary.TimeTaken = timeTaken
	summary.CompletedRequests = completed
	summary.FailedRequests = summary.Failed
	summary.TotalTransferred = totalTransferred
	if completed > 0 {
		summary.TimePerRequest = time.Duration(int64(totalTimeAllRequests) / completed)
	}
	if seconds > 0 {
		summary.Rps = float64(completed) / seconds
		summary.TransferRate = float64(totalTransferred) / seconds / 1024
	}

	sortedSummary, _ := json.MarshalIndent(summary, "", "\t")
	return string(sortedSummary)
}
