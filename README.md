# ab

A Go implementation of the [Apache HTTP server benchmarking tool](https://httpd.apache.org/docs/2.4/programs/ab.html). Stress-test web servers with concurrent HTTP requests and get a performance summary.

## Installation

```bash
git clone https://github.com/buihdk/ab.git
cd ab
go build -o ab
```

## Usage

```bash
./ab [options] <url>
```

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `-n` | 1 | Number of requests to perform |
| `-c` | 1 | Number of concurrent requests at a time |
| `-s` | 30 | Timeout in seconds per request |
| `-t` | 300 | Maximum seconds to spend benchmarking |

### Example

```bash
./ab -n 100 -c 20 -s 1 -t 10 https://example.com
```

## Output

Each response is printed as it completes. When all requests finish, a JSON summary is printed:

```json
{
    "Hostname": "example.com",
    "Port": "443",
    "DocumentPath": "https://example.com",
    "DocumentLength": 19,
    "ConcurrencyLevel": 20,
    "TimeTaken": 3200000000,
    "CompletedRequests": 100,
    "FailedRequests": 0,
    "TotalTransferred": 153600,
    "Rps": 31,
    "TimePerRequest": 32000000,
    "TransferRate": 46
}
```

`TransferRate` is in KB/s. `TimeTaken` and `TimePerRequest` are in nanoseconds.
