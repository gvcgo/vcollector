package conda

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewR())
	crawler.RegisterCondaCrawler(NewR())
}

type R struct {
	SDKName string
	CondaSearcher
}

func NewR() (r *R) {
	r = &R{
		SDKName: "r",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (r *R) GetSDKName() string {
	return r.SDKName
}

func (r *R) Start() {
	r.CondaSearcher.Search(r.SDKName)
}

func (r *R) GetVersions() []byte {
	result, _ := r.Version.Marshal()
	return result
}

func (r *R) HomePage() string {
	return "https://www.r-project.org/"
}

func (r *R) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"Scripts"}, {"Library", "bin"}},
			MacOS:   []iconf.DirPath{{"bin"}, {"sbin"}},
			Linux:   []iconf.DirPath{{"bin"}, {"sbin"}},
		},
	}
}

func TestR() {
	rr := NewR()
	rr.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		rr.SDKName,
	)
	content, _ := json.MarshalIndent(rr.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
