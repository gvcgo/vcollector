package gh

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gvcgo/goutils/pkgs/gtea/gprint"
	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/goutils/pkgs/request"
	"github.com/gvcgo/vcollector/internal/conf"
	"github.com/gvcgo/vcollector/internal/req"
)

// ReleaseItem
type Asset struct {
	Name string `json:"name"`
	Url  string `json:"browser_download_url"`
	Size int64  `json:"size"`
}

type ReleaseItem struct {
	Assets     []Asset `json:"assets"`
	TagName    string  `json:"tag_name"`
	PreRelease any     `json:"prerelease"`
}

type ReleaseList []ReleaseItem

const (
	GithubAPI           string = "https://api.github.com"
	AcceptHeader        string = "application/vnd.github.v3+json"
	AuthorizationHeader string = "token %s"
)

type Github struct {
	Repo    string
	Token   string
	Proxy   string
	fetcher *request.Fetcher
}

func NewGithub() (g *Github) {
	cnf := conf.NewConfig()
	repo := GithubRepo
	if cnf.GithubRepo != "" {
		repo = cnf.GithubRepo
	}
	token := GithubToken
	if cnf.GithubToken != "" {
		token = cnf.GithubToken
	}

	proxy := GithubProxy
	if cnf.Proxy != "" {
		proxy = cnf.Proxy
	}
	g = &Github{
		Repo:    repo,
		Token:   token,
		Proxy:   proxy,
		fetcher: request.NewFetcher(),
	}
	g.initiate()
	return
}

func (g *Github) initiate() {
	if g.Proxy != "" && gconv.Bool(os.Getenv(req.ProxyEnvName)) {
		g.fetcher.Proxy = g.Proxy
	}
	g.fetcher.Headers = map[string]string{
		"Accept":        AcceptHeader,
		"Authorization": fmt.Sprintf(AuthorizationHeader, g.Token),
	}
}

func (g *Github) GetShaStr(remotePath string) (shaStr string) {
	// https://api.github.com/repos/{user}/{repo}/contents/{remotePath}
	dUrl := fmt.Sprintf("%s/repos/%s/contents/%s", GithubAPI, g.Repo, remotePath)
	g.fetcher.Timeout = 30 * time.Second
	g.fetcher.SetUrl(dUrl)
	resp, _ := g.fetcher.GetString()
	shaStr = gjson.New(resp).Get("sha").String()
	return
}

func (g *Github) UploadFile(remotePath, localPath string) (r []byte) {
	// https://api.github.com/repos/{user}/{repo}/contents/{path}/{filename}
	if ok, _ := gutils.PathIsExist(localPath); !ok {
		gprint.PrintError("file: %s does not exist.", localPath)
		return
	}

	g.fetcher.SetUrl(fmt.Sprintf("%s/repos/%s/contents/%s", GithubAPI, g.Repo, remotePath))
	g.fetcher.Timeout = 5 * time.Minute

	content, _ := os.ReadFile(localPath)
	shaStr := g.GetShaStr(remotePath)
	g.fetcher.PostBody = map[string]interface{}{
		"message": fmt.Sprintf("update file: %s.", remotePath),
		"content": base64.StdEncoding.EncodeToString(content),
		"sha":     shaStr,
	}
	if resp := g.fetcher.Put(); resp != nil {
		defer resp.RawResponse.Body.Close()
		r, _ = io.ReadAll(resp.RawResponse.Body)
	}
	return
}

func (g *Github) getRelease(repoName string, page int) (r []byte) {
	// https://api.github.com/repos/{owner}/{repo}/releases?per_page=100&page=1
	dUrl := fmt.Sprintf("%s/repos/%s/releases?per_page=100&page=%d", GithubAPI, repoName, page)
	g.fetcher.SetUrl(dUrl)
	g.fetcher.Timeout = 180 * time.Second
	if resp := g.fetcher.Get(); resp != nil {
		defer resp.RawResponse.Body.Close()
		r, _ = io.ReadAll(resp.RawResponse.Body)
	}
	return
}

func (g *Github) GetReleases(repoName string) (rl ReleaseList) {
	page := 1
	for {
		itemList := ReleaseList{}
		r := g.getRelease(repoName, page)
		json.Unmarshal(r, &itemList)
		if len(itemList) == 0 || page >= 10 {
			break
		}
		rl = append(rl, itemList...)
		page++
	}
	return
}

func GetReleaseItems(repoName string) ReleaseList {
	gh := NewGithub()
	return gh.GetReleases(repoName)
}

func TestGithub() {
	rl := GetReleaseItems("pmmp/PHP-Binaries")
	fmt.Println(len(rl))
	for _, item := range rl {
		fmt.Println(item.TagName)
	}
}
