package conda

import (
	"encoding/json"
	"os"

	"github.com/gvcgo/vcollector/pkgs/version"
)

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

func TestPython() {
	pp := NewPython()
	pp.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/python.json"
	content, _ := json.MarshalIndent(pp.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
