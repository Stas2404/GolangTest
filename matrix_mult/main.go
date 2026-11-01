package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// Ініціалізація квадратної матриці випадковими числами
func createRandomMatrix(n int) [][]float64 {
	matrix := make([][]float64, n)
	for i := range matrix {
		matrix[i] = make([]float64, n)
		for j := range matrix[i] {
			matrix[i][j] = rand.Float64()
		}
	}
	return matrix
}

// 1. Послідовне множення матриць
func multiplySequential(a, b [][]float64, n int) [][]float64 {
	c := make([][]float64, n)
	for i := range c {
		c[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			sum := 0.0
			for k := 0; k < n; k++ {
				sum += a[i][k] * b[k][j]
			}
			c[i][j] = sum
		}
	}
	return c
}

// 2. Паралельне множення матриць за допомогою горутин
func multiplyParallel(a, b [][]float64, n int, numWorkers int) [][]float64 {
	c := make([][]float64, n)
	for i := range c {
		c[i] = make([]float64, n)
	}

	var wg sync.WaitGroup
	rowsPerWorker := (n + numWorkers - 1) / numWorkers

	for w := 0; w < numWorkers; w++ {
		startRow := w * rowsPerWorker
		endRow := startRow + rowsPerWorker
		if startRow >= n {
			break
		}
		if endRow > n {
			endRow = n
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				for j := 0; j < n; j++ {
					sum := 0.0
					for k := 0; k < n; k++ {
						sum += a[i][k] * b[k][j]
					}
					c[i][j] = sum
				}
			}
		}(startRow, endRow)
	}

	wg.Wait() // Очікування завершення роботи всіх горутин
	return c
}

func main() {
	const matrixSize = 1000 // Розмір матриці 1000x1000
	numCores := runtime.NumCPU()
	runtime.GOMAXPROCS(numCores)

	fmt.Printf("Тест множення матриць (%dx%d)\n", matrixSize, matrixSize)
	fmt.Printf("Доступна кількість логічних ядер CPU: %d\n\n", numCores)

	a := createRandomMatrix(matrixSize)
	b := createRandomMatrix(matrixSize)

	// Замір послідовної версії
	startSeq := time.Now()
	_ = multiplySequential(a, b, matrixSize)
	durationSeq := time.Since(startSeq)
	fmt.Printf("Послідовне виконання: %v\n", durationSeq)

	// Замір паралельної версії
	startPar := time.Now()
	_ = multiplyParallel(a, b, matrixSize, numCores)
	durationPar := time.Since(startPar)
	fmt.Printf("Паралельне виконання (%d потоків): %v\n", numCores, durationPar)

	speedup := float64(durationSeq)/float64(durationPar)
	fmt.Printf("Коефіцієнт прискорення (Speedup): %.2fx\n", speedup)
}