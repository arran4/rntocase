package rntocase

import (
	"errors"
	"testing"
)

// Simple algo functions for benchmarking
func algoSuccess(s string) (string, error) {
	return s, nil
}

func algoError(s string) (string, error) {
	return "", errors.New("error")
}

func algoPanic(s string) (string, error) {
	panic("panic")
}

var benchmarkAlgos = map[string]func(string) (string, error){
	"success": algoSuccess,
	"error":   algoError,
	"panic":   algoPanic,
}

// Benchmark the original Run function
func BenchmarkRun_Success(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run(benchmarkAlgos, "success", "test")
	}
}

func BenchmarkRun_Panic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run(benchmarkAlgos, "panic", "test")
	}
}

// Benchmark the batch approach used in optimized RenderUsageTable
func BenchmarkRunSafeBatch_Success_4Items(b *testing.B) {
	values := []string{"1", "2", "3", "4"}
	dummyYield := func(s string) bool { return true }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runSafeBatch(algoSuccess, values, dummyYield)
	}
}

// Benchmark the loop approach (simulating the old implementation)
func BenchmarkRunLoop_Success_4Items(b *testing.B) {
	values := []string{"1", "2", "3", "4"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			Run(benchmarkAlgos, "success", v)
		}
	}
}

func BenchmarkRunSafeBatch_Panic_4Items(b *testing.B) {
	// Panic on 2nd item
	algoMixed := func(s string) (string, error) {
		if s == "2" {
			panic("panic")
		}
		return s, nil
	}
	values := []string{"1", "2", "3", "4"}
	dummyYield := func(s string) bool { return true }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runSafeBatch(algoMixed, values, dummyYield)
	}
}

func BenchmarkRunLoop_Panic_4Items(b *testing.B) {
	algoMixed := map[string]func(string)(string, error){
		"mixed": func(s string) (string, error) {
			if s == "2" {
				panic("panic")
			}
			return s, nil
		},
	}
	values := []string{"1", "2", "3", "4"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			Run(algoMixed, "mixed", v)
		}
	}
}
