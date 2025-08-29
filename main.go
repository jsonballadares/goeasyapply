package main

import (
	"github.com/jsonballadares/goeasyapply/scraper"
	"log"
)

func main() {
	jobIDs, err := scraper.ScrapeJobIDs()
	if err != nil {
		log.Fatal(err)
	}

	for _, id := range jobIDs {
		log.Println("Job ID:", id)
	}
}
