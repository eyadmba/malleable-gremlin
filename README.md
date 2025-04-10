# Malleable Gremlin

A versatile Go HTTP server for infrastructure testing, diagnostics, and load generation.

## 🤔 Why Malleable Gremlin?

When deploying services, testing infrastructure setups, or debugging network connectivity, you often need a simple, adaptable HTTP endpoint. Malleable Gremlin provides a configurable HTTP server with various capabilities designed to help:

* Verify network reachability (HTTP forwarding).
* Check system and network information on the host where it runs.
* Test database connections and run queries (PostgreSQL).
* Simulate basic CPU, memory, or I/O load.
* Simply echo back HTTP requests for debugging proxies or clients.

It aims to be a flexible "gremlin" you can deploy into your infrastructure to probe and test its behavior.

## ✨ Features

* **System & Network Info:** Endpoints to retrieve host system (CPU, memory) and network interface details.
* **Request Echoing:** Echo back details of incoming GET/POST requests (headers, query args, body, files), optionally overriding the response status code.
* **HTTP Forwarding:** Act as a simple HTTP proxy, forwarding requests to specified URLs.
* **PostgreSQL Interaction:** Test connectivity, store connection strings, and execute queries against a PostgreSQL database.
* **Load Generation:** Simulate CPU, memory, or I/O load with configurable parameters.
* **Easy Deployment:** Single Go binary or Docker container.

## 🚀 Getting Started

### Prerequisites

* **Go:** Version 1.22 or later (required for path parameter matching like `/path/{param...}`).
* **(Optional) Docker:** If you prefer containerized deployment.

### Installation

**1. From Source:**

```bash
# Clone the repository
git clone [https://github.com/your_github_username/malleable-gremlin.git](https://github.com/your_github_username/malleable-gremlin.git)
cd malleable-gremlin

# Build the binary
go build -o malleable-gremlin ./cmd/malleable-gremlin/
```


# Malleable Gremlin

A flexible HTTP server for testing infrastructure setups and networking. This server provides various endpoints for testing different aspects of a system, including system information, network interfaces, HTTP request forwarding, PostgreSQL database operations, and load testing.

## Features

- System Information Endpoints
- Network Interface Information
- HTTP Request Echoing
- HTTP Request Forwarding
- PostgreSQL Database Operations
- Load Testing (CPU, Memory, I/O)

## Installation

```bash
go mod download
```

## Running the Server

```bash
go run main.go
```

By default, the server runs on port 8080. You can specify a different address using the `-addr` flag:

```bash
go run main.go -addr :3000
```

## API Endpoints

### About Service

#### GET /about/system
Returns general information about the host's hardware, including:
- CPU information
- Memory usage
- Disk usage
- System information

#### GET /about/network
Returns information about the network interfaces of the host, including:
- Interface names
- IP addresses
- MAC addresses
- Interface flags
- MTU values

### Echo Service

#### GET /echo/get
Echoes back the request details. Returns a JSON response with:
- Query parameters
- Headers
- URL

You can control the response status code using the `status` query parameter:
```
GET /echo/get?status=400
```

#### POST /echo/post
Similar to GET /echo/get, but also includes:
- Form data
- File uploads
- Raw request body
- JSON body (if applicable)

### Load Service

#### GET /load/cpu
Generates CPU load by running multiple goroutines.
Parameters:
- `tasks`: Number of goroutines to run (or "cpus" to use all available CPUs)
- `timeout`: Duration to keep the goroutines running (e.g., "1s", "500ms")

#### GET /load/memory
Allocates memory and optionally triggers garbage collection.
Parameters:
- `size`: Amount of memory to allocate (e.g., "500mb", "1.5gb")
- `gc_after`: When to trigger GC (duration string, "0" for immediate, "-1" for never)

#### GET /load/io
Generates I/O load by running multiple goroutines.
Parameters:
- `tasks`: Number of goroutines to run
- `wait`: Duration each goroutine will wait
- `parallel`: Number of goroutines to run in parallel

### HTTP Forwarding Service

#### GET /http/send/{domain}/{path}
Forwards a GET request to the specified domain and path.
Example:
```
GET /http/send/example.com/api/users
```

#### POST /http/send
Forwards a custom HTTP request. Request body:
```json
{
    "url": "http://example.com/api/users",
    "method": "POST",
    "headers": {
        "Content-Type": "application/json"
    },
    "body": {
        "name": "John Doe"
    }
}
```

### PostgreSQL Service

#### PUT /postgresql/connection-string
Stores a PostgreSQL connection string.
Request body:
```json
{
    "connection_string": "postgres://user:pass@localhost:5432/dbname"
}
```

#### POST /postgresql/connect
Tests a PostgreSQL connection.
Request body:
```json
{
    "connection_string": "postgres://user:pass@localhost:5432/dbname"
}
```
or
```json
{
    "connection_string_id": "stored_connection_id"
}
```

#### POST /postgresql/query
Executes PostgreSQL queries.
Request body:
```json
{
    "connection_string_id": "stored_connection_id",
    "queries": [
        {
            "query": "SELECT * FROM users",
            "connection_timeout": "1s"
        }
    ]
}
```

## Docker Support

The server is designed to work both on the host system and inside Docker containers. When running inside Docker:

- System information will reflect the container's resources
- Network information will show the container's network interfaces
- Load testing will affect the container's resources

## License

MIT 