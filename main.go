package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	colly "github.com/gocolly/colly/v2"
)

func main() {
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	ctx := context.Background()

	fmt.Println(`
  _____ ___   ___  _      __ __         _   ____ _  __ ___   ____ _  __
 / ___// _ \ / _ || | /| / // /    ___ ( ) /  _// |/ // _ \ / __/| |/_/
/ /__ / , _// __ || |/ |/ // /__  / _ \|/ _/ / /    // // // _/ _>  <  
\___//_/|_|/_/ |_||__/|__//____/ /_//_/  /___//_/|_//____//___//_/|_|`)
	PrintBlankLine()

	var knownUrls []string

	c := colly.NewCollector(colly.AllowedDomains("legacygoods.co"))
	c.OnXML("//sitemap/loc", func(e *colly.XMLElement) {
		knownUrls = append(knownUrls, e.Text)
	})

	// Start the collector
	errVisit := c.Visit("https://legacygoods.co/sitemap.xml")
	if errVisit != nil {
		CPrintError("error visiting main sitemap: ", errVisit)
	}

	CPrint("Total of Sitemaps: ", len(knownUrls))
	for _, url := range knownUrls {
		CPrint("\t", url)
	}
	PrintBlankLine()

	var pageUrls []string

	d := colly.NewCollector(colly.AllowedDomains("legacygoods.co"))
	d.OnXML("//url/loc", func(e *colly.XMLElement) {
		pageUrls = append(pageUrls, e.Text)
	})

	for i := range knownUrls {
		err := d.Visit(knownUrls[i])
		if err != nil {
			CPrintError("error visiting Sitemaps URLs: ", err)
		}
	}

	CPrint("List of all page URLs found")
	for _, url := range pageUrls {
		CPrint("\t", url)
	}
	PrintBlankLine()

	CPrint("Sending", len(pageUrls), "URLs to IndexNow...")
	code, response, errPost := POSTtoIndexNow(ctx, pageUrls)
	if errPost != nil {
		CPrintError(errPost)
	}

	CPrint("Code     :", code)
	CPrint("Response :", response)
}

type IndexNowRequestBody struct {
	Host        string   `json:"host"`
	Key         string   `json:"key"`
	KeyLocation string   `json:"keyLocation"`
	URLList     []string `json:"urlList"`
}

func POSTtoIndexNow(ctx context.Context, pageURLs []string) (string, string, error) {
	requestBody := IndexNowRequestBody{
		Host:        "legacygoods.co",
		Key:         "16e3f6d366b84fd986280e026b9fbd74",
		KeyLocation: "https://legacygoods.co/16e3f6d366b84fd986280e026b9fbd74.txt",
		URLList:     pageURLs,
	}
	requestBodyBytes, errMarshal := json.Marshal(requestBody)
	if errMarshal != nil {
		CPrintError("Error:", errMarshal)
		return "", "", errMarshal
	}

	body := bytes.NewBuffer(requestBodyBytes)

	// Create client
	client := &http.Client{}

	// Create request
	req, errNewReq := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.indexnow.org/IndexNow", body)
	if errNewReq != nil {
		CPrintError("error crafting new request:", errNewReq)
		return "", "", errNewReq
	}

	// Headers
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		CPrintError("error parsing request body bytes:", parseFormErr)
		return "", "", parseFormErr
	}

	// Fetch Request
	resp, errExec := client.Do(req)
	if errExec != nil {
		CPrintError("error posting data:", errExec)
		return "", "", errExec
	}
	defer resp.Body.Close()

	// Read Response Body
	respBody, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		CPrintError("error reading response data:", errRead)
		return "", "", errRead
	}

	/**
	HTTP Code	Response				Reasons
	---------	----------------------	--------------------------------------------------------------
	200			Ok						RL submitted successfully
	400			Bad request				Invalid format
	403			Forbidden				In case of key not valid (e.g. key not found, file found but key not in the file)
	422			Unprocessable Entity	In case of URLs don’t belong to the host or the key is not matching the schema in the protocol
	429			Too Many Requests		Too Many Requests (potential Spam)
	*/

	return resp.Status, string(respBody), nil
}

func CPrint(msg ...interface{}) {
	prefix := make([]interface{}, 0, 1)
	prefix = append(prefix, "[-]")
	msg = append(prefix, msg...)
	fmt.Println(msg...)
}

func PrintBlankLine() {
	fmt.Println("")
}

func CPrintError(msg ...interface{}) {
	prefix := make([]interface{}, 0, 1)
	prefix = append(prefix, "[❗]")

	msg = append(prefix, msg...)
	if _, err := fmt.Fprintln(os.Stderr, msg...); err != nil {
		fmt.Printf("error writing to STDERR: %v\n", err)
		os.Exit(1)
	}
}
