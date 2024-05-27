package conda

import (
	"encoding/json"
	"os"

	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewGCC())
}

type GCC struct {
	SDKName string
	CondaSearcher
}

func NewGCC() (g *GCC) {
	g = &GCC{
		SDKName: "gcc",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (g *GCC) GetSDKName() string {
	return g.SDKName
}

func (g *GCC) Start() {
	g.CondaSearcher.Search(g.SDKName)
}

func (g *GCC) GetVersions() []byte {
	r, _ := g.Version.Marshal()
	return r
}

func (g *GCC) HomePage() string {
	return "https://gcc.gnu.org/"
}

func (g *GCC) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestGCC() {
	gg := NewGCC()
	gg.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/gcc.json"
	content, _ := json.MarshalIndent(gg.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
