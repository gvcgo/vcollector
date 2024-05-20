package gh

import (
	"encoding/json"
	"path/filepath"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gvcgo/goutils/pkgs/storage"
)

type Uploader struct {
	Github *storage.GhStorage
}

func NewUploader() (u *Uploader) {
	u = &Uploader{
		Github: storage.NewGhStorage(GithubUser, GithubToken),
	}
	u.initiate()
	return
}

func (u *Uploader) initiate() {
	u.Github.Proxy = GithubProxy
}

func (u *Uploader) Upload(localFilePath string) (r []byte) {
	fileName := filepath.Base(localFilePath)
	content := u.Github.GetContents(GithubRepo, "", fileName)
	shaStr := gjson.New(content).Get("sha").String()
	return u.Github.UploadFile(GithubRepo, "", localFilePath, shaStr)
}

func (u *Uploader) GetGithubReleaseList(repoName string) (rl []ReleaseItem) {
	r := u.Github.GetReleaseList(repoName)
	if len(r) == 0 {
		return
	}
	json.Unmarshal(r, &rl)
	return
}

/*
TODO: check sha256 for version files.
*/
