# ab
https://httpd.apache.org/docs/2.4/programs/ab.html

Run `./ab -n 100 -c 20 -s 1 -t 10 https://www.grab.com/vn`

## Set of user stories

### Required
* [x] Command-line argument parsing
* [x] Input params
   * [x] Requests - Number of requests to perform
   * [x] Concurrency - Number of multiple requests to make at a time
   * [x] URL - The URL for testing
* [x] The program prints usage information if the wrong arguments are provided.
* [x] The program performs the specified HTTP requests and prints a summary of the results.
* [x] The program’s concurrency is implemented with goroutines.


### Bonus
* [x] Extend input params with: 
   * [x] Timeout - Seconds to max. wait for each response
   * [x] Timelimit - Maximum number of seconds to spend for benchmarking
* [ ] Prints key metrics of summary, such:
   * [x] Server Hostname
   * [x] Server Port
   * [x] Document Path
   * [x] Document Length
   * [ ] Concurrency Level
   * [x] Time taken for tests
   * [x] Complete requests
   * [x] Failed requests
   * [x] Total transferred
   * [ ] Requests per second
   * [x] Time per request
   * [ ] Transfer rate
