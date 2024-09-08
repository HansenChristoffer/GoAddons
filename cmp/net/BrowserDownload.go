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
	"fmt"
	"goaddons/cmp/models"
	"goaddons/cmp/system"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func StartHeadlessAndDownloadAddons(addonVault *models.AddonVault, downloadPath string) (done bool, err error) {
	if addonVault.Length() == 0 || downloadPath == "" {
		return false, fmt.Errorf("browserDownload -> StartHeadlessAndDownloadAddons(): addons or downloadPath is required\n")
	}

	ctx, cancel, err := setupContext(downloadPath)
	if err != nil {
		return false, fmt.Errorf("browserDownload -> StartHeadlessAndDownloadAddons(): %w\n", err)
	}
	defer cancel()

	// Channel to signal the completion of each download
	downloadComplete := make(chan string)
	setupContextListeners(ctx, downloadComplete)

	for idx, addon := range addonVault.Addons {
		log.Printf("[%d/%d] Will try to navigate to: %s\n", idx+1, len(addonVault.Addons), addon.Url)

		if addon.Url == "" {
			log.Printf("URL is not allowed to be empty, will ignore this addon!\n")
		} else if strings.HasSuffix(addon.Url, ".zip") {
			_, err := handleDirectDownload(addon, downloadPath, downloadComplete)
			if err != nil {
				return false, err
			}
		} else {
			_, err := handleBrowserDownload(ctx, addon, downloadComplete)
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

func handleDirectDownload(addon models.Addon, downloadPath string, dcSignal chan string) (bool, error) {
	res, err := directDownload(addon, downloadPath)
	if err != nil {
		log.Printf("Error with direct download! -> %v\n", err)
	}

	if res {
		// Wait for the signal that the current download is complete
		<-dcSignal

		log.Printf("Download finished!\n")
		return res, nil
	}
	return false, nil
}

func handleBrowserDownload(ctx context.Context, addon models.Addon, dcSignal chan string) (bool, error) {
	navigateToAddonUrl(ctx, addon)
	// Wait for the signal that the current download is complete
	<-dcSignal

	log.Printf("Download finished!\n")
	return true, nil
}

func setupContext(downloadPath string) (ctx context.Context, cancelFunc context.CancelFunc, err error) {
	systemConfig, err := models.GetSystemVaultInstance()
	if err != nil {
		return nil, nil, fmt.Errorf("setupContext: unable to get system vault instance: %w", err)
	}

	allocatorCtx, allocatorCancel := chromedp.NewRemoteAllocator(context.Background(),
		systemConfig.GetConfigByName("docker.url").Value)

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

func directDownload(addon models.Addon, downloadPath string) (done bool, err error) {
	log.Println("Because of special URL which has addon '.zip' suffix, we will use DirectDownload!")

	destination := downloadPath
	if !strings.HasSuffix(downloadPath, string(os.PathSeparator)) {
		destination += string(os.PathSeparator) + addon.Name + ".zip"
	} else {
		destination += addon.Name + ".zip"
	}

	result, err := DownloadFile(addon.Url, destination)
	if err != nil {
		return false, err
	}

	if result {
		log.Printf("Done with download addon from, %s\n", addon.Url)
		return true, nil
	}
	log.Printf("Was unable to download addon from, %s\n", addon.Url)
	return false, nil
}

func navigateToAddonUrl(ctx context.Context, addon models.Addon) {
	if err := chromedp.Run(ctx,
		chromedp.Navigate(addon.Url),
	); err != nil {
		log.Fatal(err)
	}
}

func setupContextOptions(ctx context.Context, downloadPath string) error {
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			agents, err := system.GetUserAgentConfig()
			if err != nil {
				return err
			}
			userAgent := models.RandomUserAgent(agents)
			log.Printf("Setting [%s] as user agent, with value: %s\n", userAgent.Id, userAgent.Value)

			// set user agent
			err = chromedp.Evaluate(fmt.Sprintf("navigator.userAgent = \"%s\"", userAgent.Value),
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

func setupContextListeners(ctx context.Context, dcSignal chan string) {
	chromedp.ListenTarget(ctx, func(evt interface{}) {
		switch evt := evt.(type) {
		case *browser.EventDownloadWillBegin:
			log.Printf("Download from [%s] with GUID: [%s], will soon begin, it will be known as: %s\n",
				evt.URL, evt.GUID, evt.SuggestedFilename)
		case *browser.EventDownloadProgress:
			handleDownloadProgressEvent(evt, dcSignal)
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
