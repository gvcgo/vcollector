package official

import (
	"encoding/json"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/version"
)

var GolangOsMap = map[string]string{
	"macOS":   "darwin",
	"OS X":    "darwin",
	"Windows": "windows",
	"Linux":   "linux",
}

var GolangArchMap = map[string]string{
	"x86-64": "amd64",
	"ARM64":  "arm64",
}

/*
https://golang.google.cn/dl/
https://go.dev/dl/

https://dl.google.com/go/go1.22.3.linux-arm64.tar.gz
https://go.dev/dl/go1.22.3.linux-amd64.tar.gz
*/
type Golang struct {
	DownloadUrl string
	SDKName     string
	Version     version.VersionList
	host        string
}

func NewGolang() (g *Golang) {
	g = &Golang{
		DownloadUrl: "https://golang.google.cn/dl/",
		// DownloadUrl:"https://go.dev/dl/",
		SDKName: "go",
		Version: make(version.VersionList),
		host:    "https://go.dev",
	}
	return
}

func (g *Golang) GetSDKName() string {
	return g.SDKName
}

func (g *Golang) getResult() {
	doc := req.GetDocument(g.DownloadUrl)
	if doc == nil {
		return
	}
	doc.Find(".toggle").Each(func(_ int, s *goquery.Selection) {
		idStr := s.AttrOr("id", "")
		g.parseVersion(s, idStr)
	})
	doc.Find(".toggleVisible").Each(func(_ int, s *goquery.Selection) {
		idStr := s.AttrOr("id", "")
		g.parseVersion(s, idStr)
	})
}

func (g *Golang) parseArch(rawArch string) (arch string) {
	for key, val := range GolangArchMap {
		if strings.Contains(rawArch, key) {
			return val
		}
	}
	return
}

func (g *Golang) parseVersion(ss *goquery.Selection, vStr string) {
	vStr = strings.TrimSpace(vStr)
	if !strings.HasPrefix(vStr, "go") {
		return
	}
	vStr = strings.TrimPrefix(vStr, "go")
	if _, ok := g.Version[vStr]; !ok {
		g.Version[vStr] = version.Version{}
	}
	ss.Find("table.downloadtable").Find("tr").Each(func(_ int, s *goquery.Selection) {
		tds := s.Find("td")
		packageKind := strings.TrimSpace(tds.Eq(1).Text())
		arch := g.parseArch(tds.Eq(3).Text())
		osType := GolangOsMap[strings.TrimSpace(tds.Eq(2).Text())]

		if packageKind == "Archive" && arch != "" && osType != "" {
			dUrl := tds.Eq(0).Find("a").AttrOr("href", "")
			if dUrl == "" {
				return
			}

			item := version.Item{}
			item.Arch = arch
			item.Os = osType
			item.Installer = version.Unarchiver
			item.Extra = strings.TrimSpace(tds.Eq(4).Text())
			item.Sum = strings.TrimSpace(tds.Eq(5).Text())
			if len(item.Sum) == 64 {
				item.SumType = "sha256"
			} else if len(item.Sum) == 40 {
				item.SumType = "sha1"
			}
			item.Url, _ = url.JoinPath(g.host, dUrl)
			g.Version[vStr] = append(g.Version[vStr], item)
		}
	})
}

func (g *Golang) Start() {
	g.getResult()
}

func TestGolang() {
	gg := NewGolang()
	gg.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/go.json"
	content, _ := json.MarshalIndent(gg.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
