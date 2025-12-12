package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shanehowearth/kart/promotion"
	"github.com/shanehowearth/kart/promotion/datastore/sqlite"
)

// Prepare some storage for the patterns passed in by users.
type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	// Initialise dependencies.
	promotionStore := &sqlite.Driver{}
	if err := promotionStore.InitialiseDataStore(); err != nil {
		log.Fatalf("cannot initialise promotion datastore with error: %v", err)
	}

	promotionSearch, err := promotion.NewSearch(promotionStore)
	if err != nil {
		log.Fatalf("cannot create a promotion search with error %v", err)
	}

	// Parse flags
	var patterns stringSlice
	flag.Var(&patterns, "p", "promotion code to search (can be specified multiple times)")
	flag.Parse()

	// Remaining args are files
	files := flag.Args()

	if len(patterns) == 0 || len(files) == 0 {
		// Require at least one pattern and at least one file to be passed in.
		log.Fatalf("Usage: %s -p <pattern> [-p <pattern2>...] <file1> [file2] [file3]...\n", os.Args[0])
	}

	// Crude check that the files aren't gzipped.
	// A deeper inspection is possible once the files have been opened, and the
	// magic number inspected.
	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			log.Fatalf("File does not exist: %s", file)
		}
		if strings.HasSuffix(file, ".gz") {
			log.Fatalf("ERROR: Found gzipped file %s - all files must be uncompressed.",
				file)
		}
	}

	results := promotionSearch.IsValidBatch(patterns, files)

	for pattern, isValid := range results {
		validity := "an invalid"
		if isValid {
			validity = "a valid"
		}

		fmt.Printf("%s is %s coupon\n", pattern, validity)
	}
}
