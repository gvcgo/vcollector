package tools

import (
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

type Asciinema struct {
	SDkName  string
	RepoName string
	searcher.GhSearcher
}

func NewAsciinema() (a *Asciinema) {
	a = &Asciinema{
		SDkName:  "acast",
		RepoName: "gvcgo/asciinema",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (a *Asciinema) GetSDkName() string {
	return a.SDkName
}

func (a *Asciinema) Start() {
}
