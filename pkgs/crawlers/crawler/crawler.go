package crawler

import "github.com/gvcgo/vcollector/pkgs/version"

type Crawler interface {
	Start()
	GetSDKName() string
	GetVersions() version.VersionList
}

var CrawlerList = []Crawler{}

func RegisterCrawler(c Crawler) {
	CrawlerList = append(CrawlerList, c)
}
