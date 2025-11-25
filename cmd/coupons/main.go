package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shanehowearth/kart/promotion"
	"github.com/shanehowearth/kart/promotion/datastore/sqlite"
)

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

	// Require a pattern and at least one file to be passed in.
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <pattern> <file1> [file2] [file3]...\n", os.Args[0])
	}

	pattern := os.Args[1]
	files := os.Args[2:]

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

	validity := "an invalid"
	if promotionSearch.IsValid(pattern, files) {
		validity = "a valid"
	}

	fmt.Printf("%s is %s coupon\n", pattern, validity)
}
