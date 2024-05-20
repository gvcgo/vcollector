package conda

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gvcgo/vcollector/pkgs/version"
)

type R struct {
	SDKName string
	CondaSearcher
}

func NewR() (r *R) {
	r = &R{
		SDKName: "r",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (r *R) GetSDKName() string {
	return r.SDKName
}

func (r *R) Start() {
	r.CondaSearcher.Search(r.SDKName)
}

func TestR() {
	rr := NewR()
	rr.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		rr.SDKName,
	)
	content, _ := json.MarshalIndent(rr.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
