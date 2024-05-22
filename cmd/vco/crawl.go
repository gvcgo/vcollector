package main

import (
	"fmt"

	"github.com/gvcgo/vcollector/internal/utils"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/conda"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/gh/lans"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/gh/lsp"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/gh/tools"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/mix"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/official"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/official/fixed"
)

/*
1. start crawlers.
2. upload files.
TODO: multi goroutine
*/
func Crawl() {
	uploader := utils.NewUploader()
	for _, cc := range crawler.CrawlerList {
		fmt.Printf("crawling %s\n", cc.GetSDKName())
		cc.Start()
		uploader.Upload(cc.GetSDKName(), cc.GetVersions())
	}
}
