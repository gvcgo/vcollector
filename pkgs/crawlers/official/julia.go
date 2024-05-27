package official

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewJulia())
}

var JuliaArchMap = map[string]string{
	"aarch64": "arm64",
	"x86_64":  "amd64",
}

var JuliaOSMap = map[string]string{
	"linux": "linux",
	"mac":   "darwin",
	"winnt": "windows",
}

type JuliaVersion struct {
	Url       string `json:"url"`
	Kind      string `json:"kind"`
	Arch      string `json:"arch"`
	Sum       string `json:"sha256"`
	Os        string `json:"os"`
	Extension string `json:"extension"`
}

type JuliaItem struct {
	Files  []JuliaVersion `json:"files"`
	Stable any            `json:"stable"`
}

type JuliaItemList map[string]JuliaItem

/*
https://julialang-s3.julialang.org/bin/versions.json
https://mirrors.tuna.tsinghua.edu.cn/julia-releases/bin/versions.json
*/
type Julia struct {
	DownloadUrl string
	SDKName     string
	Version     version.VersionList
	result      JuliaItemList
}

func NewJulia() (j *Julia) {
	j = &Julia{
		// DownloadUrl: "https://mirrors.tuna.tsinghua.edu.cn/julia-releases/bin/versions.json",
		DownloadUrl: "https://julialang-s3.julialang.org/bin/versions.json",
		SDKName:     "julia",
		Version:     make(version.VersionList),
		result:      make(JuliaItemList),
	}
	return
}

func (j *Julia) GetSDKName() string {
	return j.SDKName
}

func (j *Julia) parseVersion(jItem JuliaVersion, vStr string) {
	item := version.Item{}
	item.Url = jItem.Url
	item.Sum = jItem.Sum
	if item.Sum != "" {
		item.SumType = "sha256"
	}
	item.Arch = JuliaArchMap[jItem.Arch]
	item.Os = JuliaOSMap[jItem.Os]
	if item.Arch == "" || item.Os == "" {
		return
	}
	item.Installer = version.Unarchiver
	if _, ok := j.Version[vStr]; !ok {
		j.Version[vStr] = version.Version{}
	}
	j.Version[vStr] = append(j.Version[vStr], item)
}

func (j *Julia) getResult() {
	req.GetJson(j.DownloadUrl, &j.result)
	for vStr, jItemList := range j.result {
		for _, jItem := range jItemList.Files {
			if jItem.Kind != "archive" {
				continue
			}
			if !strings.Contains(jItem.Extension, "tar.gz") {
				continue
			}
			j.parseVersion(jItem, vStr)
		}
	}
}

func (j *Julia) Start() {
	j.getResult()
}

func (j *Julia) GetVersions() []byte {
	r, _ := j.Version.Marshal()
	return r
}

func (j *Julia) HomePage() string {
	return "https://julialang.org/"
}

func (j *Julia) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"LICENSE.md"},
			MacOS:   []string{"LICENSE.md"},
			Linux:   []string{"LICENSE.md"},
		},
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"bin"}},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestJulia() {
	jj := NewJulia()
	jj.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/julia.json"
	content, _ := json.MarshalIndent(jj.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
