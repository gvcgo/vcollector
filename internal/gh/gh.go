package gh

type Asset struct {
	Name string `json:"name"`
	Url  string `json:"browser_download_url"`
}

type ReleaseItem struct {
	Assets     []Asset `json:"assets"`
	TagName    string  `json:"tag_name"`
	PreRelease any     `json:"prerelease"`
}

func GetReleaseItems(repoName string) (itemList []ReleaseItem) {
	uploader := NewUploader()
	return uploader.GetGithubReleaseList(repoName)
}
