package official

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewFlutter())
}

const (
	FlutterDownloadURLPattern   string = "https://storage.googleapis.com/flutter_infra_release/releases/releases_%s.json"
	FlutterDownloadURLPatternCN string = "https://storage.flutter-io.cn/flutter_infra_release/releases/releases_%s.json"
)

var FlutterOSMap = map[string]string{
	"linux":   "linux",
	"windows": "windows",
	"macos":   "darwin",
}

var FlutterArchMap = map[string]string{
	"x64":   "amd64",
	"arm64": "arm64",
}

type FItem struct {
	Version string `json:"version"`
	Stable  string `json:"channel"`
	Arch    string `json:"dart_sdk_arch"`
	Uri     string `json:"archive"`
	Sha256  string `json:"sha256"`
}

type FItemList struct {
	BaseUrl  string  `json:"base_url"`
	Releases []FItem `json:"releases"`
}

/*
https://storage.googleapis.com/flutter_infra_release/releases/releases_{linux/macos/windows}.json
https://storage.flutter-io.cn/flutter_infra_release/releases/releases_{linux/macos/windows}.json
*/
type Flutter struct {
	DownloadUrl string
	SDKName     string
	Version     version.VersionList
	result      FItemList
}

func NewFlutter() (f *Flutter) {
	f = &Flutter{
		// DownloadUrl: FlutterDownloadURLPatternCN,
		DownloadUrl: FlutterDownloadURLPattern,
		SDKName:     "flutter",
		Version:     make(version.VersionList),
		result:      FItemList{Releases: []FItem{}},
	}
	return
}

func (f *Flutter) GetSDKName() string {
	return f.SDKName
}

func (f *Flutter) getResult() {
	platforms := []string{"linux", "macos", "windows"}
	for _, platform := range platforms {
		dUrl := fmt.Sprintf(f.DownloadUrl, platform)
		req.GetJson(dUrl, &f.result)
	INNER:
		for _, fItem := range f.result.Releases {
			item := version.Item{}
			item.Arch = FlutterArchMap[fItem.Arch]
			if item.Arch == "" {
				continue INNER
			}
			item.Os = FlutterOSMap[platform]
			item.Sum = fItem.Sha256
			if item.Sum != "" {
				item.SumType = "sha256"
			}
			item.Url, _ = url.JoinPath(f.result.BaseUrl, fItem.Uri)
			item.Installer = version.Unarchiver
			vStr := fItem.Version
			if _, ok := f.Version[vStr]; !ok {
				f.Version[vStr] = version.Version{}
			}
			f.Version[vStr] = append(f.Version[vStr], item)
		}
	}
}

func (f *Flutter) Start() {
	f.getResult()
}

func (f *Flutter) GetVersions() []byte {
	r, _ := f.Version.Marshal()
	return r
}

func TestFlutter() {
	f := NewFlutter()
	f.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/flutter.json"
	content, _ := json.MarshalIndent(f.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
