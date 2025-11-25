package promotion

import (
	"database/sql"
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
	Count    int
	Err      error
}

func (s *Search) IsValid(pattern string, files []string) bool {
	// Check the cache.
	cachedResult, err := s.repo.GetCodeFileMatchCount(pattern)
	if err != nil {
		if err == sql.ErrNoRows {
			// Cache miss - continue to search files.
		} else {
			// TODO: use slog to make this debug level logging.
			log.Printf("cache lookup failed, performing fresh search: %v", err)
		}
	} else {
		// Cache hit.
		return cachedResult >= minFileCount
	}

	resultsChan := make(chan FileResult, len(files))
	var fileWg sync.WaitGroup

	// Process each file concurrently
	for _, filepath := range files {
		fileWg.Add(1)
		// Launch a goroutine for each file search
		go func(fp string) {
			defer fileWg.Done()
			count, err := SearchFileParallel(fp, pattern)
			resultsChan <- FileResult{FilePath: fp, Count: count, Err: err}
		}(filepath)
	}

	fileWg.Wait()
	close(resultsChan)

	filesWithMatchCount := 0

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

		if res.Count > 0 {
			filesWithMatchCount++
		}
	}

	// Populate the cache.
	if err := s.repo.AddCodeFileMatchCount(pattern, filesWithMatchCount); err != nil {
		// only log the issue, the result has already been calculated.
		log.Printf("Caching result failed with error: %v", err)
	}

	return filesWithMatchCount >= minFileCount
}
