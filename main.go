package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	successfulCode = 200
	timeoutCode    = 408
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

func createSummary(link string, timeTaken time.Duration, totalTransferred int64, totalTimeAllRequests time.Duration, summary summaryInfo) string {
	u, _ := url.Parse(link)

	summary.DocumentPath = link
	summary.DocumentLength = len(link)
	summary.Hostname = u.Hostname()
	summary.Port = u.Port()
	summary.TimeTaken = timeTaken
	summary.CompletedRequests = summary.Responded
	summary.FailedRequests = summary.Requested - summary.Responded
	summary.TotalTransferred = totalTransferred
	summary.TimePerRequest = time.Duration(int64(totalTimeAllRequests) / summary.Responded)

	sortedSummary, _ := json.MarshalIndent(summary, "", "\t")
	formattedSummary := string(sortedSummary)

	return formattedSummary
}

func main() {
	fmt.Println("Hello from my app")

	requests := flag.Int64("n", 1, "Number of requests to perform")
	concurrency := flag.Int64("c", 1, "Number of multiple requests to make at a time")
	timeout := flag.Int64("s", 1000, "Seconds to max. wait for each response")
	timelimit := flag.Int64("t", 10000, "Seconds to max. to spend on benchmarking")
	fmt.Println(requests, concurrency, timeout, timelimit)

	flag.Parse()
	if flag.NArg() == 0 || *requests == 0 || *requests < *concurrency {
		flag.PrintDefaults()
		os.Exit(-1)
	}

	start := time.Now()
	summary := summaryInfo{}
	c := make(chan responseInfo)
	link := flag.Arg(0)

	for i := int64(0); i < *concurrency; i++ {
		go checkLink(link, c)
		summary.Requested++
	}

	totalTransferred := int64(0)
	totalTimeAllRequests := time.Duration(0)
	for response := range c {
		if summary.Requested < *requests {
			go checkLink(link, c)
			summary.Requested++
		}
		totalTransferred += response.bytes
		totalTimeAllRequests += response.duration
		fmt.Println(response)
		summary.Responded++

		if summary.Requested == summary.Responded {
			timeTaken := time.Now().Sub(start)
			formattedSummary := createSummary(link, timeTaken, totalTransferred, totalTimeAllRequests, summary)
			fmt.Println(formattedSummary)
			break
		}
	}
}
