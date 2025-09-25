# Go Internet Speed Test

A simple Go application to test your internet speed via CLI or as an HTTP API server.

## Features

- **CLI Mode:** Run a speed test directly from your terminal.
- **API Mode:** Expose a `/speedtest` HTTP endpoint for remote speed testing.
- **ISP Info:** Fetches and displays your ISP information.
- **Automatic Best Server Selection:** Finds the lowest-latency server for accurate results.

## Requirements

- Go 1.18+
- Internet connection

## Installation

Clone the repository:

```sh
git clone https://github.com/chrislim1914/go-internet-speed-test.git
cd go-internet-speed-test
```

Install dependencies:

```sh
go mod tidy
```

## Usage

### CLI Mode

Run a speed test from your terminal:

```sh
go run main.go
```

### API Mode

Start the HTTP API server (default port 8080):

```sh
go run main.go --api
```

Specify a custom port:

```sh
go run main.go --api --port=9090
```

Then, access the speed test endpoint:

```
GET http://localhost:8080/speedtest
```

## Output

- **CLI:** Prints best server, ISP info, download and upload speeds.
- **API:** Returns download and upload speeds as plain text.

## Project Structure

- `main.go` - Entry point, handles CLI/API mode.
- `speedtester/` - Speed test logic and server selection.
- `utilities/` - Helper functions (spinner, URL normalization, etc).

## License

MIT

---

*Created by [chrislim1914](https://github.com/chrislim1914)*