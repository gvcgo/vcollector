package conda

import (
	"strings"

	"github.com/gvcgo/vcollector/internal/conda"
	"github.com/gvcgo/vcollector/pkgs/version"
)

type CondaSearcher struct {
	Version version.VersionList
}

func (c *CondaSearcher) Search(sdkName string) {
	versions := conda.SearchVersions(sdkName)

	for platform, versionList := range versions {
		pList := strings.Split(platform, "/")
		osStr := pList[0]
		archStr := pList[1]
		for _, vv := range versionList {
			if vv == "" {
				continue
			}
			item := version.Item{
				Os:        osStr,
				Arch:      archStr,
				Installer: version.Conda,
			}
			if _, ok := c.Version[vv]; !ok {
				c.Version[vv] = version.Version{}
			}
			c.Version[vv] = append(c.Version[vv], item)
		}
	}
}

func (c *CondaSearcher) GetVersions() version.VersionList {
	return c.Version
}
