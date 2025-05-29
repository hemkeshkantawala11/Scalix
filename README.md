# Key-Value Cache System (Scalix) 🚀

Hey there! Welcome to my distributed cache system project. This is a Redis-inspired, high-performance in-memory cache built with Go. If you've ever wondered how distributed caching works or wanted to see a practical implementation of consistent hashing, you're in the right place!

## What Is This Project? 🤔

In simple terms, this is a distributed memory cache that spreads data across multiple nodes. Instead of storing everything on one server (which can become a bottleneck), this system cleverly distributes your data using something called "consistent hashing" (more on that below).

Think of it like having multiple warehouses instead of one giant warehouse - it's more efficient and less prone to failure.

## Key Features 💎

- **Truly Distributed**: Data is spread across multiple nodes intelligently
- **Blazing Fast**: Written in Go for maximum performance 
- **Memory Smart**: Uses LRU (Least Recently Used) policy to manage memory
- **Self-Monitoring**: Automatically watches memory usage and cleans up when needed
- **Simple API**: Easy-to-use HTTP interface for all operations
- **Flexible**: Add or remove nodes whenever you want without disruption
- **Docker Ready**: Deploy anywhere with optimized container setup

## How It Works 🔍

### The Architecture

This project is built with three main components:

1. **Cache Layer** (The Storage Engine)
   - Data is stored in "shards" across multiple nodes
   - We use an LRU cache to efficiently manage memory
   - A background process constantly monitors memory usage
   - Go's concurrency features (channels, sync.Map) make it incredibly fast

2. **Consistent Hashing** (The Secret Sauce)
   - This is the magic that decides which key goes to which node
   - When nodes are added or removed, only a minimal number of keys need to be reshuffled
   - We use virtual nodes to ensure even distribution

3. **API Layer** (The Interface)
   - A RESTful API built with the Gin framework
   - Simple endpoints for getting, setting, and managing nodes
   - Optimized for high throughput and minimal latency

## Technical Deep Dive 🧠

### Smart Optimizations

1. **Non-blocking Write Operations**
   - Instead of blocking when writing data, we use a buffered channel (queue)
   - Multiple worker goroutines process this queue concurrently
   - This approach significantly improves performance during high write loads

2. **Proactive Memory Management**
   - Every 5 seconds, a background goroutine checks memory usage
   - If memory usage exceeds 70%, the LRU cache is purged
   - This prevents those dreaded "out of memory" crashes

3. **Thread-Safe Data Access**
   - Uses Go's `sync.Map` for highly concurrent read/write operations
   - Much more efficient than traditional locks for read-heavy workloads

4. **Configurable Cache Limits**
   - You can set the maximum cache size to match your hardware
   - Prevents unbounded memory growth that could crash your application

## Performance Highlights 📊

I've included a Locust script (`locust.py`) for load testing that shows:

- Handles thousands of requests per second on modest hardware
- Response times typically under 5ms
- Read operations (GET) are prioritized 3:2 over writes (PUT)

## Let's Get Started! 🏁

### Prerequisites

Before you begin, make sure you have:

- Docker installed on your system (for containerized deployment)
- OR Go 1.18+ installed (for local development)
- Basic familiarity with terminal/command line

### Option 1: Running with Docker (Recommended) 🐳

Docker makes it super easy to run this project without worrying about dependencies. Let me walk you through it:

#### Step 1: Clone the repository

```bash
git clone https://github.com/yourusername/HLD-Redis-Assignment.git
cd HLD-Redis-Assignment
```

#### Step 2: Pull the public Docker image

```bash
docker pull hemkesh11/hld-redis-assignment:latest
```

This command pulls the Docker image using the `Dockerfile` in the project. Here's what this Dockerfile does:

```dockerfile
# Build Stage - This stage compiles the Go application
FROM golang:1.23.1 AS builder

WORKDIR /app

# Copy go.mod and go.sum first for efficient caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Change directory to where the main.go file exists
WORKDIR /app/cmds/server

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

# Final Stage - This stage creates the actual runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /main .

# Ensure it's executable
RUN chmod +x /root/main

EXPOSE 7171

CMD ["/root/main"]
```

