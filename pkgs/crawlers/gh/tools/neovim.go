package tools

import (
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewNeovim())
}

type Neovim struct {
	SDKName  string
	RepoName string
	searcher.GhSearcher
}

func NewNeovim() (n *Neovim) {
	n = &Neovim{
		SDKName:  "neovim",
		RepoName: "neovim/neovim",
		GhSearcher: searcher.GhSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (n *Neovim) GetSDKName() string {
	return n.SDKName
}

func (n *Neovim) HomePage() string {
	return "https://neovim.io/"
}

func (n *Neovim) tagFilter(ri gh.ReleaseItem) bool {
	return searcher.GhVersionRegexp.FindString(ri.TagName) != ""
}

func (n *Neovim) fileFilter(a gh.Asset) bool {
	if strings.Contains(a.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(a.Name, ".sha256sum") {
		return false
	}
	if strings.HasSuffix(a.Name, ".msi") {
		return false
	}
	if strings.HasSuffix(a.Name, ".appimage") {
		return false
	}
	if strings.HasSuffix(a.Name, ".zsync") {
		return false
	}
	return true
}

func (n *Neovim) osParser(fName string) (osStr string) {
	if fName == "nvim-macos.tar.gz" {
		return "darwin"
	}
	if strings.Contains(fName, "macos") {
		return "darwin"
	}
	if strings.Contains(fName, "linux") {
		return "linux"
	}
	if strings.Contains(fName, "-win") {
		return "windows"
	}
	return
}

func (n *Neovim) archParser(fName string) (archStr string) {
	if fName == "nvim-macos.tar.gz" {
		return "any"
	}
	if strings.Contains(fName, "linux64") {
		return "amd64"
	}
	if strings.Contains(fName, "win64") {
		return "amd64"
	}
	if strings.Contains(fName, "x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "arm64") {
		return "amr64"
	}
	return
}

func (v *Neovim) vParser(tagName string) (vStr string) {
	return strings.TrimPrefix(tagName, "v")
}

func (n *Neovim) insParser(fName string) (insStr string) {
	return version.Unarchiver
}

func (n *Neovim) Start() {
	n.GhSearcher.Search(
		n.RepoName,
		n.tagFilter,
		n.fileFilter,
		n.vParser,
		n.archParser,
		n.osParser,
		n.insParser,
		nil,
	)
}

func (n *Neovim) GetVersions() []byte {
	r, _ := n.Version.Marshal()
	return r
}

func (n *Neovim) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"bin", "share"},
			MacOS:   []string{"bin", "share"},
			Linux:   []string{"bin", "share"},
		},
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}
