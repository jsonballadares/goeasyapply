package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/jsonballadares/goeasyapply/scraper"
	"github.com/jsonballadares/goeasyapply/util"
)

func main() {
	// TODO: pull these from a vault or config
	username := "jsonballadares@gmail.com"
	password := "<yL21n8!O&b$&X5?"
	// --- Chrome options (headful) ---
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // show browser
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("start-maximized", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer allocCancel()

	// --- Create browser context ---
	ctx, ctxCancel := chromedp.NewContext(allocCtx)
	defer ctxCancel()

	// Optional: Set a timeout for the entire scraping session
	ctx, timeoutCancel := context.WithTimeout(ctx, 15*time.Second)
	defer timeoutCancel()

	fmt.Println("Launching Chrome...")

	// Start the browser
	if err := chromedp.Run(ctx); err != nil {
		log.Fatalf("Failed to start Chrome: %v", err)
	}

	// Login
	if err := util.LinkedInLogin(ctx, username, password); err != nil {
		log.Fatalf("Failed to login to LinkedIn: %v", err)
	}

	// --- Start scraping ---
	fmt.Println("Starting scraping job IDs...")

	// Pass the context into the scraper
	jobIDs, err := scraper.ScrapeJobIDs(ctx, "https://www.linkedin.com/jobs/search/?currentJobId=4291616774&f_AL=true&f_WT=2&geoId=103644278&keywords=golang", 25)
	if err != nil {
		log.Fatalf("Scraping failed: %v", err)
	}

	fmt.Printf("Scraping finished. Found %d new job IDs:\n", len(jobIDs))
	for _, id := range jobIDs {
		fmt.Println(id)
	}
}
