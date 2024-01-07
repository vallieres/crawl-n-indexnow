package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const CrawlNIndexNowASCII = `
  _____ ___   ___  _      __ __         _   ____ _  __ ___   ____ _  __ _  __ ____  _      __
 / ___// _ \ / _ || | /| / // /    ___ ( ) /  _// |/ // _ \ / __/| |/_// |/ // __ \| | /| / /
/ /__ / , _// __ || |/ |/ // /__  / _ \|/ _/ / /    // // // _/ _>  < /    // /_/ /| |/ |/ /
\___//_/|_|/_/ |_||__/|__//____/ /_//_/  /___//_/|_//____//___//_/|_|/_/|_/ \____/ |__/|__/

`

const IndexNowEndpoint = "https://api.indexnow.org/IndexNow"

func Root(version string, commit string, date string) *cobra.Command {
	if date != "unknown" {
		t, err := time.Parse(time.RFC3339, date)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		date = t.Format("2006-01-02")
	}
	rootCmd := &cobra.Command{
		Use:     "crawl-n-indexnow",
		Version: fmt.Sprintf("%s (%s) [%s]", version, date, commit),
		Short:   "Get the goods and ship 'em to the indexes!",
		Long: CrawlNIndexNowASCII + `
Crawl n' Index is a simple CLI that pulls your Shopify site's URL, and 
submits them to various indexes to speed up the indexing process.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.SilenceUsage = true

	return rootCmd
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
