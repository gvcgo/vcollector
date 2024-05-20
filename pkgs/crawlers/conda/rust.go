package conda

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gvcgo/vcollector/pkgs/version"
)

type Rust struct {
	SDKName string
	CondaSearcher
}

func NewRust() (r *Rust) {
	r = &Rust{
		SDKName: "rust",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (r *Rust) GetSDKName() string {
	return r.SDKName
}

func (r *Rust) Start() {
	r.CondaSearcher.Search(r.SDKName)
}

func TestRust() {
	rr := NewRust()
	rr.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		rr.SDKName,
	)
	content, _ := json.MarshalIndent(rr.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
