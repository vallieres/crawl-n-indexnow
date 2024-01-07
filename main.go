package main

import (
	"context"
	"fmt"
	"os"

	"github.com/vallieres/crawl-n-indexnow/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var exitCode int
	defer func() {
		if r := recover(); r != nil {
			exitCode = 1
			fmt.Fprintln(os.Stderr, "Panic: Recovered in main: ", r)
		}
		os.Exit(exitCode)
	}()

	ctx := context.Background()

	root := cmd.Root(version, commit, date)

	root.AddCommand(
		cmd.Shopify(),
		cmd.Sitemap(),
	)

	if err := root.ExecuteContext(ctx); err != nil {
		fmt.Println("failed to run command", err.Error())
		exitCode = 1
	}
}
