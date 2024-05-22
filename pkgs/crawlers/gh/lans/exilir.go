package lans

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/pkgs/crawlers/gh/searcher"
	"github.com/gvcgo/vcollector/pkgs/version"
)

var (
	elixirOTPRegExp             = regexp.MustCompile(`otp-([0-9]+)`)
	elixirVersionPattern string = "%s-%s"
)

/*
https://github.com/elixir-lang/elixir
*/
type Elixir struct {
	SDKName  string
	RepoName string
	Version  version.VersionList
}

func NewElixir() (e *Elixir) {
	e = &Elixir{
		SDKName:  "elixir",
		RepoName: "elixir-lang/elixir",
		Version:  make(version.VersionList),
	}
	return
}

func (e *Elixir) GetSDKName() string {
	return e.SDKName
}

func (e *Elixir) filter(aa gh.Asset) bool {
	if strings.Contains(aa.Url, "archive/refs/") {
		return false
	}
	if strings.HasSuffix(aa.Name, ".sha1sum") {
		return false
	}
	if strings.HasSuffix(aa.Name, ".sha256sum") {
		return false
	}
	if !strings.HasPrefix(aa.Name, "elixir-otp-") {
		return false
	}
	return true
}

func (e *Elixir) parseVersion(eItem gh.ReleaseItem) {
	if searcher.GhVersionRegexp.FindString(eItem.TagName) == "" {
		return
	}

	for _, aa := range eItem.Assets {
		if !e.filter(aa) {
			continue
		}
		otpVersion := elixirOTPRegExp.FindString(aa.Name)
		if otpVersion == "" {
			continue
		}
		vStr := fmt.Sprintf(elixirVersionPattern, strings.TrimPrefix(eItem.TagName, "v"), otpVersion)
		if strings.HasSuffix(aa.Name, ".exe") {
			item := version.Item{}
			item.Url = aa.Url
			item.Arch = "amd64"
			item.Os = "windows"
			item.Size = aa.Size
			item.Installer = version.Executable
			if _, ok := e.Version[vStr]; !ok {
				e.Version[vStr] = version.Version{}
			}
			e.Version[vStr] = append(e.Version[vStr], item)
		} else if strings.HasSuffix(aa.Name, ".zip") {
			platforms := []string{"darwin", "linux"}
			for _, platform := range platforms {
				item := version.Item{}
				item.Url = aa.Url
				item.Arch = "any"
				item.Os = platform
				item.Size = aa.Size
				item.Installer = version.Unarchiver
				if _, ok := e.Version[vStr]; !ok {
					e.Version[vStr] = version.Version{}
				}
				e.Version[vStr] = append(e.Version[vStr], item)
			}
		}
	}
}

func (e *Elixir) getResult() {
	eItemList := gh.GetReleaseItems(e.RepoName)
	for _, eItem := range eItemList {
		e.parseVersion(eItem)
	}
}

func (e *Elixir) Start() {
	e.getResult()
}

func TestElixir() {
	nn := NewElixir()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
