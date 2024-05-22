package mix

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/conda"
	"github.com/gvcgo/vcollector/pkgs/version"
)

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
		vStr := p.parseVersion(hh)
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

func (p *PHP) parseVersion(dUrl string) (vStr string) {
	vStr = PHPVersionRegexp.FindString(dUrl)
	vStr = strings.Trim(vStr, "-")
	return
}

func (p *PHP) GetVersions() version.VersionList {
	return p.Version
}

func (p *PHP) Start() {
	p.getFromConda()
	p.getFromOfficialReleases()
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
