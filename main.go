package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	requests := flag.Int64("n", 1, "Number of requests to perform")
	concurrency := flag.Int64("c", 1, "Number of multiple requests to make at a time")
	timeout := flag.Int64("s", 30, "Seconds to max. wait for each response")
	timelimit := flag.Int64("t", 300, "Seconds to max. to spend on benchmarking")
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
		fmt.Println("Number of requests cannot be less than number of concurency!")
		os.Exit(-1)
	}

	fmt.Println("requests :", *requests, requests)
	fmt.Println("concurrency :", *concurrency, concurrency)
	fmt.Println("timeout :", *timeout, timeout)
	fmt.Println("timelimit :", *timelimit, timelimit)
	fmt.Println()

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
