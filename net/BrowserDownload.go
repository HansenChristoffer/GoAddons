// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package net

import (
	"context"
	"database/sql"
	"fmt"
	"goaddons/database"
	"goaddons/models"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

const (
	dockerUrl = "ws://localhost:9222"
)

var userAgents = []string{
	// Desktop Browsers
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:86.0) Gecko/20100101 Firefox/86.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.705.50 Safari/537.36 Edg/88.0.705.50",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_16_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.2 Safari/605.1.15",

	// Other Devices
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.2 Safari/537.36 Edge/12.10136",
}

func StartHeadlessAndDownloadAddons(runId string, arr []models.Addon, dp string, db *sql.DB) (done bool, err error) {
	if arr == nil || dp == "" {
		return false, fmt.Errorf("addons or downloadPath is required\n")
	}

	ctx, cancel, err := setupContext(dp)
	if err != nil {
		return false, err
	}
	defer cancel()

	// Channel to signal the completion of each download
	downloadComplete := make(chan string)
	setupContextListeners(ctx, dp, downloadComplete)

	for i, a := range arr {
		log.Printf("[%d/%d] Will try to navigate to: %s\n", i+1, len(arr), a.DownloadUrl)

		if a.DownloadUrl == "" {
			log.Printf("DownloadURL is not allowed to be empty, will ignore this addon!\n")
		} else if strings.HasSuffix(a.DownloadUrl, ".zip") {
			_, err := handleDirectDownload(a, dp)
			if err != nil {
				log.Printf("Error with direct download! -> %v\n", err)
			}

			continue
		}

		navigateToAddonUrl(ctx, a)

		// Wait for the signal that the current download is complete
		<-downloadComplete

		err = handleDownloadLogging(db, a, runId)
		if err != nil {
			log.Println(err)
		}

		// Update the column 'last_downloaded' in kaasufouji.addons table
		_, err := database.UpdateAddon(db, a)
		if err != nil {
			return false, err
		}

		// Trying to accommodate for any type of latency issues
		time.Sleep(2 * time.Second)
		log.Printf("Continuing to the next addon...\n")
	}

	log.Println("Download of all addons is now done!")
	return true, nil
}

func setupContext(dp string) (ctx context.Context, cancelFunc context.CancelFunc, err error) {
	allocatorCtx, allocatorCancel := chromedp.NewRemoteAllocator(context.Background(), dockerUrl)

	// initialize a controllable Chrome instance
	ctx, cancel := chromedp.NewContext(allocatorCtx)

	if err = setupContextOptions(ctx, dp); err != nil {
		cancel()
		allocatorCancel()
		return nil, nil, err
	}

	combinedCancelFunc := func() {
		cancel()          // Cancel the chromedp context first.
		allocatorCancel() // Then cancel the allocator context.
	}

	return ctx, combinedCancelFunc, nil
}

func getRandomUserAgent() string {
	index := rand.Intn(len(userAgents))
	return userAgents[index]
}

func handleDirectDownload(a models.Addon, dp string) (done bool, err error) {
	log.Println("Because of special URL which has a '.zip' suffix, we will use DirectDownload!")

	destination := dp

	if !strings.HasSuffix(dp, string(os.PathSeparator)) {
		destination += string(os.PathSeparator) + a.Filename + ".zip"
	} else {
		destination += a.Filename + ".zip"
	}

	b, err := DownloadFile(a.DownloadUrl, destination)
	if err != nil {
		return false, err
	}

	if b {
		log.Printf("Done with download addon from, %s\n", a.DownloadUrl)
		return true, nil
	}

	log.Printf("Was unable to download addon from, %s\n", a.DownloadUrl)
	return false, nil
}

func handleDownloadLogging(db *sql.DB, a models.Addon, runId string) error {
	dLog, err := database.InsertDLog(db, models.DLog{RunId: runId,
		Url: a.DownloadUrl})
	if err != nil {
		return fmt.Errorf("failed to store download logging into database for the URL: %s -> %v\n",
			a.DownloadUrl, err)
	}

	log.Printf("Stored download logging with ID: %d, for the URL: %s\n", dLog, a.DownloadUrl)
	return nil
}

func navigateToAddonUrl(ctx context.Context, a models.Addon) {
	if err := chromedp.Run(ctx,
		chromedp.Navigate(a.DownloadUrl),
	); err != nil {
		log.Fatal(err)
	}
}

func setupContextOptions(ctx context.Context, dp string) error {
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			userAgent := getRandomUserAgent()
			log.Printf("Setting [%s] as user agent!\n", userAgent)

			// set user agent
			err := chromedp.Evaluate(fmt.Sprintf("navigator.userAgent = \"%s\"", userAgent),
				nil).Do(ctx)

			// set download behavior and path to downloadPath
			_ = browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow).
				WithDownloadPath(dp).
				WithEventsEnabled(true).
				Do(ctx)

			// enable ad-blocker
			params := page.SetAdBlockingEnabledParams{Enabled: true}
			_ = params.Do(ctx)
			return err
		}),
	)

	return err
}

func setupContextListeners(ctx context.Context, dp string, dc chan string) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *browser.EventDownloadWillBegin:
			log.Printf("Download from [%s] with GUID: [%s], will soon begin, it will be known as: %s\n",
				ev.URL, ev.GUID, ev.SuggestedFilename)
		case *browser.EventDownloadProgress:
			handleDownloadProgressEvent(ev, dc)
		case *page.EventDocumentOpened:
			handleDocumentOpenedEvent(ev, ctx, dp)
		}
	})
}

func handleDownloadProgressEvent(ev *browser.EventDownloadProgress, dc chan string) {
	if ev.State == browser.DownloadProgressStateInProgress {
		log.Printf("[%s] Progress: %f/%f}\n", ev.GUID, ev.ReceivedBytes, ev.ReceivedBytes)
	} else if ev.State == browser.DownloadProgressStateCanceled {
		log.Println("Download was cancelled!")
		dc <- ev.GUID
	} else if ev.State == browser.DownloadProgressStateCompleted {
		log.Println("Download completed!")
		dc <- ev.GUID
	} else {
		log.Printf("Unknown EventDownloadProgress state! [%s]\n", ev.State.String())
	}
}

func handleDocumentOpenedEvent(ev *page.EventDocumentOpened, ctx context.Context, dp string) {
	log.Printf("Document at [%s] has been opened, capturing screenshot...\n", ev.Frame.URL)

	// Execute the screenshot action in a separate goroutine to not block the event listener
	go func() {
		var buf []byte

		if err := chromedp.Run(ctx,
			chromedp.Sleep(time.Millisecond*200), // Ensure the page is ready
			chromedp.CaptureScreenshot(&buf),
		); err != nil {
			log.Fatal(err)
		}

		if err := os.WriteFile(dp+string(os.PathSeparator)+"screenshot.jpg", buf, 0644); err != nil {
			log.Fatal(err)
		}

		log.Println("Screenshot captured.")
	}()
}
