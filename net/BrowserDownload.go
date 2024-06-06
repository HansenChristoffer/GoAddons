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
}

func StartHeadlessAndDownloadAddons(runId string, addons []models.Addon, downloadPath string,
	db *sql.DB) (done bool, err error) {

	if addons == nil || downloadPath == "" {
		return false, fmt.Errorf("addons or downloadPath is required\n")
	}

	ctx, cancel, err := setupContext(downloadPath)
	if err != nil {
		return false, err
	}
	defer cancel()

	// Channel to signal the completion of each download
	downloadComplete := make(chan string)
	setupContextListeners(ctx, downloadPath, downloadComplete)

	for idx, addon := range addons {
		log.Printf("[%d/%d] Will try to navigate to: %s\n", idx+1, len(addons), addon.DownloadUrl)

		if addon.DownloadUrl == "" {
			log.Printf("DownloadURL is not allowed to be empty, will ignore this addon!\n")
		} else if strings.HasSuffix(addon.DownloadUrl, ".zip") {
			_, err := handleDirectDownload(addon, downloadPath, db, runId, downloadComplete)
			if err != nil {
				return false, err
			}
		} else {
			_, err := handleBrowserDownload(ctx, addon, db, runId, downloadComplete)
			if err != nil {
				return false, err
			}
		}

		// Trying to accommodate for any type of latency issues
		time.Sleep(1 * time.Second)
		log.Printf("Continuing to the next addon...\n")
	}
	log.Println("Download of all addons is now done!")
	return true, nil
}

func handleDirectDownload(addon models.Addon, downloadPath string, db *sql.DB,
	runId string, dcSignal chan string) (bool, error) {

	res, err := directDownload(addon, downloadPath)
	if err != nil {
		log.Printf("Error with direct download! -> %v\n", err)
	}

	if res {
		// Wait for the signal that the current download is complete
		<-dcSignal

		err = databaseOperations(db, addon, runId)
		if err != nil {
			return false, err
		}
		return res, nil
	}
	return false, nil
}

func handleBrowserDownload(ctx context.Context, addon models.Addon, db *sql.DB,
	runId string, dcSignal chan string) (bool, error) {

	navigateToAddonDownloadUrl(ctx, addon)
	// Wait for the signal that the current download is complete
	<-dcSignal

	err := databaseOperations(db, addon, runId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func databaseOperations(db *sql.DB, addon models.Addon, runId string) (err error) {
	err = handleDownloadLogging(db, addon, runId)
	if err != nil {
		return err
	}

	// Update the column 'last_downloaded'
	_, err = database.UpdateAddon(db, addon)
	if err != nil {
		return err
	}
	return nil
}

func setupContext(downloadPath string) (ctx context.Context, cancelFunc context.CancelFunc, err error) {
	allocatorCtx, allocatorCancel := chromedp.NewRemoteAllocator(context.Background(), dockerUrl)
	// initialize a controllable Chrome instance
	ctx, cancel := chromedp.NewContext(allocatorCtx)

	if err = setupContextOptions(ctx, downloadPath); err != nil {
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

func directDownload(addon models.Addon, downloadPath string) (done bool, err error) {
	log.Println("Because of special URL which has addon '.zip' suffix, we will use DirectDownload!")

	destination := downloadPath
	if !strings.HasSuffix(downloadPath, string(os.PathSeparator)) {
		destination += string(os.PathSeparator) + addon.Filename + ".zip"
	} else {
		destination += addon.Filename + ".zip"
	}

	result, err := DownloadFile(addon.DownloadUrl, destination)
	if err != nil {
		return false, err
	}

	if result {
		log.Printf("Done with download addon from, %s\n", addon.DownloadUrl)
		return true, nil
	}
	log.Printf("Was unable to download addon from, %s\n", addon.DownloadUrl)
	return false, nil
}

func handleDownloadLogging(db *sql.DB, addon models.Addon, runId string) error {
	dLog, err := database.InsertDLog(db, models.DLog{RunId: runId,
		Url: addon.DownloadUrl})
	if err != nil {
		return fmt.Errorf("failed to store download logging into database for the URL: %s -> %v\n",
			addon.DownloadUrl, err)
	}
	log.Printf("Stored download logging with ID: %d, for the URL: %s\n", dLog, addon.DownloadUrl)
	return nil
}

func navigateToAddonDownloadUrl(ctx context.Context, addon models.Addon) {
	if err := chromedp.Run(ctx,
		chromedp.Navigate(addon.DownloadUrl),
	); err != nil {
		log.Fatal(err)
	}
}

func setupContextOptions(ctx context.Context, downloadPath string) error {
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			userAgent := getRandomUserAgent()
			log.Printf("Setting [%s] as user agent!\n", userAgent)

			// set user agent
			err := chromedp.Evaluate(fmt.Sprintf("navigator.userAgent = \"%s\"", userAgent),
				nil).Do(ctx)

			// set download behavior and path to downloadPath
			_ = browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow).
				WithDownloadPath(downloadPath).
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

func setupContextListeners(ctx context.Context, downloadPath string, dcSignal chan string) {
	chromedp.ListenTarget(ctx, func(evt interface{}) {
		switch evt := evt.(type) {
		case *browser.EventDownloadWillBegin:
			log.Printf("Download from [%s] with GUID: [%s], will soon begin, it will be known as: %s\n",
				evt.URL, evt.GUID, evt.SuggestedFilename)
		case *browser.EventDownloadProgress:
			handleDownloadProgressEvent(evt, dcSignal)
		case *page.EventDocumentOpened:
			handleDocumentOpenedEvent(evt, ctx, downloadPath)
		}
	})
}

func handleDownloadProgressEvent(evt *browser.EventDownloadProgress, dcSignal chan string) {
	if evt.State == browser.DownloadProgressStateInProgress {
		log.Printf("[%s] Progress: %f/%f}\n", evt.GUID, evt.ReceivedBytes, evt.ReceivedBytes)
	} else if evt.State == browser.DownloadProgressStateCanceled {
		log.Println("Download was cancelled!")
		dcSignal <- evt.GUID
	} else if evt.State == browser.DownloadProgressStateCompleted {
		log.Println("Download completed!")
		dcSignal <- evt.GUID
	} else {
		log.Printf("Unknown EventDownloadProgress state! [%s]\n", evt.State.String())
	}
}

func handleDocumentOpenedEvent(evt *page.EventDocumentOpened, ctx context.Context, downloadPath string) {
	log.Printf("Document at [%s] has been opened, capturing screenshot...\n", evt.Frame.URL)

	// Execute the screenshot action in a separate goroutine to not block the event listener
	go func() {
		var buf []byte

		if err := chromedp.Run(ctx,
			chromedp.Sleep(time.Millisecond*200), // Ensure the page is ready
			chromedp.CaptureScreenshot(&buf),
		); err != nil {
			log.Fatal(err)
		}

		if err := os.WriteFile(downloadPath+string(os.PathSeparator)+"screenshot.jpg", buf, 0644); err != nil {
			log.Fatal(err)
		}
		log.Println("Screenshot captured.")
	}()
}
