package utils

import (
	"fmt"
	"runtime"
	"time"
)

// ProfilingData holds profiling information for each execution phase
type ProfilingData struct {
	MemBefore  runtime.MemStats
	MemAfter   runtime.MemStats
	Goroutines int
	TimeTaken  time.Duration
}

// CaptureMemStats captures current memory statistics
func CaptureMemStats() runtime.MemStats {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return mem
}

// CaptureGoroutines returns current number of goroutines
func CaptureGoroutines() int {
	return runtime.NumGoroutine()
}

// formatMemory handles both positive and negative memory changes
func formatMemory(bytes float64) string {
	mb := bytes / 1024 / 1024
	if mb < 0 {
		return fmt.Sprintf("GC freed %.2f MB", -mb)
	}
	return fmt.Sprintf("+%.2f MB", mb)
}

func PrintProfilingReport(
	seqData ProfilingData,
	wgData ProfilingData,
	chData ProfilingData,
	memBefore runtime.MemStats,
	memAfter runtime.MemStats,
) {

	fmt.Println("\n================= PROFILING REPORT =================")

	// Memory Stats — Before and After entire execution
	fmt.Println("\n------------------ MEMORY STATS ------------------")
	fmt.Println("Before Execution:")
	fmt.Printf("  Alloc = %.2f MB | Sys = %.2f MB | NumGC = %d\n",
		float64(memBefore.Alloc)/1024/1024,
		float64(memBefore.Sys)/1024/1024,
		memBefore.NumGC,
	)
	fmt.Println("After Execution:")
	fmt.Printf("  Alloc = %.2f MB | Sys = %.2f MB | NumGC = %d\n",
		float64(memAfter.Alloc)/1024/1024,
		float64(memAfter.Sys)/1024/1024,
		memAfter.NumGC,
	)

	fmt.Println("\n------------------ PHASE PROFILING ------------------")

	// Sequential
	seqMemUsed := float64(seqData.MemAfter.Alloc) - float64(seqData.MemBefore.Alloc)
	fmt.Println("\n[1] Sequential Execution")
	fmt.Printf("  Time Taken     : %d ms\n", seqData.TimeTaken.Milliseconds())
	fmt.Printf("  Memory Used    : %s\n", formatMemory(seqMemUsed))
	fmt.Printf("  Goroutines     : %d\n", seqData.Goroutines)

	// WaitGroup
	wgMemUsed := float64(wgData.MemAfter.Alloc) - float64(wgData.MemBefore.Alloc)
	fmt.Println("\n[2] WaitGroup Execution")
	fmt.Printf("  Time Taken     : %d ms\n", wgData.TimeTaken.Milliseconds())
	fmt.Printf("  Memory Used    : %s\n", formatMemory(wgMemUsed))
	fmt.Printf("  Goroutines     : %d\n", wgData.Goroutines)

	// Channel
	chMemUsed := float64(chData.MemAfter.Alloc) - float64(chData.MemBefore.Alloc)
	fmt.Println("\n[3] Channel Execution")
	fmt.Printf("  Time Taken     : %d ms\n", chData.TimeTaken.Milliseconds())
	fmt.Printf("  Memory Used    : %s\n", formatMemory(chMemUsed))
	fmt.Printf("  Goroutines     : %d\n", chData.Goroutines)

	fmt.Println("\n------------------ CPU PROFILE ------------------")
	fmt.Println("  CPU Profiling Started...")
	fmt.Println("  CPU Profiling Stopped.")
	fmt.Println("  CPU Profile saved to: cpu.prof")
	fmt.Println("  Run: go tool pprof cpu.prof")

	fmt.Println("\n------------------ SUMMARY ------------------")

	// Find fastest method
	fastest := "Sequential"
	fastestTime := seqData.TimeTaken
	if wgData.TimeTaken < fastestTime {
		fastest = "WaitGroup"
		fastestTime = wgData.TimeTaken
	}
	if chData.TimeTaken < fastestTime {
		fastest = "Channels"
	}

	// Find highest memory usage
	highestMem := "Sequential"
	highestMemVal := seqMemUsed
	if wgMemUsed > highestMemVal {
		highestMem = "WaitGroup"
		highestMemVal = wgMemUsed
	}
	if chMemUsed > highestMemVal {
		highestMem = "Channels"
	}

	fmt.Printf("  Fastest Method        : %s\n", fastest)
	fmt.Printf("  Highest Memory Usage  : %s\n", highestMem)
	fmt.Printf("  Most Efficient Method : %s\n", fastest)
	fmt.Printf("\n=================================================\n")
}