package scraper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"math/rand"
)

// ScrapeJobIDs scrapes job IDs across multiple pages and skips jobs already applied
func ScrapeJobIDs(ctx context.Context, searchURL string, maxJobs int) ([]string, error) {
	var jobIDs []string

	// Calculate total pages to scrape
	perPage := 25
	totalPages := (maxJobs / perPage) + 1
	fmt.Printf("Total pages to scrape: %d\n", totalPages)

	for page := 0; page < totalPages; page++ {
		start := page * perPage
		pageURL := fmt.Sprintf("%s&start=%d", searchURL, start)
		fmt.Printf("Visiting page %d/%d: %s\n", page+1, totalPages, pageURL)

		// Navigate to the page
		if err := chromedp.Run(ctx,
			chromedp.Navigate(pageURL),
			chromedp.WaitVisible(`ul.jobs-search-results__list`, chromedp.ByQuery),
		); err != nil {
			log.Printf("Failed to navigate or wait for jobs list: %v", err)
			continue
		}

		// Scroll gradually to load all jobs
		if err := scrollJobList(ctx); err != nil {
			log.Printf("Scrolling failed: %v", err)
		}

		// Extract job IDs and applied messages
		var jobs []map[string]string
		err := chromedp.Run(ctx,
			chromedp.Evaluate(`
				Array.from(document.querySelectorAll("li[data-occludable-job-id]")).map(el => {
					let appliedEl = el.querySelector("span.artdeco-inline-feedback__message");
					let appliedText = appliedEl ? appliedEl.innerText.trim() : "";
					return { id: el.getAttribute("data-occludable-job-id"), applied: appliedText };
				})
			`, &jobs),
		)
		if err != nil {
			log.Printf("Evaluation error: %v", err)
			continue
		}

		// Filter jobs
		for _, job := range jobs {
			if job["applied"] != "" && job["applied"] != " " {
				fmt.Printf("Skipping job because of applied message: %s\n", job["applied"])
				continue
			}
			if job["id"] != "" {
				// Ensure uniqueness
				exists := false
				for _, j := range jobIDs {
					if j == job["id"] {
						exists = true
						break
					}
				}
				if !exists {
					jobIDs = append(jobIDs, job["id"])
				}
			}
		}

		fmt.Printf("Collected job IDs so far (%d): %v\n", len(jobIDs), jobIDs)

		// Random sleep to emulate human behavior
		time.Sleep(time.Duration(rand.Intn(4000)+3000) * time.Millisecond)
	}

	fmt.Printf("Final collected job IDs (%d): %v\n", len(jobIDs), jobIDs)
	return jobIDs, nil
}

// scrollJobList scrolls the jobs list to simulate human scrolling
func scrollJobList(ctx context.Context) error {
	for i := 0; i < 8; i++ {
		err := chromedp.Run(ctx,
			chromedp.Evaluate(`
				let ul = document.querySelector('ul.jobs-search-results__list');
				if (ul) { ul.scrollTop += ul.offsetHeight / 2; }
			`, nil),
		)
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(rand.Intn(700)+500) * time.Millisecond)
	}
	return nil
}
