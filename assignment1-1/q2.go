package cos418_hw1_1

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"sync"
)

// Sum numbers from channel `nums` and output sum to `out`.
// You should only output to `out` once.
func sumWorker(nums chan int, out chan int) {
	sum := 0
	for num := range nums {
		sum += num
	}
	out <- sum // Send the sum once when done
}

// Read integers from the file `fileName` and return sum of all values.
// This function must launch `num` goroutines running `sumWorker` to find the sum concurrently.
func sum(num int, fileName string) int {
	// Open the file
	file, err := os.Open(fileName)
	checkError(err)
	defer file.Close()

	// Read all integers from the file
	ints, err := readInts(file)
	checkError(err)

	// Create channels
	nums := make(chan int, len(ints)) // Buffered channel to distribute numbers
	out := make(chan int, num)        // Channel for worker outputs

	// Launch `num` goroutines
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sumWorker(nums, out)
		}()
	}

	// Send numbers to the `nums` channel
	for _, n := range ints {
		nums <- n
	}
	close(nums) // Signal that no more numbers will be sent

	// Wait for all goroutines to finish and close the output channel
	go func() {
		wg.Wait()
		close(out)
	}()

	// Aggregate the results from all workers
	totalSum := 0
	for partialSum := range out {
		totalSum += partialSum
	}

	return totalSum
}



// Read a list of integers separated by whitespace from `r`.
// Return the integers successfully read with no error, or
// an empty slice of integers and the error that occurred.
// Do NOT modify this function.
func readInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var elems []int
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return elems, err
		}
		elems = append(elems, val)
	}
	return elems, nil
}
