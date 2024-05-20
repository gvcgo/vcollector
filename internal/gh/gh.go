package gh

type Assets struct {
	Name string `json:"name"`
	Url  string `json:"browser_download_url"`
}

type ReleaseItem struct {
	Assets     []*Assets `json:"assets"`
	TagName    string    `json:"tag_name"`
	PreRelease any       `json:"prerelease"`
}

type Filter func(dUrl string) bool

func GetReleaseItems(repoName string) (itemList []ReleaseItem) {
	uploader := NewUploader()
	return uploader.GetGithubReleaseList(repoName)
}
