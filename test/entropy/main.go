// test/entropytest/main.go
package main

import (
	"datflux/internal/entropy"
	"datflux/internal/password"
	"flag"
	"fmt"
	"math"
	"os"
	"time"
)

func main() {
	standardSamples := flag.Int("standard", 1000, "Number of standard mode passwords to test")
	paranoiaSamples := flag.Int("paranoia", 100, "Number of paranoia mode passwords to test")
	paranoiaCandidates := flag.Int("candidates", 25, "Candidates per paranoia password")
	outputFile := flag.String("output", "entropy_results.txt", "Output file for results")
	flag.Parse()

	file, err := os.Create(*outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	collector := entropy.NewCollector(time.Millisecond*100, 50)
	defer collector.Close()
	generator := password.NewGenerator(collector)

	logf := func(format string, args ...any) {
		fmt.Printf(format, args...)
		fmt.Fprintf(file, format, args...)
	}

	logf("datFlux Entropy Test Results\n")
	logf("========================\n")
	logf("Test run at: %s\n\n", time.Now().Format(time.RFC1123))

	// STANDARD MODE TEST
	logf("STANDARD MODE TESTING\n")
	logf("-------------------\n")
	results := testMode(generator, false, 0, *standardSamples, logf)
	printResults(results, "Standard Mode", logf)

	// PARANOIA MODE TEST
	logf("\nPARANOIA MODE TESTING\n")
	logf("-------------------\n")
	logf("Testing with %d candidates per password\n\n", *paranoiaCandidates)
	results = testMode(generator, true, *paranoiaCandidates, *paranoiaSamples, logf)
	printResults(results, "Paranoia Mode", logf)

	logf("\nTest completed successfully.\n")
	fmt.Println("Results saved to", *outputFile)
}

type TestResults struct {
	Passwords      []string
	EntropyValues  []float64
	LengthCounts   map[int]int
	TotalEntropy   float64
	MinEntropy     float64
	MaxEntropy     float64
	TotalPasswords int
}

func testMode(generator *password.Generator, paranoia bool, samples, count int, logf func(string, ...interface{})) TestResults {
	generator.SetParanoiaMode(paranoia, samples)

	results := TestResults{
		Passwords:      make([]string, 0, count),
		EntropyValues:  make([]float64, 0, count),
		LengthCounts:   make(map[int]int),
		MinEntropy:     math.MaxFloat64,
		TotalPasswords: count,
	}

	logf("Generating and analyzing %d passwords...\n", count)

	for i := 0; i < count; i++ {
		if i%max(count/10, 1) == 0 {
			logf("Progress: %d/%d passwords\n", i, count)
		}

		pwd := generator.Generate()
		strength := generator.AnalyzeStrength(pwd)
		entropyBits := strength.EntropyBits

		results.Passwords = append(results.Passwords, pwd)
		results.EntropyValues = append(results.EntropyValues, entropyBits)
		results.TotalEntropy += entropyBits
		results.LengthCounts[len(pwd)]++

		if entropyBits < results.MinEntropy {
			results.MinEntropy = entropyBits
		}
		if entropyBits > results.MaxEntropy {
			results.MaxEntropy = entropyBits
		}
	}

	return results
}

func printResults(results TestResults, label string, logf func(string, ...any)) {
	avgEntropy := results.TotalEntropy / float64(results.TotalPasswords)

	// standard deviation
	var sumSquaredDiff float64
	for _, entropy := range results.EntropyValues {
		diff := entropy - avgEntropy
		sumSquaredDiff += diff * diff
	}
	stdDev := math.Sqrt(sumSquaredDiff / float64(results.TotalPasswords))

	// results
	logf("\n%s Entropy Analysis (%d passwords):\n", label, results.TotalPasswords)
	logf("Average Entropy: %.2f bits\n", avgEntropy)
	logf("Minimum Entropy: %.2f bits\n", results.MinEntropy)
	logf("Maximum Entropy: %.2f bits\n", results.MaxEntropy)
	logf("Standard Deviation: %.2f bits\n", stdDev)

	// 10 bins
	logf("\nHistogram data:\n")
	bins := 10
	binSize := (results.MaxEntropy - results.MinEntropy) / float64(bins)
	histogram := make([]int, bins)

	for _, entropy := range results.EntropyValues {
		binIndex := int((entropy - results.MinEntropy) / binSize)
		if binIndex >= bins {
			binIndex = bins - 1
		}
		histogram[binIndex]++
	}

	for i := 0; i < bins; i++ {
		lowerBound := results.MinEntropy + float64(i)*binSize
		upperBound := results.MinEntropy + float64(i+1)*binSize
		count := histogram[i]
		percentage := float64(count) * 100 / float64(results.TotalPasswords)
		logf("%.1f-%.1f bits: %d passwords (%.1f%%)\n", lowerBound, upperBound, count, percentage)
	}

	// length stats
	logf("\nPassword Length Distribution:\n")
	var minLen, maxLen int = math.MaxInt, 0
	for length := range results.LengthCounts {
		if length < minLen {
			minLen = length
		}
		if length > maxLen {
			maxLen = length
		}
	}

	for length := minLen; length <= maxLen; length++ {
		count := results.LengthCounts[length]
		if count > 0 {
			percentage := float64(count) * 100 / float64(results.TotalPasswords)
			logf("Length %d: %d passwords (%.1f%%)\n", length, count, percentage)
		}
	}

	if results.TotalPasswords > 0 {
		sampleIndex := results.TotalPasswords / 2 // middle password
		samplePwd := results.Passwords[sampleIndex]
		sampleEntropy := results.EntropyValues[sampleIndex]
		logf("\nSample Password (median):\n")
		logf("Password: %s\n", samplePwd)
		logf("Length: %d characters\n", len(samplePwd))
		logf("Entropy: %.2f bits\n", sampleEntropy)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
