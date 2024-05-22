package crawler

type Crawler interface {
	Start()
	GetSDKName() string
	GetVersions() []byte
}

var CrawlerList = []Crawler{}

func RegisterCrawler(c Crawler) {
	CrawlerList = append(CrawlerList, c)
}
