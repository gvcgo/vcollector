package official

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gvcgo/vcollector/internal/req"
	"github.com/gvcgo/vcollector/pkgs/version"
)

const (
	KubectlDownloadUrlPattern  string = `https://dl.k8s.io/release/v%s/bin/%s/%s/kubectl`
	KubectlSha256UrlPattern    string = `https://dl.k8s.io/release/v%s/bin/%s/%s/kubectl.sha256`
	KubectlExeSha256UrlPattern string = `https://dl.k8s.io/v%s/bin/%s/%s/kubectl.exe.sha256`
)

var KubectlVersionRegexp = regexp.MustCompile(`\d+(.\d+){2}`)

/*
https://kubernetes.io/releases/patch-releases/
https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/
*/
type Kubectl struct {
	DownloadUrl string
	SDKName     string
	Version     version.VersionList
	latestUrl   string
}

func NewKubectl() (k *Kubectl) {
	k = &Kubectl{
		DownloadUrl: "https://kubernetes.io/releases/patch-releases/",
		SDKName:     "kubectl",
		Version:     make(version.VersionList),
		latestUrl:   "https://dl.k8s.io/release/stable.txt",
	}
	return
}

func (k *Kubectl) GetSDKName() string {
	return k.SDKName
}

func (k *Kubectl) parseVersion(vStr, archStr, osStr string) {
	sha256Url := fmt.Sprintf(KubectlSha256UrlPattern, vStr, osStr, archStr)
	if osStr == "windows" {
		sha256Url = fmt.Sprintf(KubectlExeSha256UrlPattern, vStr, osStr, archStr)
	}

	sha256 := req.GetResp(sha256Url)
	if strings.Contains(sha256, "NoSuchKey") {
		return
	}
	sha256 = strings.TrimSpace(sha256)

	dUrl := fmt.Sprintf(KubectlDownloadUrlPattern, vStr, osStr, archStr)
	if osStr == "windows" {
		dUrl += ".exe"
	}
	item := version.Item{}
	item.Arch = archStr
	item.Os = osStr
	item.Url = dUrl
	item.Sum = sha256
	item.SumType = "sha256"
	item.Installer = version.Unarchiver
	if _, ok := k.Version[vStr]; !ok {
		k.Version[vStr] = version.Version{}
	}
	k.Version[vStr] = append(k.Version[vStr], item)
}

func (k *Kubectl) getResult() {
	doc := req.GetDocument(k.DownloadUrl)
	if doc == nil {
		return
	}
	versionList := []string{}

	doc.Find("tr").Find("td").Each(func(_ int, s *goquery.Selection) {
		ss := KubectlVersionRegexp.FindString(strings.TrimSpace(s.Text()))
		if ss != "" && strings.Contains(ss, ".") {
			versionList = append(versionList, ss)
		}
	})

	s := req.GetResp(k.latestUrl)
	latestVersion := KubectlVersionRegexp.FindString(s)
	if latestVersion != "" {
		versionList = append(versionList, latestVersion)
	}
	archOSList := []string{
		"darwin/amd64",
		"darwin/arm64",
		"linux/amd64",
		"linux/arm64",
		"windows/amd64",
	}
	for _, archOS := range archOSList {
		aoList := strings.Split(archOS, "/")
		for _, vStr := range versionList {
			k.parseVersion(vStr, aoList[1], aoList[0])
		}
	}
}

func (k *Kubectl) Start() {
	k.getResult()
}

func (k *Kubectl) GetVersions() version.VersionList {
	return k.Version
}

func TestKubectl() {
	os.Setenv(req.ProxyEnvName, "http://127.0.0.1:2023")
	kk := NewKubectl()
	kk.Start()

	ff := "/Volumes/data/projects/go/src/gvcgo_org/vcollector/test/kubectl.json"
	content, _ := json.MarshalIndent(kk.Version, "", "    ")
	os.WriteFile(ff, content, os.ModePerm)
}
