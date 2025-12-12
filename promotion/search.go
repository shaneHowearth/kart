package promotion

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/sys/unix"
)

// SearchFileParallel searches files using concurrency.
// File is mmapped for faster access (kernel manages access), and then broken up
// into chunks that are then passed to goroutines to be searched.
func SearchFileParallel(filepath string, patterns []string) (map[string]int, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}
	size := int(fi.Size())

	if size == 0 {
		return map[string]int{}, nil
	}

	// Mmap the entire file.
	// TODO: Add thoughts on advantages/disadvantages of using mmap.
	fullData, err := unix.Mmap(int(f.Fd()), 0, size, unix.PROT_READ, unix.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("failed to mmap file: %w", err)
	}
	defer unix.Munmap(fullData)

	// Note: This assumes that the file only contains uppercase codes.
	// Need a requirement check that is correct.
	upperPatterns := make([][]byte, 0, len(patterns))
	for _, pattern := range patterns {
		upperPatterns = append(upperPatterns, []byte(strings.ToUpper(pattern)))
	}

	// Determine parallelism
	numCores := runtime.GOMAXPROCS(0)
	if numCores == 0 {
		numCores = 4
	}

	chunkSize := size / numCores

	var wg sync.WaitGroup

	countChan := make(chan map[string]int)

	var sumWg sync.WaitGroup
	sumWg.Add(1)

	totalMatches := map[string]int{}
	go func() {
		defer sumWg.Done()
		for counts := range countChan {
			for pattern, count := range counts {
				totalMatches[pattern] += count
			}
		}
	}()

	// Divide file into chunks and launch goroutines
	currentStart := 0
	for i := 0; i < numCores; i++ {
		currentEnd := currentStart + chunkSize
		if i == numCores-1 {
			currentEnd = size
		}

		// Adjust the end boundary to land on a newline for line-accurate counting
		if currentEnd < size {
			for currentEnd > currentStart && fullData[currentEnd] != '\n' {
				currentEnd--
			}
			currentEnd++
		}

		if currentEnd > currentStart {
			wg.Add(1)
			chunk := fullData[currentStart:currentEnd]
			go SearchChunks(chunk, upperPatterns, &wg, countChan)

			currentStart = currentEnd
		}

		if currentStart >= size {
			break
		}
	}

	wg.Wait()
	close(countChan)
	sumWg.Wait() // Wait for the aggregation goroutine to finish

	return totalMatches, nil
}

// SearchChunk searches the provided file chunk for the string.
// The string matching is by bytes, converting to runes will cause allocations
// and slow things down, and would only be useful is we were looking for the nth
// character.
func SearchChunks(data []byte, patterns [][]byte, wg *sync.WaitGroup, countChan chan<- map[string]int) {
	defer wg.Done()

	currentOffset := 0
	matches := map[string]int{}

	for currentOffset < len(data) {
		// Find the next newline character (end of the current line within the chunk)
		lineEnd := bytes.IndexByte(data[currentOffset:], '\n')

		var line []byte
		if lineEnd == -1 {
			// End of file.
			line = data[currentOffset:]
			currentOffset = len(data) // Exit loop after processing the last line
		} else {
			// Extract line, without newline.
			line = data[currentOffset : currentOffset+lineEnd]
			currentOffset = currentOffset + lineEnd + 1 // Move offset past the newline
		}

		for i, pattern := range patterns {
			// Case insensitive search - pattern is already UCase, and files are
			// assumed to be UCase codes only.
			if bytes.Equal(line, pattern) {
				matches[string(patterns[i])]++
			}
		}
	}

	// Send the total non-overlapping occurrences found in this chunk
	countChan <- matches
}
