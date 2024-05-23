package conda

import (
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
)

func init() {
	crawler.RegisterCrawler(NewCondaForgePackages())
}

const (
	CondaForgeSDKName = "conda-forge-pkgs"
)

/*
https://raw.githubusercontent.com/conda-forge/feedstock-outputs/single-file/feedstock-outputs.json
*/
type CondaForgePackages struct {
	DownloadUrl string
	SDKName     string
	result      []byte
}

func NewCondaForgePackages() (c *CondaForgePackages) {
	c = &CondaForgePackages{
		DownloadUrl: "https://raw.githubusercontent.com/conda-forge/feedstock-outputs/single-file/feedstock-outputs.json",
		SDKName:     CondaForgeSDKName,
	}
	return
}

func (c *CondaForgePackages) GetSDKName() string {
	return c.SDKName
}

func (c *CondaForgePackages) GetVersions() []byte {
	return c.result
}

func (c *CondaForgePackages) HomePage() string {
	return "https://conda-forge.org/"
}

func (c *CondaForgePackages) Start() {
	resp := req.GetResp(c.DownloadUrl, 180)
	c.result = []byte(resp)
}
