package fixed

import (
	"encoding/json"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewMiniconda())
}

/*
https://repo.anaconda.com/miniconda/
*/
type Miniconda struct {
	DownloadUrl string
	SDKName     string
	Version     version.VersionList
}

func NewMiniconda() (m *Miniconda) {
	m = &Miniconda{
		DownloadUrl: "https://repo.anaconda.com/miniconda/",
		SDKName:     "miniconda",
		Version:     make(version.VersionList),
	}
	return
}

func (m *Miniconda) GetSDKName() string {
	return m.SDKName
}

func (m *Miniconda) filterMinicondaByFName(fname string) bool {
	r := false
	toBunList := []string{
		".pkg",
		"Miniconda2-latest-", // for miniconda
		"Miniconda-latest-",  // for miniconda
	}
	for _, b := range toBunList {
		if strings.Contains(fname, b) {
			return true
		}
	}
	return r
}

func (m *Miniconda) parseOS(fname string) (platform string) {
	if strings.Contains(fname, "Windows") {
		return "windows"
	}
	if strings.Contains(fname, "Linux") {
		return "linux"
	}
	if strings.Contains(fname, "MacOSX") {
		return "darwin"
	}
	return
}

func (m *Miniconda) parseArch(fname string) (archStr string) {
	if strings.Contains(fname, "x86_64") {
		return "amd64"
	}
	if strings.Contains(fname, "arm64") {
		return "arm64"
	}
	if strings.Contains(fname, "aarch64") {
		return "arm64"
	}
	return
}

func (m *Miniconda) getResult() {
	doc := req.GetDocument(m.DownloadUrl)
	if doc == nil {
		return
	}

	doc.Find("table").Find("tr").Each(func(ii int, s *goquery.Selection) {
		u := s.Find("td").Eq(0).Find("a").AttrOr("href", "")
		if u == "" {
			return
		}
		fName := s.Find("td").Eq(0).Find("a").Text()
		if m.filterMinicondaByFName(fName) {
			return
		}
		sha256Str := s.Find("td").Eq(3).Text()

		if strings.Contains(fName, "latest") {
			item := version.Item{}
			if !strings.HasPrefix(u, "http") {
				u, _ = url.JoinPath(m.DownloadUrl, u)
			}
			item.Url = u
			item.Arch = m.parseArch(fName)
			item.Os = m.parseOS(fName)
			if item.Os == "" || item.Arch == "" {
				return
			}

			item.Sum = sha256Str
			if item.Sum != "" {
				item.SumType = "sha256"
			}
			item.Installer = version.Executable
			vStr := "latest"
			// fmt.Printf("%+v\n", item)
			if _, ok := m.Version[vStr]; !ok {
				m.Version[vStr] = version.Version{}
			}
			m.Version[vStr] = append(m.Version[vStr], item)
		}
	})
}

func (m *Miniconda) Start() {
	m.getResult()
}

func (m *Miniconda) GetVersions() []byte {
	r, _ := m.Version.Marshal()
	return r
}

func (m *Miniconda) HomePage() string {
	return "https://docs.anaconda.com/free/miniconda/index.html"
}

func (m *Miniconda) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{{"bin"}, {"condabin"}},
			MacOS:   []iconf.DirPath{{"bin"}, {"condabin"}},
			Linux:   []iconf.DirPath{{"bin"}, {"condabin"}},
		},
	}
}

func TestMiniconda() {
	mm := NewMiniconda()
	mm.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/miniconda.json"
	content, _ := json.MarshalIndent(mm.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
