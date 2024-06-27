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
	crawler.RegisterCrawler(NewLFortran())
	crawler.RegisterCondaCrawler(NewLFortran())
}

type LFortran struct {
	SDKName string
	CondaSearcher
}

func NewLFortran() (l *LFortran) {
	l = &LFortran{
		SDKName: "lfortran",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (l *LFortran) GetSDKName() string {
	return l.SDKName
}

func (l *LFortran) Start() {
	l.CondaSearcher.Search(l.SDKName)
}

func (l *LFortran) GetVersions() []byte {
	r, _ := l.Version.Marshal()
	return r
}

func (l *LFortran) HomePage() string {
	return "https://lfortran.org/"
}

func (l *LFortran) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"Library", "bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestLFortran() {
	ll := NewLFortran()
	ll.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		ll.SDKName,
	)
	content, _ := json.MarshalIndent(ll.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
