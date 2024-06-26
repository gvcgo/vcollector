package conda

import (
	"encoding/json"
	"os"

	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewRuby())
	crawler.RegisterCondaCrawler(NewRuby())
}

type Ruby struct {
	SDKName string
	CondaSearcher
}

func NewRuby() (r *Ruby) {
	r = &Ruby{
		SDKName: "ruby",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (r *Ruby) GetSDKName() string {
	return r.SDKName
}

func (r *Ruby) Start() {
	r.CondaSearcher.Search(r.SDKName)
}

func (r *Ruby) GetVersions() []byte {
	result, _ := r.Version.Marshal()
	return result
}

func (r *Ruby) HomePage() string {
	return "https://www.ruby-lang.org/en/"
}

func (r *Ruby) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestRuby() {
	rr := NewRuby()
	rr.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/ruby.json"
	content, _ := json.MarshalIndent(rr.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
