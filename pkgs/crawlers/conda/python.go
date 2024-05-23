package conda

import (
	"encoding/json"
	"os"

	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewPython())
}

type Python struct {
	SDKName string
	CondaSearcher
}

func NewPython() (p *Python) {
	return &Python{
		SDKName: "python",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
}

func (p *Python) GetSDKName() string {
	return p.SDKName
}

func (p *Python) Start() {
	p.CondaSearcher.Search(p.SDKName)
}

func (p *Python) GetVersions() []byte {
	r, _ := p.Version.Marshal()
	return r
}

func (p *Python) HomePage() string {
	return "https://www.python.org/"
}

func TestPython() {
	pp := NewPython()
	pp.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/python.json"
	content, _ := json.MarshalIndent(pp.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
