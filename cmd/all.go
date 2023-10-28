package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vallieres/crawl-n-indexnow/util"
)

var AllCmd *cobra.Command

func All() *cobra.Command {
	AllCmd = &cobra.Command{
		Use:   "all",
		Short: "Sends all of the Shopify's URLs to every single indexes.",
		Long: CrawlNIndexNowASCII + ` 

Gathers all of the Shopify's URL by parsing every single sitemap pages, 
packages them nicely and posts them to every index on file.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeAll(cmd.Context())
		},
	}

	AllCmd.PersistentFlags().StringVar(&domain, "domain", "", "the Shopify domain")
	AllCmd.PersistentFlags().StringVar(&indexNowKey, "key", "", "the IndexNow key")

	cliContext := util.NewCLIContext(&domain, &indexNowKey)
	ctx := util.WithService(cliContext, &domain, &indexNowKey)
	AllCmd.SetContext(ctx)

	return AllCmd
}

func executeAll(ctx context.Context) error {
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
