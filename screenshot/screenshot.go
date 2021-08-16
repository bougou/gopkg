package screenshot

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// TakeScreenShotToImageFile takes a screenshot of the specified url and save to filename in specified imgFormat.
// supported format: "jpeg", "png", "webp"
func TakeScreenShotToImageFile(url string, filename string, imgFormat string) {

	// Start Chrome
	// Remove the 2nd param if you don't need debug information logged
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Printf))
	defer cancel()

	// Run Tasks
	// List of actions to run in sequence (which also fills our image buffer)
	var imageBuf []byte
	if err := chromedp.Run(ctx, fullScreenshotTasks2(url, &imageBuf, "png")); err != nil {
		log.Fatal(err)
	}

	// Write our image to file
	if err := ioutil.WriteFile(filename, imageBuf, 0644); err != nil {
		log.Fatal(err)
	}

}

func screenshotTasks(url string, imageBuf *[]byte, imageFormat string) chromedp.Tasks {
	wait := 2 * time.Second
	//	imgFormat := page.CaptureScreenshotFormat(imageFormat)
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.EmulateViewport(1280, 1280, chromedp.EmulateOrientation(emulation.OrientationTypePortraitPrimary, 30)),
		chromedp.Sleep(wait),
		chromedp.FullScreenshot(imageBuf, 90),
		// chromedp.ActionFunc(func(ctx context.Context) (err error) {
		// 	*imageBuf, err = page.CaptureScreenshot().WithQuality(90).WithFormat(imgFormat).Do(ctx)
		// 	return err
		// }),
	}
}

func screenshotTasksToPDF(url string, imageBuf *[]byte, imageFormat string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) (err error) {
			*imageBuf, _, err = page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			return err
		}),
	}
}

func fullScreenshotTasks2(url string, imageBuf *[]byte, imageFormat string) chromedp.Tasks {

	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(5 * time.Second),
		chromedp.EmulateViewport(1280, 640),
		showWidthHeight(),
		chromedp.FullScreenshot(imageBuf, 90),
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Reset
func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(10 * time.Second),
		chromedp.FullScreenshot(res, quality),
	}
}

func scrollToBottom() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		_, exp, err := runtime.Evaluate(`window.scrollTo(0,document.body.scrollHeight);`).Do(ctx)
		if err != nil {
			return err
		}
		if exp != nil {
			return exp
		}
		return nil
	})
}

func showWidthHeight() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// get layout metrics
		_, _, contentSize, _, _, cssContentSize, err := page.GetLayoutMetrics().Do(ctx)
		if err != nil {
			fmt.Println(err)
		}

		if cssContentSize != nil {
			contentSize = cssContentSize
		}
		width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

		fmt.Println("viewport width:", width)
		fmt.Println("viewport height:", height)
		return nil
	})
}

// elementScreenshot takes a screenshot of the first element node matching the selector.
func elementScreenshot(urlstr, selector string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(selector, res, chromedp.NodeVisible),
	}
}
