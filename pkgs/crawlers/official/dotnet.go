package official

import (
	"encoding/json"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewDotnet())
}

func filterDotNetSDKByUrl(vUrl string) bool {
	if vUrl == "" {
		return false
	}
	excludeList := []string{
		"alpine",
		"installer",
		"winget",
		"arm32",
		"install",
		"scripts",
	}
	for _, ee := range excludeList {
		if strings.Contains(vUrl, ee) {
			return false
		}
	}
	return true
}

/*
https://dotnet.microsoft.com/en-us/download/dotnet
*/
type Dotnet struct {
	DownloadUrl string
	SDKName     string
	Version     version.VersionList
	host        string
}

func NewDotnet() (d *Dotnet) {
	d = &Dotnet{
		DownloadUrl: "https://dotnet.microsoft.com/download/dotnet",
		SDKName:     "dotnet",
		Version:     make(version.VersionList),
		host:        "https://dotnet.microsoft.com",
	}
	return
}

func (d *Dotnet) GetSDKName() string {
	return d.SDKName
}

func (d *Dotnet) parseArch(dUrl string) (arch string) {
	if strings.Contains(dUrl, "arm64") {
		return "arm64"
	}
	if strings.Contains(dUrl, "-x64") {
		return "amd64"
	}
	return
}

func (d *Dotnet) parseOS(dUrl string) (platform string) {
	if strings.Contains(dUrl, "linux") {
		return "linux"
	}
	if strings.Contains(dUrl, "-win-") {
		return "windows"
	}
	if strings.Contains(dUrl, "osx") {
		return "darwin"
	}
	return
}

func (d *Dotnet) fetchVersion(vUrl, vStr string) {
	doc := req.GetDocument(vUrl)
	if doc == nil {
		return
	}
	link := doc.Find("a#directLink").AttrOr("href", "")
	sha512Str := doc.Find("input#checksum").AttrOr("value", "")
	sumType := "sha512"
	if sha512Str == "" {
		sumType = sha512Str
	}

	item := version.Item{}
	item.Url = link
	item.SumType = sumType
	item.Sum = sha512Str
	item.Arch = d.parseArch(link)
	item.Os = d.parseOS(link)
	if item.Arch == "" || item.Os == "" {
		return
	}
	item.Installer = version.Unarchiver
	if _, ok := d.Version[vStr]; !ok {
		d.Version[vStr] = version.Version{}
	}
	d.Version[vStr] = append(d.Version[vStr], item)
}

func (d *Dotnet) getResult() {
	supportedVersionUrls := []string{}
	doc := req.GetDocument(d.DownloadUrl)
	if doc == nil {
		return
	}

	doc.Find("div#supported-versions-table").Find("table").Find("a").Each(func(_ int, s *goquery.Selection) {
		u := s.AttrOr("href", "")
		if u != "" && !strings.Contains(u, d.host) {
			u, _ = url.JoinPath(d.host, u)
		}
		supportedVersionUrls = append(supportedVersionUrls, u)
	})

	for _, u := range supportedVersionUrls {
		doc2 := req.GetDocument(u)
		if doc2 == nil {
			continue
		}
		doc2.Find("div.download-panel").Find("div").Find("table").Each(func(i int, ss *goquery.Selection) {
			vInfo := ss.Find("caption").AttrOr("id", "")
			vList := strings.Split(vInfo, "-sdk-")
			if len(vList) < 2 {
				return
			}
			vName := vList[len(vList)-1]
			ss.Find("a").Each(func(i int, sa *goquery.Selection) {
				uu := sa.AttrOr("href", "")
				if filterDotNetSDKByUrl(uu) {
					if !strings.Contains(uu, d.host) {
						uu, _ = url.JoinPath(d.host, uu)
					}
					d.fetchVersion(uu, vName)
				}
			})
		})
	}
}

func (d *Dotnet) Start() {
	d.getResult()
}

func (d *Dotnet) GetVersions() []byte {
	r, _ := d.Version.Marshal()
	return r
}

func TestDotnet() {
	dd := NewDotnet()
	dd.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/dotnet.json"
	content, _ := json.MarshalIndent(dd.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
