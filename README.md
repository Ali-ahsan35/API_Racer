# API_Racer

A Beego-based Go project that demonstrates and benchmarks **concurrent vs sequential API calling** using goroutines, channels, and WaitGroup. Built as a learning exercise to deeply understand Go concurrency patterns and performance profiling.

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
- CPU profiling using `runtime/pprof`
- Memory profiling using `runtime.ReadMemStats()`
- Goroutine count tracking using `runtime.NumGoroutine()`

---

## Project Structure

```
apiracer/
├── conf/
│   └── app.conf                   # Beego configuration + API key
├── controllers/
│   ├── benchmark_controller.go    # Benchmark endpoint controller
│   └── default.go                 # Default Beego controller
├── models/                        # Models (reserved for future use)
├── request/
│   └── api_request.go             # Reusable HTTP request layer with response validation
├── routers/
│   └── router.go                  # Route definitions
├── service/
│   ├── sequential.go              # Phase 1: Sequential execution
│   ├── waitgroup.go               # Phase 2: Concurrent with WaitGroup
│   └── channel.go                 # Phase 3: Concurrent with Channels
├── utils/
│   ├── visualizer.go              # Terminal output and performance comparison
│   └── profiler.go                # Profiling logic (memory, goroutines, CPU)
├── static/                        # Static files
├── tests/                         # Test files
├── views/                         # Beego views
├── cpu.prof                       # CPU profile output (generated at runtime)
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

### 4. Add API Key

Open `conf/app.conf` and add your API key:

```ini
appname = apiracer
httpport = 8080
runmode = dev
apikey = your_secret_key_here
```

### 5. Run the project

```bash
bee run
```

Server starts at: `http://localhost:8080`

---

## API Endpoint

| Method | URL | Description |
|--------|-----|-------------|
| GET | `/benchmark` | Runs all 3 execution strategies, prints comparison and profiling report |

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
All 12 APIs are called **simultaneously** using goroutines. Each goroutine sends its result into a **buffered channel**. A separate goroutine closes the channel when all senders finish. Main function collects results using `range`.

```
API 1  ──▶ goroutine ──▶ ch <- true
API 2  ──▶ goroutine ──▶ ch <- true
...
API 12 ──▶ goroutine ──▶ ch <- true
           close(ch) when all done
           range ch collects all results
```

---

## Response Validation

Every API call goes through 3 levels of validation:

```
Level 1 → HTTP status code must be 200 OK
Level 2 → Response must be valid JSON
Level 3 → Response must contain GeoInfo and Result fields
```

---

## Profiling

### What is Profiled?

| Metric | Tool | What it measures |
|---|---|---|
| Memory | `runtime.ReadMemStats()` | Alloc, Sys, NumGC before and after |
| Goroutines | `runtime.NumGoroutine()` | Active goroutines at peak |
| CPU | `runtime/pprof` | CPU activity saved to cpu.prof |

### How Profiling Works

**Memory Profiling:**
```
Before running APIs  →  take memory snapshot
Run all 3 phases
After running APIs   →  take memory snapshot
Compare snapshots    →  see memory consumed per phase
```

**Goroutine Profiling:**
```
Sequential  →  5 goroutines  (background only, no goroutines created)
WaitGroup   →  19 goroutines (5 background + 12 yours)
Channels    →  20 goroutines (5 background + 12 yours + 1 channel closer)
```

**CPU Profiling:**
```
pprof.StartCPUProfile() →  start recording CPU activity
... all 3 phases run ...
pprof.StopCPUProfile()  →  stop recording
cpu.prof                →  saved recording file
```

### Analyze CPU Profile

After running the benchmark, analyze the CPU profile:

```bash
go tool pprof cpu.prof
```

Useful commands inside pprof:

```bash
top           # shows top CPU consuming functions
list FetchAPI # shows CPU usage line by line in FetchAPI
web           # opens visual graph in browser
exit          # quit pprof
```

---

## Sample Terminal Output

