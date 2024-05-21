package gh

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gvcgo/goutils/pkgs/request"
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
	User    string
	Repo    string
	Token   string
	Proxy   string
	fetcher *request.Fetcher
}

func NewGithub() (g *Github) {
	g = &Github{
		User:    GithubUser,
		Repo:    GithubRepo,
		Token:   GithubToken,
		Proxy:   GithubProxy,
		fetcher: request.NewFetcher(),
	}
	g.initiate()
	return
}

func (g *Github) initiate() {
	if g.Proxy != "" {
		g.fetcher.Proxy = g.Proxy
	}
	g.fetcher.Headers = map[string]string{
		"Accept":        AcceptHeader,
		"Authorization": fmt.Sprintf(AuthorizationHeader, g.Token),
	}
}

func (g *Github) GetContents(repoName, remotePath string) (r []byte) {
	// TODO: get shaStr
	return
}

func (g *Github) UploadFile(repoName, remotePath, localPath, shaStr string) (r []byte) {
	// TODO: upload file
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
	rl := GetReleaseItems("exaloop/codon")
	fmt.Println(len(rl))
	for _, item := range rl {
		fmt.Println(item.TagName)
	}
}
