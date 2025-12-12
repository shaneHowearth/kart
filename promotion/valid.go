package promotion

import (
	"fmt"
	"log"
	"sync"

	"github.com/shanehowearth/kart/internal/validation"
)

const minFileCount = 2

type Search struct {
	repo Store
}

func NewSearch(repo Store) (*Search, error) {
	if validation.IsNil(repo) {
		// TODO sentinel error
		return nil, fmt.Errorf("supplied store is nil")
	}

	return &Search{repo: repo}, nil
}

// Result struct used by main to hold results from parallel file processing
type FileResult struct {
	FilePath string
	Counts   map[string]int
	Err      error
}

func (s *Search) IsValidBatch(patterns []string, files []string) map[string]bool {
	missedPatterns := []string{}
	results := map[string]bool{}

	// Check the cache.
	// TODO: What to do if there are partial misses.
	cachedResults, err := s.repo.GetCodeFileMatchCounts(patterns)
	if err != nil {
		log.Println("Cache error, please fix %w", err)
	}

	for pattern, result := range cachedResults {
		if result.Found {
			// Cache hit - use cached value
			results[pattern] = result.MatchCount >= minFileCount
		} else {
			// Cache miss - need to search
			missedPatterns = append(missedPatterns, pattern)
		}
	}

	if len(missedPatterns) == 0 {
		// nothing left to do.
		return results
	}

	resultsChan := make(chan FileResult, len(files))
	var fileWg sync.WaitGroup

	// Process each file concurrently
	for _, filepath := range files {
		fileWg.Add(1)
		// Launch a goroutine for each file search
		go func(fp string) {
			defer fileWg.Done()
			counts, err := SearchFileParallel(fp, missedPatterns)
			resultsChan <- FileResult{FilePath: fp, Counts: counts, Err: err}
		}(filepath)
	}

	fileWg.Wait()
	close(resultsChan)

	tmpResults := map[string]int{}
	for res := range resultsChan {
		if res.Err != nil {
			// TODO: Need a requirement here, should the search keep on
			// trucking, which could result in false positives/negatives
			// (depending on the way things are calculated), or should
			// everything die in a fire.
			// Am just producing a log so that the user can make a decision on
			// what to do - fix, or rerun (with the cache needing to be
			// rectified).
			log.Printf("Error processing file %s: %v", res.FilePath, res.Err)
			continue
		}

		for code, count := range res.Counts {
			if count > 0 {
				tmpResults[code]++
			}
		}
	}

	// Populate the cache.
	if err := s.repo.AddCodeFileMatchCounts(tmpResults); err != nil {
		// only log the issue, the result has already been calculated.
		log.Printf("Caching result failed with error: %v", err)
	}

	// Add results for patterns that were searched.
	for k, v := range tmpResults {
		results[k] = v >= minFileCount
	}

	// Add false for patterns that had zero matches.
	for _, pattern := range missedPatterns {
		if _, exists := results[pattern]; !exists {
			results[pattern] = false
		}
	}

	return results
}
