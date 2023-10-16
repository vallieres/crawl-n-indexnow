package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "crawl-n-index",
		Short: "Get the goods and ship 'em to the indexes!",
		Long: `
  _____ ___   ___  _      __ __         _   ____ _  __ ___   ____ _  __
 / ___// _ \ / _ || | /| / // /    ___ ( ) /  _// |/ // _ \ / __/| |/_/
/ /__ / , _// __ || |/ |/ // /__  / _ \|/ _/ / /    // // // _/ _>  <
\___//_/|_|/_/ |_||__/|__//____/ /_//_/  /___//_/|_//____//___//_/|_|

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

func CPrintNNL(msg ...interface{}) {
	prefix := make([]interface{}, 0, 1)
	prefix = append(prefix, "[-]")

	msg = append(prefix, msg...)
	fmt.Print(msg...)
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
