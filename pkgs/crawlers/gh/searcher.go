package gh

import (
	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/pkgs/version"
)

/*
Search for github release items.
*/
type GithubSearcher struct {
	Version version.VersionList
}

type (
	ArchParser      func(fileName string) string
	OsParser        func(fileName string) string
	SumGetter       func(gh.Asset) (sum string, sumType string)
	InstallerParser func(fileName string) string
	Filter          func(gh.Asset) bool
)

func (s *GithubSearcher) Search(
	repoName string,
	filter Filter,
	archParser ArchParser,
	osParser OsParser,
	sumGetter SumGetter,
	installerParser InstallerParser,
) {
	itemList := gh.GetReleaseItems(repoName)
	for _, gItem := range itemList {

		vStr := gItem.TagName

	INNER:
		for _, a := range gItem.Assets {
			if filter != nil && !filter(a) {
				continue INNER
			}
			item := version.Item{}
			item.Arch = archParser(a.Name)
			item.Os = osParser(a.Name)
			if item.Arch == "" || item.Os == "" {
				continue INNER
			}
			item.Sum, item.SumType = sumGetter(a)
			item.Installer = installerParser(a.Name)
			if _, ok := s.Version[vStr]; !ok {
				s.Version[vStr] = version.Version{}
			}
			s.Version[vStr] = append(s.Version[vStr], item)
		}
	}
}

func (s *GithubSearcher) GetVersions() version.VersionList {
	return s.Version
}
