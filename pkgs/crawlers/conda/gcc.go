package conda

import (
	"encoding/json"
	"os"

	"github.com/gvcgo/vcollector/pkgs/version"
)

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

func TestGCC() {
	gg := NewGCC()
	gg.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/gcc.json"
	content, _ := json.MarshalIndent(gg.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