```
------------------ CPU PROFILE ------------------
  CPU Profiling Started...

[1] Running Sequential Execution...
-----------------------------------
  [API 1] Success | Location: Hawaii, USA
  [API 2] Success | Location: North America
  ...
  [API 12] Success | Location: Texas, USA

[2] Running Concurrent (WaitGroup)...
-----------------------------------
  [API 9]  Success | Location: Texas, USA
  [API 3]  Success | Location: North America
  ...

[3] Running Concurrent (Channels)...
-----------------------------------
  [API 7]  Success | Location: Texas, USA
  [API 11] Success | Location: Texas, USA
  ...

  Total Execution Time: 1908 ms

================= API PERFORMANCE TEST =================

Total APIs Called: 12

[1] Sequential Execution:
-----------------------------------
Time Taken : 1842 ms
Success    : 12/12

[2] Concurrent (WaitGroup):
-----------------------------------
Time Taken : 37 ms
Success    : 12/12

[3] Concurrent (Channels):
-----------------------------------
Time Taken : 25 ms
Success    : 12/12

================= COMPARISON =================

Performance Gain:
- WaitGroup vs Sequential  : ~98% faster
- Channels vs Sequential   : ~99% faster
- WaitGroup vs Channels    : ~48% slower than Channels

================= PROFILING REPORT =================

------------------ MEMORY STATS ------------------
Before Execution:
  Alloc = 2.20 MB | Sys = 12.21 MB | NumGC = 0

After Execution:
  Alloc = 5.44 MB | Sys = 21.33 MB | NumGC = 8

------------------ PHASE PROFILING ------------------

[1] Sequential Execution
  Time Taken     : 1842 ms
  Memory Used    : +0.36 MB
  Goroutines     : 5

[2] WaitGroup Execution
  Time Taken     : 37 ms
  Memory Used    : +2.94 MB
  Goroutines     : 19

[3] Channel Execution
  Time Taken     : 25 ms
  Memory Used    : GC freed 1.06 MB
  Goroutines     : 20

------------------ CPU PROFILE ------------------
  CPU Profiling Started...
  CPU Profiling Stopped.
  CPU Profile saved to: cpu.prof
  Run: go tool pprof cpu.prof

------------------ SUMMARY ------------------
  Fastest Method        : Channels
  Highest Memory Usage  : WaitGroup
  Most Efficient Method : Channels

=================================================
```

---

## Key Observations

**1. Random order in concurrent execution**
> Sequential prints API 1, 2, 3... in order. WaitGroup and Channels print in random order. This proves goroutines truly run simultaneously — whichever finishes first prints first.

**2. Why Channels is faster than WaitGroup**
> WaitGroup requires a `sync.Mutex` lock for every result update — extra overhead. Channels have built-in synchronization — no mutex needed, making it slightly faster.

**3. Why Sequential is slowest**
> Each of the 12 APIs waits for the previous one. Total time = sum of all API response times.

**4. Why Concurrent is fastest**
> All 12 APIs run at the same time. Total time ≈ slowest single API response time.

**5. Why Sequential uses least memory**
> Only 1 API runs at a time. No goroutine overhead. Memory stays low (+0.36 MB).

**6. Why WaitGroup uses most memory**
> 12 goroutines all alive at the same time, each holding their own response data. Peak memory usage is highest.

**7. Why Channels shows GC freed memory**
> Go's garbage collector ran during channel execution and freed more memory than was allocated, resulting in net memory reduction.

**8. Goroutine count explained**
> Sequential = 5 (background only), WaitGroup = 19 (5 + 12 yours), Channels = 20 (5 + 12 yours + 1 channel closer goroutine).

**9. CPU profile insight**
> Program is network bound — CPU was idle 61% of the time waiting for API responses. This explains why concurrency gives such dramatic speedup (98-99% faster).

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
| `runtime/pprof` | stdlib | CPU profiling |
| `runtime` | stdlib | Memory and goroutine profiling |

---

## Author

Built as an internship learning task to understand Go concurrency patterns — goroutines, channels, WaitGroup, performance benchmarking, and Go profiling tools.