package utils

import (
	"fmt"
	"time"
)

func ShowResults(
	seqDuration time.Duration, seqSuccess int,
	wgDuration time.Duration, wgSuccess int,
	chDuration time.Duration, chSuccess int,
) {
	totalAPIs := 12

	seqMs := seqDuration.Milliseconds()
	wgMs := wgDuration.Milliseconds()
	chMs := chDuration.Milliseconds()

	// Header
	fmt.Println("\n================= API PERFORMANCE TEST =================")
	fmt.Printf("\nTotal APIs Called: %d\n", totalAPIs)

	// Sequential
	fmt.Println("\n[1] Sequential Execution:")
	fmt.Println("-----------------------------------")
	fmt.Printf("Time Taken : %d ms\n", seqMs)
	fmt.Printf("Success    : %d/%d\n", seqSuccess, totalAPIs)

	// WaitGroup
	fmt.Println("\n[2] Concurrent (WaitGroup):")
	fmt.Println("-----------------------------------")
	fmt.Printf("Time Taken : %d ms\n", wgMs)
	fmt.Printf("Success    : %d/%d\n", wgSuccess, totalAPIs)

	// Channels
	fmt.Println("\n[3] Concurrent (Channels):")
	fmt.Println("-----------------------------------")
	fmt.Printf("Time Taken : %d ms\n", chMs)
	fmt.Printf("Success    : %d/%d\n", chSuccess, totalAPIs)

	fmt.Println("\n================= COMPARISON =================")

	wgVsSeq := float64(seqMs-wgMs) / float64(seqMs) * 100
	chVsSeq := float64(seqMs-chMs) / float64(seqMs) * 100
	wgVsCh := float64(chMs-wgMs) / float64(chMs) * 100

	fmt.Println("\nPerformance Gain:")

	if wgMs < seqMs {
		fmt.Printf("- WaitGroup vs Sequential  : ~%.0f%% faster\n", wgVsSeq)
	} else {
		fmt.Printf("- WaitGroup vs Sequential  : ~%.0f%% slower\n", -wgVsSeq)
	}

	if chMs < seqMs {
		fmt.Printf("- Channels vs Sequential   : ~%.0f%% faster\n", chVsSeq)
	} else {
		fmt.Printf("- Channels vs Sequential   : ~%.0f%% slower\n", -chVsSeq)
	}

	if wgMs < chMs {
		fmt.Printf("- WaitGroup vs Channels    : ~%.0f%% faster than Channels\n", -wgVsCh)
	} else {
		fmt.Printf("- WaitGroup vs Channels    : ~%.0f%% slower than Channels\n", wgVsCh)
	}

}
