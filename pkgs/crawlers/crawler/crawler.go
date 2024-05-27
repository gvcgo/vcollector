package crawler

import "github.com/gvcgo/vcollector/internal/iconf"

type Crawler interface {
	Start()
	GetSDKName() string
	HomePage() string
	GetVersions() []byte
	GetInstallConf() (ic iconf.InstallerConfig)
}

var CrawlerList = []Crawler{}

func RegisterCrawler(c Crawler) {
	CrawlerList = append(CrawlerList, c)
}
