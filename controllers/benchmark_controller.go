package controllers

import (
	"apiracer/service"
	"apiracer/utils"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type BenchmarkController struct {
	beego.Controller
}

func (c *BenchmarkController) RunBenchmark() {

	// Record total request start time
	totalStart := time.Now()

	// Capture memory BEFORE everything starts
	memBefore := utils.CaptureMemStats()

	// CPU Profiling — write to file instead of terminal
	fmt.Println("\n------------------ CPU PROFILE ------------------")
	f, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("  Could not create CPU profile file:", err)
	} else {
		pprof.StartCPUProfile(f)
		fmt.Println("  CPU Profiling Started...")
		defer func() {
			pprof.StopCPUProfile()
			f.Close()
			fmt.Println("  CPU Profiling Stopped.")
			fmt.Println("  CPU Profile saved to: cpu.prof")
			fmt.Println("  Run: go tool pprof cpu.prof")
		}()
	}

	// ============ Phase 1: Sequential ============
	fmt.Println("\n[1] Running Sequential Execution...")
	fmt.Println("-----------------------------------")
	seqDuration, seqSuccess, seqData := service.RunSequential()

	// ============ Phase 2: WaitGroup ============
	fmt.Println("\n[2] Running Concurrent (WaitGroup)...")
	fmt.Println("-----------------------------------")
	wgDuration, wgSuccess, wgData := service.RunWaitGroup()

	// ============ Phase 3: Channels ============
	fmt.Println("\n[3] Running Concurrent (Channels)...")
	fmt.Println("-----------------------------------")
	chDuration, chSuccess, chData := service.RunChannel()

	// Capture memory AFTER everything finishes
	memAfter := utils.CaptureMemStats()

	// Total execution time
	totalDuration := time.Since(totalStart)
	fmt.Printf("\n  Total Execution Time: %d ms\n", totalDuration.Milliseconds())

	// Show performance comparison
	utils.ShowResults(
		seqDuration, seqSuccess,
		wgDuration, wgSuccess,
		chDuration, chSuccess,
	)

	// Show profiling report
	utils.PrintProfilingReport(
		seqData,
		wgData,
		chData,
		memBefore,
		memAfter,
	)

	// Return JSON response
	c.Data["json"] = map[string]interface{}{
		"status":        "completed",
		"total_time_ms": totalDuration.Milliseconds(),
		"results": map[string]interface{}{
			"sequential": map[string]interface{}{
				"time_ms":    seqDuration.Milliseconds(),
				"success":    seqSuccess,
				"goroutines": seqData.Goroutines,
			},
			"waitgroup": map[string]interface{}{
				"time_ms":    wgDuration.Milliseconds(),
				"success":    wgSuccess,
				"goroutines": wgData.Goroutines,
			},
			"channels": map[string]interface{}{
				"time_ms":    chDuration.Milliseconds(),
				"success":    chSuccess,
				"goroutines": chData.Goroutines,
			},
		},
	}
	c.ServeJSON()
}