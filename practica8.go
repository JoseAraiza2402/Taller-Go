package main

import (
	"fmt"
	"sync"
)

type Result struct {
	Number    int
	Factorial uint64
	Error     error
}

func calculateFactorial(n int, ch chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	if n < 0 {
		ch <- Result{Number: n, Error: fmt.Errorf("factorial no definido para números negativos: %d", n)}
		return
	}

	if n > 20 {
		ch <- Result{Number: n, Error: fmt.Errorf("número demasiado grande para calcular factorial: %d", n)}
		return
	}

	result := uint64(1)
	for i := 1; i <= n; i++ {
		result *= uint64(i)
	}

	ch <- Result{Number: n, Factorial: result}
}

func main() {
	numbers := []int{5, 7, 3, 10, 4, -1, 21}

	resultChan := make(chan Result, len(numbers))

	var wg sync.WaitGroup

	for _, num := range numbers {
		wg.Add(1)
		go calculateFactorial(num, resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	results := make(map[int]Result)
	for result := range resultChan {
		results[result.Number] = result
	}

	for _, num := range numbers {
		result := results[num]
		if result.Error != nil {
			fmt.Printf("Error para %d: %v\n", num, result.Error)
		} else {
			fmt.Printf("Factorial de %d = %d\n", num, result.Factorial)
		}
	}
}
