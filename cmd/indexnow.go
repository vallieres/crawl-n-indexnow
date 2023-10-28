package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/spf13/cobra"

	"github.com/vallieres/crawl-n-indexnow/models"
	"github.com/vallieres/crawl-n-indexnow/util"
)

var (
	IndexNowCmd         *cobra.Command
	domain, indexNowKey string
)

const (
	Ten   = 10
	Sixty = 60
)

func IndexNow() *cobra.Command {
	IndexNowCmd = &cobra.Command{
		Use:   "indexnow",
		Short: "Sends all of the Shopify's URLs to IndexNow.",
		Long: CrawlNIndexNowASCII + `

Gathers all of the Shopify's URL by parsing every single sitemap pages, 
packages them nicely and posts them to IndexNow's API.'.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeIndexNow(cmd.Context())
		},
	}

	IndexNowCmd.PersistentFlags().StringVar(&domain, "domain", "", "the Shopify domain")
	IndexNowCmd.PersistentFlags().StringVar(&indexNowKey, "key", "", "the IndexNow key")

	cliContext := util.NewCLIContext(&domain, &indexNowKey)
	ctx := util.WithService(cliContext, &domain, &indexNowKey)
	IndexNowCmd.SetContext(ctx)

	return IndexNowCmd
}

func executeIndexNow(ctx context.Context) error {
	fmt.Println(CrawlNIndexNowASCII) //nolint:govet

	domainCtx, errGetDomain := util.GetDomain(ctx)
	if errGetDomain != nil {
		return fmt.Errorf("error pulling the Domain from the context: %w", errGetDomain)
	}
	if *domainCtx == "" {
		return fmt.Errorf("domain is required to execute this command")
	}

	indexNowKeyCtx, errGetIndexNowKey := util.GetIndexNowKey(ctx)
	if errGetIndexNowKey != nil {
		return fmt.Errorf("error pulling the IndexNowKey from the context: %w", errGetIndexNowKey)
	}
	if *indexNowKeyCtx == "" {
		return fmt.Errorf("indexNowKey is required to execute this command")
	}

	listURLs, errGetURLs := GetListOfShopifyURLs(*domainCtx)
	if errGetURLs != nil {
		CPrintError(errGetURLs)
	}

	CPrint("Sending", len(listURLs), "URLs to IndexNow...")
	code, response, errPost := POSTtoIndexNow(*domainCtx, *indexNowKeyCtx, listURLs)
	if errPost != nil {
		CPrintError(errPost)
	}

	CPrint("Code     :", code)
	CPrint("Response :", response)

	return nil
}

func GetListOfShopifyURLs(domain string) ([]string, error) {
	var knownUrls []string

	c := colly.NewCollector(colly.AllowedDomains(domain))
	c.OnXML("//sitemap/loc", func(e *colly.XMLElement) {
		knownUrls = append(knownUrls, e.Text)
	})

	// Start the collector
	errVisit := c.Visit("https://" + domain + "/sitemap.xml")
	if errVisit != nil {
		CPrintError("error visiting main sitemap: ", errVisit)
	}

	CPrint("Total of Sitemaps: ", len(knownUrls))
	for _, url := range knownUrls {
		CPrint("\t", url)
	}
	PrintBlankLine()

	var pageUrls []string

	d := colly.NewCollector(colly.AllowedDomains(domain))
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

	return pageUrls, nil
}

func POSTtoIndexNow(domain string, indexNowKey string, pageURLs []string) (string, string, error) {
	ctx, cancelFnc := context.WithTimeout(context.Background(), Sixty*time.Second)
	defer cancelFnc()

	requestBody := models.IndexNowRequestBody{
		Host:        domain,
		Key:         indexNowKey,
		KeyLocation: "https://" + domain + "/" + indexNowKey + ".txt",
		URLList:     pageURLs,
	}
	requestBodyBytes, errMarshal := json.Marshal(requestBody)
	if errMarshal != nil {
		CPrintError("Error:", errMarshal)
		return "", "", errMarshal
	}

	body := bytes.NewBuffer(requestBodyBytes)

	// Create client
	client := &http.Client{
		Timeout: Ten * time.Second,
	}

	// Create request
	req, errNewReq := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.indexnow.org/IndexNow", body)
	if errNewReq != nil {
		CPrintError("error crafting new request:", errNewReq)
		return "", "", errNewReq
	}

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
