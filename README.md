# ab
https://httpd.apache.org/docs/2.4/programs/ab.html

1. **Submitted by: KHOA BUI**
2. **Time spent: 3 hours**

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
* [ ] Extend input params with: 
   * [ ] Timeout - Seconds to max. wait for each response
   * [ ] Timelimit - Maximum number of seconds to spend for benchmarking
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
