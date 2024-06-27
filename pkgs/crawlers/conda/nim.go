package conda

import (
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewNim())
	crawler.RegisterCondaCrawler(NewNim())
}

type Nim struct {
	SDKName string
	CondaSearcher
}

func NewNim() (n *Nim) {
	return &Nim{
		SDKName: "nim",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
}

func (n *Nim) GetSDKName() string {
	return n.SDKName
}

func (n *Nim) Start() {
	n.CondaSearcher.Search(n.SDKName)
}

func (n *Nim) GetVersions() []byte {
	r, _ := n.Version.Marshal()
	return r
}

func (n *Nim) HomePage() string {
	return "https://nim-lang.org/"
}

func (n *Nim) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}
