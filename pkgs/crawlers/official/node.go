package official

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gvcgo/vcollector/internal/iconf"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
	"github.com/gvcgo/vcollector/pkgs/version"
)

func init() {
	crawler.RegisterCrawler(NewNode())
}

const (
	NodeDownloadUrl   string = "https://nodejs.org/download/release"
	NodeSumUrlPattern string = "https://nodejs.org/download/release/%s/SHASUMS256.txt"
)

type NodeItem struct {
	Version string `json:"version"`
	LTS     any    `json:"lts"`
	Date    string `json:"date"`
}

type NodeResult []NodeItem

/*
https://nodejs.org/dist/index.json
https://nodejs.org/download/release
*/
type Node struct {
	DownloadUrl string
	SDKName     string
	result      NodeResult
	Version     version.VersionList
}

func NewNode() (n *Node) {
	n = &Node{
		DownloadUrl: "https://nodejs.org/dist/index.json",
		SDKName:     "node",
		result:      NodeResult{},
		Version:     make(version.VersionList),
	}
	return
}

func (n *Node) GetSDKName() string {
	return n.SDKName
}

func (n *Node) filterYear(dateStr string) bool {
	dList := strings.Split(dateStr, "-")
	year, _ := strconv.Atoi(dList[0])
	return year > 2017
}

func (n *Node) parseArch(fName string) (arch string) {
	if strings.Contains(fName, "-arm64") {
		return "arm64"
	}
	if strings.Contains(fName, "-x64") {
		return "amd64"
	}
	return
}

func (n *Node) parseOS(fName string) (platform string) {
	if strings.Contains(fName, "-darwin") {
		return "darwin"
	}
	if strings.Contains(fName, "-linux") {
		return "linux"
	}
	if strings.Contains(fName, "-win") {
		return "windows"
	}
	return
}

func (n *Node) getVersion(vName string, isLTS bool) {
	sumUrl := fmt.Sprintf(NodeSumUrlPattern, vName)
	data := req.GetResp(sumUrl)
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "/") {
			continue
		}
		if strings.HasSuffix(line, ".tar.gz") || strings.HasSuffix(line, ".zip") {
			uList := strings.Split(line, " ")
			if len(uList) >= 2 {
				item := version.Item{}
				item.Sum = strings.TrimSpace(uList[0])
				item.SumType = "sha256"

				fName := strings.TrimSpace(uList[len(uList)-1])
				item.Url = fmt.Sprintf(
					"%s/%s/%s",
					NodeDownloadUrl,
					vName,
					fName,
				)
				item.Arch = n.parseArch(fName)
				item.Os = n.parseOS(fName)
				if item.Arch == "" || item.Os == "" {
					continue
				}
				if isLTS {
					item.LTS = "1"
				}
				item.Installer = version.Unarchiver
				vStr := strings.TrimPrefix(vName, "v")
				if _, ok := n.Version[vStr]; !ok {
					n.Version[vStr] = version.Version{}
				}
				n.Version[vStr] = append(n.Version[vStr], item)
			}
		}
	}
}

func (n *Node) getResult() {
	req.GetJson(n.DownloadUrl, &n.result)
	if len(n.result) == 0 {
		return
	}
	for _, nItem := range n.result {
		if n.filterYear(nItem.Date) {
			n.getVersion(nItem.Version, gconv.Bool(nItem.LTS))
		}
	}
}

func (n *Node) Start() {
	n.getResult()
}

func (n *Node) GetVersions() []byte {
	r, _ := n.Version.Marshal()
	return r
}

func (n *Node) HomePage() string {
	return "https://nodejs.org/en"
}

func (n *Node) GetInstallConf() (ic iconf.InstallerConfig) {
	return iconf.InstallerConfig{
		FlagFiles: &iconf.FileItems{
			Windows: []string{"LICENSE", "README.md"},
			MacOS:   []string{"LICENSE", "README.md"},
			Linux:   []string{"LICENSE", "README.md"},
		},
		BinaryDirs: &iconf.DirItems{
			Windows: []iconf.DirPath{},
			MacOS:   []iconf.DirPath{{"bin"}},
			Linux:   []iconf.DirPath{{"bin"}},
		},
	}
}

func TestNode() {
	nn := NewNode()
	nn.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/node.json"
	content, _ := json.MarshalIndent(nn.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
