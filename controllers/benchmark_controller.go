package controllers

import (
	"apiracer/service"
	"apiracer/utils"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

type BenchmarkController struct {
	beego.Controller
}

func (c *BenchmarkController) RunBenchmark() {

	fmt.Println("\n[1] Running Sequential Execution...")
	fmt.Println("-----------------------------------")
	seqDuration, seqSuccess := service.RunSequential()

	fmt.Println("\n[2] Running Concurrent (WaitGroup)...")
	fmt.Println("-----------------------------------")
	wgDuration, wgSuccess := service.RunWaitGroup()

	fmt.Println("\n[3] Running Concurrent (Channels)...")
	fmt.Println("-----------------------------------")
	chDuration, chSuccess := service.RunChannel()

	utils.ShowResults(
		seqDuration, seqSuccess,
		wgDuration, wgSuccess,
		chDuration, chSuccess,
	)

	c.Data["json"] = map[string]interface{}{
		"status": "completed",
		"results": map[string]interface{}{
			"sequential": map[string]interface{}{
				"time_ms": seqDuration.Milliseconds(),
				"success": seqSuccess,
			},
			"waitgroup": map[string]interface{}{
				"time_ms": wgDuration.Milliseconds(),
				"success": wgSuccess,
			},
			"channels": map[string]interface{}{
				"time_ms": chDuration.Milliseconds(),
				"success": chSuccess,
			},
		},
	}
	c.ServeJSON()
}