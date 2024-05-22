package conda

import "github.com/gvcgo/vcollector/internal/req"

/*
https://raw.githubusercontent.com/conda-forge/feedstock-outputs/single-file/feedstock-outputs.json
*/
type CondaForgePackages struct {
	DownloadUrl string
	SDKName     string
	result      []byte
}

func NewCondaForgePackages() (c *CondaForgePackages) {
	c = &CondaForgePackages{
		DownloadUrl: "https://raw.githubusercontent.com/conda-forge/feedstock-outputs/single-file/feedstock-outputs.json",
		SDKName:     "conda-forge-pkgs",
	}
	return
}

func (c *CondaForgePackages) GetSDKName() string {
	return c.SDKName
}

func (c *CondaForgePackages) GetVersions() []byte {
	return c.result
}

func (c *CondaForgePackages) Start() {
	resp := req.GetResp(c.DownloadUrl, 180)
	c.result = []byte(resp)
}
