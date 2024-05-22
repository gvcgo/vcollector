package conda

/*
TODO: Conda Forge
https://raw.githubusercontent.com/conda-forge/feedstock-outputs/single-file/feedstock-outputs.json
*/
type CondaForgePackages struct {
	DownloadUrl string
	SDKName     string
}

func NewCondaForgePackages() (c *CondaForgePackages) {
	c = &CondaForgePackages{
		DownloadUrl: "https://raw.githubusercontent.com/conda-forge/feedstock-outputs/single-file/feedstock-outputs.json",
		SDKName:     "conda-forge-pkgs",
	}
	return
}
