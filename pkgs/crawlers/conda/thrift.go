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
	crawler.RegisterCrawler(NewThrift())
}

type Thrif struct {
	SDKName string
	CondaSearcher
}

func NewThrift() (t *Thrif) {
	t = &Thrif{
		SDKName: "thrift-compiler",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (t *Thrif) GetSDKName() string {
	return t.SDKName
}

func (t *Thrif) Start() {
	t.CondaSearcher.Search(t.SDKName)
}

func (t *Thrif) GetVersions() []byte {
	r, _ := t.Version.Marshal()
	return r
}

func (t *Thrif) HomePage() string {
	return "https://thrift.apache.org/"
}

func (t *Thrif) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestThrift() {
	tt := NewThrift()
	tt.Start()

	ff := fmt.Sprintf(
		"/home/moqsien/projects/go/src/gvcgo/vcollector/test/%s.json",
		tt.SDKName,
	)
	content, _ := json.MarshalIndent(tt.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
