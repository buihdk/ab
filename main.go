package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	requests := flag.Int64("n", 1, "Number of requests to perform")
	concurrency := flag.Int64("c", 1, "Number of multiple requests to make at a time")
	timeout := flag.Int64("s", 30, "Seconds to max. wait for each response")
	timelimit := flag.Int64("t", 300, "Maximum seconds to spend benchmarking")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Invalid number of arguments!")
		fmt.Println("Usage: ab [options] [http[s]://]hostname[:port]/path")
		fmt.Println("Example: ./ab -n 100 -c 20 -s 1 -t 10 https://www.grab.com/vn")
		flag.PrintDefaults()
		os.Exit(-1)
	}

	if *requests == 0 {
		fmt.Println("Number of requests cannot be zero!")
		os.Exit(-1)
	}

	if *requests < *concurrency {
		fmt.Println("Number of requests cannot be less than number of concurrency!")
		os.Exit(-1)
	}

	client := &http.Client{Timeout: time.Duration(*timeout) * time.Second}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timelimit)*time.Second)
	defer cancel()

	start := time.Now()
	summary := summaryInfo{}
	c := make(chan responseInfo)
	link := flag.Arg(0)

	for i := int64(0); i < *concurrency; i++ {
		go checkLink(ctx, client, link, c)
		summary.Requested++
	}

	documentLength := int64(0)
	totalTransferred := int64(0)
	totalTimeAllRequests := time.Duration(0)
	for response := range c {
		if summary.Requested < *requests {
			go checkLink(ctx, client, link, c)
			summary.Requested++
		}
		if documentLength == 0 && response.status == successfulCode {
			documentLength = response.bytes
		}
		totalTransferred += response.bytes
		totalTimeAllRequests += response.duration
		fmt.Println(response)
		summary.Responded++

		if summary.Requested == summary.Responded {
			timeTaken := time.Since(start)
			formattedSummary := createSummary(link, timeTaken, totalTransferred, totalTimeAllRequests, *concurrency, documentLength, summary)
			fmt.Println(formattedSummary)
			break
		}
	}
}
