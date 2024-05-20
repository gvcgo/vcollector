package crawler

type Crawler interface {
	Start()
	GetSDKName() string
}

var CrawlerList = []Crawler{}

func RegisterCrawler(c Crawler) {
	CrawlerList = append(CrawlerList, c)
}
