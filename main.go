package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/PuerkitoBio/goquery"
)

type ApiResponse struct {
    ResultCount int    `json:"ResultCount"`
    SearchParams struct {
        DateFrom string `json:"DateFrom"`
        DateTo   string `json:"DateTo"`
        Query    string `json:"Query"`
        Language string `json:"Language"`
        // ... (other fields in the SearchParams struct)
    } `json:"SearchParams"`
    Replies []struct {
        Pengumuman struct {
            Id2          string `json:"Id2"`
            // ... (other fields in the Pengumuman struct)
        } `json:"pengumuman"`
        Attachments []struct {
            Id          int    `json:"Id"`
            PDFFilename string `json:"PDFFilename"`
            // ... (other fields in the Attachments struct)
        } `json:"attachments"`
    } `json:"Replies"`
}
func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	ctxs, _ := chromedp.NewExecAllocator(context.Background(), opts...)

	// Initialize a controllable Chrome instance
	ctx, cancel := chromedp.NewContext(
		ctxs,
	)
	// To release the browser resources when
	// it is no longer needed
	defer cancel()

	var html string
	err := chromedp.Run(ctx,
		// Visit the target page
		chromedp.Navigate("https://idx.co.id/primary/ListedCompany/GetAnnouncement?kodeEmiten=&emitenType=*&indexFrom=0&pageSize=10&dateFrom=20240201&dateTo=20240209&lang=id&keyword=HMETD"),
		// Wait for the page to load
		chromedp.Sleep(10 * time.Second),
		// Extract the raw HTML from the page
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Select the root node on the page
			rootNode, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			html, err = dom.GetOuterHTML().WithNodeID(rootNode.NodeID).Do(ctx)
			return err
		}),
	)
	if err != nil {
		log.Fatal("Error while performing the automation logic:", err)
	}

	// Parse HTML and extract JSON
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal("Error while parsing HTML:", err)
	}

	var jsonStr string
	doc.Find("pre").Each(func(i int, s *goquery.Selection) {
		jsonStr = s.Text()
	})

	// Unmarshal the JSON data into the ApiResponse struct
	var apiResponse ApiResponse
	err = json.Unmarshal([]byte(jsonStr), &apiResponse)
	if err != nil {
		log.Fatal("Error while unmarshalling JSON:", err)
	}

	// Print the API response
	fmt.Println(apiResponse)
}