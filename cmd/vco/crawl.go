package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gvcgo/vcollector/internal/conf"
	"github.com/gvcgo/vcollector/internal/utils"
	"github.com/gvcgo/vcollector/pkgs/crawlers/conda"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/gh/lans"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/gh/lsp"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/gh/tools"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/mix"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/official"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/official/fixed"
)

const SDKListFileName = "sdk-homepage.json"

var (
	sender chan crawler.Crawler
	wg     *sync.WaitGroup = &sync.WaitGroup{}
)

func send() {
	for _, cc := range crawler.CrawlerList {
		sender <- cc
	}
	close(sender)
}

func runCrawler(cc crawler.Crawler) {
	if cc == nil {
		return
	}
	fmt.Println("start crawler:", cc.GetSDKName())
	cc.Start()
	uploader := utils.NewUploader()
	if cc.GetSDKName() == conda.CondaForgeSDKName {
		uploader.DisableSaveSha256()
	}
	uploader.Upload(cc.GetSDKName(), cc.GetVersions())
}

func crawl() {
OUTTER:
	for {
		select {
		case cc, ok := <-sender:
			if !ok && cc == nil {
				wg.Done()
				break OUTTER
			}
			if cc != nil {
				runCrawler(cc)
			}
		default:
			time.Sleep(time.Microsecond * 100)
		}
	}
}

func RunMultiGoroutine() {
	sender = make(chan crawler.Crawler, 10)
	go send()
	time.Sleep(time.Millisecond * 500)
	// multi goroutines.
	num := 1
	for i := 0; i < num; i++ {
		wg.Add(1)
		go crawl()
	}
	wg.Wait()
}

func RunSingleGoroutine() (hList map[string]string) {
	homepageList := map[string]string{}
	for _, cc := range crawler.CrawlerList {
		runCrawler(cc)
		homepageList[cc.GetSDKName()] = cc.HomePage()
	}
	return homepageList
}

/*
1. start crawlers.
2. upload files.
*/
func start() {
	hList := RunSingleGoroutine()
	// upload sdklist file.
	fPath := filepath.Join(conf.GetWorkDir(), utils.ShaFileName)
	content, _ := os.ReadFile(fPath)
	upl := utils.NewUploader()
	if len(content) > 0 {
		upl.DisableSaveSha256()
		upl.Upload("sdk-list", content)
	}
	if len(hList) > 0 {
		upl.DisableSaveSha256()
		content, _ := json.MarshalIndent(hList, "", "  ")
		upl.Upload(SDKListFileName, content)
	}
}
