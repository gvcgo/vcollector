package mix

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/conda"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewPHP())
}

var PHPVersionRegexp = regexp.MustCompile(`-\d+(.\d+){2}-`)

/*
for windows:
https://windows.php.net/downloads/releases/archives/
https://github.com/pmmp/PHP-Binaries
*/
type PHP struct {
	SDKName     string
	RepoName    string
	DownloadUrl string
	host        string
	conda.CondaSearcher
}

func NewPHP() (p *PHP) {
	p = &PHP{
		SDKName:     "php",
		RepoName:    "pmmp/PHP-Binaries",
		DownloadUrl: "https://windows.php.net/downloads/releases/archives/",
		host:        "https://windows.php.net",
		CondaSearcher: conda.CondaSearcher{
			Version: make(version.VersionList),
		},
	}
	return
}

func (p *PHP) GetSDKName() string {
	return p.SDKName
}

func (p *PHP) getFromConda() {
	p.CondaSearcher.Search(p.SDKName)
}

func (p *PHP) getFromOfficialReleases() {
	doc := req.GetDocument(p.DownloadUrl)
	if doc == nil {
		return
	}
	doc.Find("a").Each(func(_ int, ss *goquery.Selection) {
		text := strings.TrimSpace(ss.Text())
		if !strings.HasSuffix(text, "zip") {
			return
		}
		if strings.Contains(text, "-nts-") {
			return
		}
		if !strings.Contains(text, "x64") {
			return
		}
		if strings.Contains(text, "-devel-") {
			return
		}
		if strings.Contains(text, "-test-") {
			return
		}
		if strings.Contains(text, "-debug-") {
			return
		}
		hh := ss.AttrOr("href", "")
		vStr := p.parseVersionForWin(hh)
		if vStr == "" || strings.HasPrefix(vStr, "5.") {
			return
		}
		if hh != "" {
			item := version.Item{}
			item.Url, _ = url.JoinPath(p.host, hh)
			item.Arch = "amd64"
			item.Os = "windows"
			item.Installer = version.Unarchiver
			if _, ok := p.Version[vStr]; !ok {
				p.Version[vStr] = version.Version{}
			}
			p.Version[vStr] = append(p.Version[vStr], item)
		}
	})
}

func (p *PHP) parseVersionForWin(dUrl string) (vStr string) {
	vStr = PHPVersionRegexp.FindString(dUrl)
	vStr = strings.Trim(vStr, "-")
	return
}

func (p *PHP) parseVersionForGithubItems(tagName string) (vStr string) {
	if !strings.HasPrefix(tagName, "php-") {
		return
	}
	vStr = strings.TrimPrefix(tagName, "php-")
	return
}

func (p *PHP) getFromGithub() {
	pItemList := gh.GetReleaseItems(p.RepoName)
	for _, pItem := range pItemList {
		vStr := p.parseVersionForGithubItems(pItem.TagName)
		if vStr == "" {
			continue
		}
	INNER:
		for _, aa := range pItem.Assets {
			if strings.Contains(aa.Name, "archive/refs/") {
				continue INNER
			}
			if strings.Contains(aa.Name, "-symbols") {
				continue INNER
			}
			if !strings.Contains(aa.Name, "-PM5") {
				continue INNER
			}
			item := version.Item{}
			item.Url = aa.Url
			item.Size = aa.Size
			item.Os = p.parseOsForGithubItem(aa.Name)
			item.Arch = p.parseArchForGithubItem(aa.Name)
			if item.Os == "" || item.Arch == "" {
				continue INNER
			}
			item.Installer = version.Unarchiver
			if _, ok := p.Version[vStr]; !ok {
				p.Version[vStr] = version.Version{}
			}
			p.Version[vStr] = append(p.Version[vStr], item)
		}
	}
}

func (p *PHP) parseOsForGithubItem(fName string) (platform string) {
	fName = strings.ToLower(fName)
	if strings.Contains(fName, "linux") {
		return "linux"
	}
	if strings.Contains(fName, "macos") {
		return "darwin"
	}
	return
}

func (p *PHP) parseArchForGithubItem(fName string) (arch string) {
	if strings.Contains(fName, "-x86_64") {
		return "amd64"
	}
	if strings.Contains(fName, "-arm64") {
		return "arm64"
	}
	if strings.Contains(fName, "-x64") {
		return "amd64"
	}
	return
}

func (p *PHP) GetVersions() []byte {
	r, _ := p.Version.Marshal()
	return r
}

func (p *PHP) Start() {
	p.getFromConda()
	p.getFromOfficialReleases()
	p.getFromGithub()
}

func (p *PHP) HomePage() string {
	return "https://www.php.net/"
}

func (p *PHP) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"php.exe"},
			MacOS:   []string{"bin", "lib"},
			Linux:   []string{"bin", "lib"},
		},
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestPHP() {
	nn := NewPHP()
	nn.Start()

	ff := fmt.Sprintf(
		"/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/%s.json",
		nn.SDKName,
	)
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
