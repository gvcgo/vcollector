package conda

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewTypst())
}

type Typst struct {
	SDKName string
	CondaSearcher
}

func NewTypst() (t *Typst) {
	t = &Typst{
		SDKName: "typst",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (t *Typst) GetSDKName() string {
	return t.SDKName
}

func (t *Typst) Start() {
	t.CondaSearcher.Search(t.SDKName)
}

func (t *Typst) GetVersions() []byte {
	r, _ := t.Version.Marshal()
	return r
}

func TestTypst() {
	tt := NewTypst()
	tt.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		tt.SDKName,
	)
	content, _ := json.MarshalIndent(tt.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
