package official

/*
https://bell-sw.com/pages/downloads/
https://api.bell-sw.com/v1/liberica/releases/?&
*/
type LibericaItem struct {
	Bitness        int    `json:"bitness"`
	UpdateVersion  int    `json:"updateVersion"`
	DownloadUrl    string `json:"downloadUrl"`
	BundleType     string `json:"bundleType"`
	FeatureVersion int    `json:"featureVersion"`
	PackageType    string `json:"packageType"`
	Architecture   string `json:"architecture"`
	ExtraVersion   int    `json:"extraVersion"`
	BuildVersion   int    `json:"buildVersion"`
	Os             string `json:"os"`
	InterimVersion int    `json:"interimVersion"`
	Version        string `json:"version"`
	Sha1           string `json:"sha1"`
	Size           int64  `json:"size"`
	PatchVersion   int    `json:"patchVersion"`
}
