package conda

import (
	"encoding/json"
	"os"

	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewPyPy())
}

type PyPy struct {
	SDKName string
	CondaSearcher
}

func NewPyPy() (p *PyPy) {
	p = &PyPy{
		SDKName: "pypy",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (p *PyPy) GetSDKName() string {
	return p.SDKName
}

func (p *PyPy) Start() {
	p.CondaSearcher.Search(p.SDKName)
}

func (p *PyPy) GetVersions() []byte {
	r, _ := p.Version.Marshal()
	return r
}

func (p *PyPy) HomePage() string {
	return "https://www.pypy.org/"
}

func (p *PyPy) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		BinaryDirs: &iconf.FileItems{
			Windows: []string{"bin"},
			MacOS:   []string{"bin"},
			Linux:   []string{"bin"},
		},
	}
}

func TestPyPy() {
	pp := NewPyPy()
	pp.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/pypy.json"
	content, _ := json.MarshalIndent(pp.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
