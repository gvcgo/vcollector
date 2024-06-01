package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/gogf/gf/v2/util/gutil"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/conda"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/gh/lans"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/gh/lsp"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/gh/tools"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/mix"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/official"
	_ "github.com/gvcgo/vcollector/pkgs/crawlers/official/fixed"
)

func GenerateSDKHompePageListForDocs() {
	homePageList := map[string]string{}
	for _, cc := range crawler.CrawlerList {
		homePageList[cc.GetSDKName()] = cc.HomePage()
	}

	nameList := []string{}
	for name := range homePageList {
		nameList = append(nameList, name)
	}

	sort.Slice(nameList, func(i, j int) bool {
		return gutil.ComparatorString(nameList[i], nameList[j]) < 0
	})

	content := ""
	for _, name := range nameList {
		content += fmt.Sprintf("- [%s](%s)\n\n", name, homePageList[name])
	}

	os.WriteFile("/home/moqsien/golang/src/gvcgo/vcollector/docs_sdk_list.md", []byte(content), os.ModePerm)
}
