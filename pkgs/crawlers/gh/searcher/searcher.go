package searcher

import (
	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/pkgs/version"
)

type (
	ArchParser func(fileName string) string                                   // Arch
	OsParser   func(fileName string) string                                   // Os
	VParser    func(tagName string) string                                    // Version Name
	InsParser  func(fileName string) string                                   // Installer
	SumGetter  func(fileName string, assets []gh.Asset) (sum, sumType string) // Sum
	TagFilter  func(gh.ReleaseItem) bool
	FileFilter func(gh.Asset) bool
)

/*
Search release items for github repo.
*/
type GhSearcher struct {
	Version version.VersionList
}

func (g *GhSearcher) Search(
	repoName string,
	tagFilter TagFilter,
	fileFilter FileFilter,
	versionParser VParser,
	archParser ArchParser,
	osParser OsParser,
	insParser InsParser,
	sumGetter SumGetter,
) {
	gItemList := gh.GetReleaseItems(repoName)
	for _, gItem := range gItemList {
		if tagFilter != nil && !tagFilter(gItem) {
			continue
		}
		vStr := versionParser(gItem.TagName)
		if vStr == "" {
			continue
		}
		var (
			sumStr  string
			sumType string
		)
	INNER:
		for _, a := range gItem.Assets {
			if fileFilter != nil && !fileFilter(a) {
				continue INNER
			}
			if (sumStr == "" || sumType == "") && sumGetter != nil {
				sumStr, sumType = sumGetter(a.Name, gItem.Assets)
			}
			item := version.Item{}
			item.Arch = archParser(a.Name)
			item.Os = osParser(a.Name)
			if item.Arch == "" || item.Os == "" {
				continue INNER
			}
			item.Installer = insParser(a.Name)
			item.Url = a.Url
			item.Sum = sumStr
			item.SumType = sumType
			if _, ok := g.Version[vStr]; !ok {
				g.Version[vStr] = version.Version{}
			}
			g.Version[vStr] = append(g.Version[vStr], item)
		}
	}
}
