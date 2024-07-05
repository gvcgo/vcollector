package conda

import (
	"encoding/json"
	"os"

	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/version"
)

/*
conda g++
*/

type GXX struct {
	SDKName string
	CondaSearcher
}

func NewGXX() (g *GXX) {
	g = &GXX{
		SDKName: "gxx",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (g *GXX) GetSDKName() string {
	return g.SDKName
}

func (g *GXX) Start() {
	g.CondaSearcher.Search(g.SDKName)
}

func (g *GXX) GetVersions() []byte {
	r, _ := g.Version.Marshal()
	return r
}

func (g *GXX) HomePage() string {
	return "https://gcc.gnu.org/"
}

func (g *GXX) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestGXX() {
	gg := NewGXX()
	gg.Start()

	ff := "/home/moqsien/projects/go/src/gvcgo/vcollector/test/gxx.json"
	content, _ := json.MarshalIndent(gg.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
