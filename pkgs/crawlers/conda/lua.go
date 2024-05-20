package conda

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gvcgo/vcollector/pkgs/version"
)

type Lua struct {
	SDKName string
	CondaSearcher
}

func NewLua() (l *Lua) {
	l = &Lua{
		SDKName: "lua",
		CondaSearcher: CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (l *Lua) GetSDKName() string {
	return l.SDKName
}

func (l *Lua) Start() {
	l.CondaSearcher.Search(l.SDKName)
}

func TestLua() {
	ll := NewLua()
	ll.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		ll.SDKName,
	)
	content, _ := json.MarshalIndent(ll.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
