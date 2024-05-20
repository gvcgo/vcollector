package conda

import (
	"encoding/json"
	"os"

	"github.com/gvcgo/vcollector/pkgs/version"
)

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

func TestPyPy() {
	pp := NewPyPy()
	pp.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/pypy.json"
	content, _ := json.MarshalIndent(pp.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