The Dockerfile uses a multi-stage build approach:
- First stage: Compiles the Go application
- Second stage: Creates a minimal Alpine Linux image with just the compiled binary
- This makes the final image very small and secure

#### Step 3: Run the container

```bash
docker run -p 7171:7171 hemkesh11/hld-redis-assignment:latest
```

This command:
- Starts a container from the image we just built
- Maps port 7171 from the container to port 7171 on your host
- The `-p 7171:7171` part means "connect port 7171 inside the container to port 7171 on my computer"

You should see output indicating the server has started:

```
[GIN] 2023/01/01 - 12:00:00 | 200 | 123.456µs | 127.0.0.1 | GET "/get"
Starting server on :7171
```

#### Step 4: Verify it's working

Open a new terminal window and run:

```bash
curl http://127.0.0.1:7171/get?key=test
```

You should get a response indicating the key wasn't found (which is expected since we haven't added any keys yet).

## Using the Distributed Cache 💻

Now that your server is running, let's try some operations:

### 1. Setting Values in the Cache

To add a key-value pair to the cache:

```bash
curl -X POST http://127.0.0.1:7171/put \
  -H "Content-Type: application/json" \
  -d '{"Key": "user:1234", "Value": "John Doe"}'
```

You should receive:
```json
{
  "status": "OK",
  "message": "Key inserted/updated successfully."
}

```

### 2. Retrieving Values from the Cache

To retrieve a value by its key:

```bash
curl http://127.0.0.1:7171/get?key=user:1234
```

You should receive:
```json
{
  "status": "OK",
  "key": "exampleKey",
  "value": "the corresponding value"
}

```

## Load Testing Your Cache ⚡

Want to see how your cache performs under load? I've included a Locust script for stress testing:

### Step 1: Install Locust

```bash
pip install locust
```

### Step 2: Run the load test

```bash
locust -f locust.py --host=http://127.0.0.1:7171
```

### Step 3: Access the Locust web interface

Open your browser and go to http://localhost:8089

### Step 4: Configure and start the test

1. Enter the number of users to simulate (e.g., 100)
2. Enter the spawn rate (users started/second, e.g., 10)
3. Click "Start swarming"

You'll see real-time graphs and statistics showing how your cache is performing!

## Understanding the Code Structure 📁

If you're curious about the code organization:

- `cmds/server/main.go` - The entry point of the application
- `internal/cache/cache.go` - The core cache implementation
- `internal/consistentHash/consistentHash.go` - Consistent hashing algorithm
- `internal/cache_handlers/cache_handlers.go` - API handlers
- `internal/models/request_models.go` - Request and response data structures
- `locust.py` - Load testing script
- `Dockerfile` - Container definition for Docker

## Best Practices I've Implemented 🌟

Throughout this project, I've followed several best practices:

1. **Clean Architecture**
   - Separation of concerns with distinct layers
   - Each component does one thing and does it well
   - Easy to understand, test, and maintain

2. **Thorough Error Handling**
   - Every operation has appropriate error checks
   - HTTP responses include helpful status codes
   - Potential failures are gracefully managed

3. **Proper Resource Management**
   - All locks are properly acquired and released
   - Channels are correctly managed to prevent leaks
   - Memory usage is constantly monitored

4. **Efficient Docker Configuration**
   - Multi-stage builds for minimal image size
   - Proper layer caching for faster builds
   - Secure base image (Alpine Linux)

5. **Performance Testing**
   - Included load testing with realistic patterns
   - Both read and write operations are tested
   - Tests simulate real-world usage patterns

## Ideas for Future Enhancements 💡

Here are some ways this project could be extended:

- **Data Persistence**: Add a way to save data to disk so it survives restarts
- **Time-To-Live (TTL)**: Make cache entries expire after a certain time
- **Security Layer**: Add authentication and authorization
- **Circuit Breaker**: Implement fault tolerance patterns
- **Metrics Dashboard**: Add Prometheus and Grafana for monitoring
- **Auto-scaling**: Dynamically add/remove nodes based on load

## Contact and Support 📞

Found a bug or have a question? Please open an issue on GitHub or reach out directly.

---

Thanks for checking out my distributed cache system! I hope it helps you understand how distributed systems work or serves as a useful component in your architecture. Feel free to use, modify, and learn from this code! 
