package main

import (
	"fmt"
	"sync"
	"time"
)

// Структура задачі
type Job struct {
	ID    int
	Input int
}

// Структура результату
type Result struct {
	JobID  int
	Output int
	Worker int
}

// Функція воркера: читає задачі з каналу jobs і записує результат у results
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		// Імітація обчислювальної або мережевої роботи
		time.Sleep(50 * time.Millisecond)
		results <- Result{
			JobID:  job.ID,
			Output: job.Input * job.Input,
			Worker: id,
		}
	}
}

func main() {
	const totalJobs = 20
	const numWorkers = 4

	jobs := make(chan Job, totalJobs)
	results := make(chan Result, totalJobs)
	var wg sync.WaitGroup

	// 1. Запуск пулу воркерів
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// 2. Відправка задач у буферизований канал
	for j := 1; j <= totalJobs; j++ {
		jobs <- Job{ID: j, Input: j}
	}
	close(jobs) // Закриття вхідного каналу сигналізує воркерам про кінець задач

	// 3. Окремий потік для безпечного закриття результатів після завершення всіх воркерів
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. Збір та агрегація результатів
	fmt.Println("Обробка задач пулом воркерів")
	for res := range results {
		fmt.Printf("Задача №%2d оброблена Воркером #%d. Результат: %d\n",
			res.JobID, res.Worker, res.Output)
	}
}
