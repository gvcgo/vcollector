package conda

import (
	"encoding/json"
	"os"

	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewClang())
}

type Clang struct {
	SDKName string
	CondaSearcher
}

func NewClang() (c *Clang) {
	c = &Clang{
		SDKName: "clang",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (c *Clang) GetSDKName() string {
	return c.SDKName
}

func (c *Clang) Start() {
	c.CondaSearcher.Search(c.SDKName)
}

func (c *Clang) GetVersions() []byte {
	r, _ := c.Version.Marshal()
	return r
}

func TestClang() {
	cc := NewClang()
	cc.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/clang.json"
	content, _ := json.MarshalIndent(cc.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
