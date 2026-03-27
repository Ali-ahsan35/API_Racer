# apiracer

A Beego-based Go project that demonstrates and benchmarks **concurrent vs sequential API calling** using goroutines, channels, and WaitGroup. Built as a learning exercise to deeply understand Go concurrency patterns.

---

## What This Project Covers

- Goroutines
- Buffered and unbuffered channels
- `sync.WaitGroup`
- `sync.Mutex`
- Sequential vs concurrent execution
- Synchronization between goroutines
- Response collection from multiple goroutines
- Error handling in concurrent processes
- Execution time measurement
- Performance comparison and terminal visualization

---

## Project Structure

```
apiracer/
├── conf/                          # Beego configuration
├── controllers/
│   ├── benchmark_controller.go    # Benchmark endpoint controller
│   └── default.go                 # Default Beego controller
├── models/                        # Models (reserved for future use)
├── request/
│   └── api_request.go             # Reusable HTTP request layer
├── routers/
│   └── router.go                  # Route definitions
├── service/
│   ├── sequential.go              # Phase 1: Sequential execution
│   ├── waitgroup.go               # Phase 2: Concurrent with WaitGroup
│   └── channel.go                 # Phase 3: Concurrent with Channels
├── utils/
│   └── visualizer.go              # Terminal output and comparison
├── static/                        # Static files
├── tests/                         # Test files
├── views/                         # Beego views
├── main.go                        # Entry point
├── go.mod
└── go.sum
```

---

## Prerequisites

- Go `1.25+`
- Beego `v2.1.0`
- Bee CLI tool

---

## Installation & Setup

### 1. Clone the repository

```bash
git clone https://github.com/Ali-ahsan35/API_Racer
cd API_Racer
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Install Bee CLI (if not already installed)

```bash
go install github.com/beego/bee/v2@latest
```

### 4. Run the project

```bash
bee run
```

Server starts at: `http://localhost:8080`

---

## API Endpoint

| Method | URL | Description |
|--------|-----|-------------|
| GET | `/benchmark` | Runs all 3 execution strategies and prints comparison |

### Hit the endpoint

```bash
curl http://localhost:8080/benchmark
```

Or open in browser / Postman:
```
GET http://localhost:8080/benchmark
```

---

## How It Works

### Phase 1 — Sequential Execution
All 12 external APIs are called **one by one** in a simple loop. Each API waits for the previous one to finish before starting.

```
API 1 ──▶ API 2 ──▶ API 3 ──▶ ... ──▶ API 12
Total time = sum of all individual times
```

### Phase 2 — Concurrent with WaitGroup
All 12 APIs are called **simultaneously** using goroutines. `sync.WaitGroup` waits for all goroutines to finish. `sync.Mutex` protects the shared success counter.

```
API 1  ──▶ (goroutine)
API 2  ──▶ (goroutine)
...          all run at the same time
API 12 ──▶ (goroutine)
           WaitGroup.Wait() ──▶ done
Total time ≈ slowest single API
```

### Phase 3 — Concurrent with Channels
All 12 APIs are called **simultaneously** using goroutines. Each goroutine sends its result (`true`/`false`) into a **buffered channel**. Main function collects all 12 results from the channel.

```
API 1  ──▶ goroutine ──▶ ch <- true
API 2  ──▶ goroutine ──▶ ch <- true
...
API 12 ──▶ goroutine ──▶ ch <- true
           collect 12 results from channel
```

---

## Sample Terminal Output

```
[1] Running Sequential Execution...
-----------------------------------
  [API 1] Success
  [API 2] Success
  ...
  [API 12] Success

[2] Running Concurrent (WaitGroup)...
-----------------------------------
  [API 9] Success
  [API 3] Success
  ...

[3] Running Concurrent (Channels)...
-----------------------------------
  [API 7] Success
  [API 11] Success
  ...

================= API PERFORMANCE TEST =================

Total APIs Called: 12

[1] Sequential Execution:
-----------------------------------
Time Taken : 4679 ms
Success    : 12/12

[2] Concurrent (WaitGroup):
-----------------------------------
Time Taken : 596 ms
Success    : 12/12

[3] Concurrent (Channels):
-----------------------------------
Time Taken : 373 ms
Success    : 12/12

================= COMPARISON =================

Performance Gain:
- WaitGroup vs Sequential  : ~87% faster
- Channels vs Sequential   : ~92% faster
- WaitGroup vs Channels    : ~60% slower than Channels

=======================================================
```

---

## Key Observations

**1. Random order in concurrent execution**
> Sequential prints API 1, 2, 3... in order. WaitGroup and Channels print in random order. This proves goroutines truly run simultaneously — whichever finishes first prints first.

**2. Why Channels is faster than WaitGroup**
> WaitGroup requires a `sync.Mutex` lock for every result update — extra overhead. Channels have built-in synchronization — no mutex needed, making it slightly faster.

**3. Why Sequential is slowest**
> Each of the 12 APIs waits for the previous one. Total time = sum of all API response times (~4679ms).

**4. Why Concurrent is fastest**
> All 12 APIs run at the same time. Total time ≈ slowest single API response time (~373ms).

---

## Tech Stack

| Technology | Version | Purpose |
|---|---|---|
| Go | 1.25 | Core language |
| Beego | v2.1.0 | Web framework |
| Bee CLI | latest | Development tool |
| `sync.WaitGroup` | stdlib | Goroutine synchronization |
| `sync.Mutex` | stdlib | Shared data protection |
| `net/http` | stdlib | HTTP client |

---

## Author

Built as an internship learning task to understand Go concurrency patterns — goroutines, channels, WaitGroup, and performance benchmarking.
