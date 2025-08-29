package util

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

// LinkedInLogin logs in to LinkedIn using the provided context
func LinkedInLogin(ctx context.Context, username, password string) error {
	// Run login steps within the provided context
	return chromedp.Run(ctx,
		chromedp.Navigate("https://www.linkedin.com/login"),
		chromedp.WaitVisible(`#username`),
		chromedp.SendKeys(`#username`, username),
		chromedp.SendKeys(`#password`, password),
		chromedp.Click(`button[type=submit]`),
		chromedp.Sleep(3*time.Second), // wait for redirect
	)
}
