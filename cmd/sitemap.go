package cmd

import (
	"context"
	"fmt"

	sitemap "github.com/oxffaa/gopher-parse-sitemap"
	"github.com/spf13/cobra"

	"github.com/vallieres/crawl-n-indexnow/util"
)

const FirstTen = 10

func Sitemap() *cobra.Command {
	IndexNowCmd = &cobra.Command{
		Use:   "sitemap",
		Short: "Sends all of the Sitemap URLs to IndexNow.",
		Long: CrawlNIndexNowASCII + `

Gathers all of the website's URLs by parsing every single sitemap pages, 
packages them nicely and posts them to IndexNow's API.'.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeSitemap(cmd.Context())
		},
	}

	IndexNowCmd.PersistentFlags().StringVar(&domain, "domain", "", "the website's domain")
	IndexNowCmd.PersistentFlags().StringVar(&indexNowKey, "key", "", "the IndexNow key")

	cliContext := util.NewCLIContext(&domain, &indexNowKey)
	ctx := util.WithService(cliContext, &domain, &indexNowKey)
	IndexNowCmd.SetContext(ctx)

	return IndexNowCmd
}

func executeSitemap(ctx context.Context) error {
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

	listURLs, errGetURLs := GetListSitemapURLs(*domainCtx)
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

func GetListSitemapURLs(domain string) ([]string, error) {
	pageURLs := make([]string, 0)

	errGetUrls := sitemap.ParseFromSite("https://"+domain+"/sitemap.xml", func(e sitemap.Entry) error {
		pageURLs = append(pageURLs, e.GetLocation())
		return nil
	})
	if errGetUrls != nil {
		CPrintError("error parsing Sitemap URLs: ", errGetUrls)
	}

	CPrint("Total of Sitemap URLs: ", len(pageURLs))
	PrintBlankLine()

	CPrint("List of first 10 page URLs found")
	for i, url := range pageURLs {
		if i > FirstTen {
			break
		}
		CPrint("\t", url)
	}
	PrintBlankLine()

	return pageURLs, nil
}
