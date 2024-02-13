package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

type ApiResponse struct {
	ResultCount  int `json:"ResultCount"`
	SearchParams struct {
		DateFrom string `json:"DateFrom"`
		DateTo   string `json:"DateTo"`
		Query    string `json:"Query"`
		Language string `json:"Language"`
		// ... other fields
	} `json:"SearchParams"`
	Replies []struct {
		Announcement struct {
			Id2               string      `json:"Id2"`
			EfekEmiten_DIRE   bool        `json:"EfekEmiten_DIRE"`
			EfekEmiten_DINFRA bool        `json:"EfekEmiten_DINFRA"`
			Id                int         `json:"Id"`
			FinalId           interface{} `json:"FinalId"`
			OldFinalId        int         `json:"OldFinalId"`
			NoPengumuman      string      `json:"NoPengumuman"`
			TglPengumuman     string      `json:"TglPengumuman"`
			JudulPengumuman   string      `json:"JudulPengumuman"`
			JenisPengumuman   string      `json:"JenisPengumuman"`
			Kode_Emiten       string      `json:"Kode_Emiten"`
			CreatedDate       string      `json:"CreatedDate"`
			Form_Id           string      `json:"Form_Id"`
			PerihalPengumuman string      `json:"PerihalPengumuman"`
			JMSXGroupID       string      `json:"JMSXGroupID"`
			Divisi            string      `json:"Divisi"`
			KodeDivisi        string      `json:"KodeDivisi"`
			// ... other fields
		} `json:"pengumuman"`
		Attachments []struct {
			Id               int         `json:"Id"`
			PDFFilename      string      `json:"PDFFilename"`
			FullSavePath     string      `json:"FullSavePath"`
			JMSXGroupID      string      `json:"JMSXGroupID"`
			CorrelationID    interface{} `json:"CorrelationID"`
			IsAttachment     bool        `json:"IsAttachment"`
			OriginalFilename string      `json:"OriginalFilename"`
		} `json:"attachments"`
	} `json:"Replies"`
}

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	ctxs, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	// initialize a controllable Chrome instance
	ctx, cancel := chromedp.NewContext(
		ctxs,
	)
	// to release the browser resources when
	// it is no longer needed
	defer cancel()

	var html string
	err := chromedp.Run(ctx,
		// visit the target page
		chromedp.Navigate("https://idx.co.id/primary/ListedCompany/GetAnnouncement?kodeEmiten=&emitenType=*&indexFrom=0&pageSize=10&dateFrom=20240201&dateTo=20240209&lang=id&keyword=HMETD"),
		// wait for the page to load
		chromedp.Sleep(2000*time.Millisecond),
		// extract the raw HTML from the page
		chromedp.ActionFunc(func(ctx context.Context) error {
			// select the root node on the page
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

	// Remove HTML tags from the HTML string
	html = striphtml.StripTags(html)

	// Unmarshal the JSON data into the ApiResponse struct
	var apiResponse ApiResponse
	err = json.Unmarshal([]byte(html), &apiResponse)
	if err != nil {
		log.Fatal("Error while unmarshalling JSON:", err)
	}

	// Print the API response
	fmt.Println(apiResponse)
}
