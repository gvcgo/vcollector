package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gvcgo/goutils/pkgs/gutils"
	"github.com/gvcgo/vcollector/internal/conf"
	"github.com/gvcgo/vcollector/internal/gh"
	"github.com/gvcgo/vcollector/pkgs/crawlers/crawler"
)

const (
	ShaFileName            string = "sdklist_sha256.json"
	VersionFileNamePattern string = "%s.version.json"
)

// TODO: add installation config file sha256
type Sha256 struct {
	Sha256            string `json:"sha256"`
	HomePage          string `json:"homepage"`
	InstallConfSha256 string `json:"install_conf_sha256"`
}

type Sha256List map[string]Sha256

/*
1. Check sha256
2. Upload file to remote repo
3. Delete file from remote repo
*/
type Uploader struct {
	ShaFile      string
	VersionDir   string
	Github       *gh.Github
	Sha256List   Sha256List
	doNotSaveSha bool
}

func NewUploader() (u *Uploader) {
	u = &Uploader{
		ShaFile:    filepath.Join(conf.GetWorkDir(), ShaFileName),
		VersionDir: conf.GetVersionDir(),
		Github:     gh.NewGithub(),
		Sha256List: make(Sha256List),
	}
	u.loadSha256Info()
	return
}

func (u *Uploader) loadSha256Info() {
	if ok, _ := gutils.PathIsExist(u.ShaFile); ok {
		content, _ := os.ReadFile(u.ShaFile)
		json.Unmarshal(content, &u.Sha256List)
	}
}

func (u *Uploader) saveSha256Info() {
	content, _ := json.MarshalIndent(u.Sha256List, "", "  ")
	os.WriteFile(u.ShaFile, content, os.ModePerm)
}

func (u *Uploader) getVersionFilePath(sdkName string) string {
	fName := fmt.Sprintf(VersionFileNamePattern, sdkName)
	return filepath.Join(u.VersionDir, fName)
}

func (u *Uploader) saveVersionFile(sdkName string, content []byte) {
	os.WriteFile(u.getVersionFilePath(sdkName), content, os.ModePerm)
}

func (u *Uploader) checkSha256(sdkName, homepage string, content []byte) (ok bool) {
	h := sha256.New()
	h.Write(content)
	shaStr := fmt.Sprintf("%x", h.Sum(nil))

	if len(u.Sha256List) == 0 {
		u.loadSha256Info()
	}

	if ss, ok1 := u.Sha256List[sdkName]; !ok1 {
		if !u.doNotSaveSha && homepage != "" {
			u.Sha256List[sdkName] = Sha256{
				Sha256:   shaStr,
				HomePage: homepage,
			}
		}
		u.saveSha256Info()
		u.saveVersionFile(sdkName, content)
		return true
	} else {
		if ss.Sha256 == shaStr {
			return false
		} else {
			if !u.doNotSaveSha && homepage != "" {
				u.Sha256List[sdkName] = Sha256{
					Sha256:   shaStr,
					HomePage: homepage,
				}
			}
			u.saveSha256Info()
			u.saveVersionFile(sdkName, content)
			return true
		}
	}
}

func (u *Uploader) Upload(sdkName, homepage string, content []byte) {
	if len(content) == 0 {
		return
	}
	// "{}"
	if len(string(content)) < 10 {
		return
	}
	if u.checkSha256(sdkName, homepage, content) {
		localFilePath := u.getVersionFilePath(sdkName)
		remoteFilePath := filepath.Base(localFilePath)
		u.Github.UploadFile(remoteFilePath, localFilePath)
	}
}

func (u *Uploader) UploadSDKInfo(cc crawler.Crawler) {
	content := cc.GetVersions()
	homepage := cc.HomePage()
	sdkName := cc.GetSDKName()
	if len(string(content)) > 10 {
		if u.checkSha256(sdkName, homepage, content) {
			localFilePath := u.getVersionFilePath(sdkName)
			remoteFilePath := filepath.Base(localFilePath)
			u.Github.UploadFile(remoteFilePath, localFilePath)
		}
	}

	installConfContent, _ := json.Marshal(cc.GetInstallConf())
	installConfFile := filepath.Join(conf.GetInstallConfigFileDir(), fmt.Sprintf("%s.toml", sdkName))
	if ok, installConfSha := u.checkInstallConfFileSha256(installConfFile, installConfContent); ok && !u.doNotSaveSha {
		ss := u.Sha256List[sdkName]
		ss.InstallConfSha256 = installConfSha
		u.Sha256List[sdkName] = ss
		u.Github.UploadFile(fmt.Sprintf("install/%s.toml", sdkName), installConfFile)
	}
}

func (u *Uploader) checkInstallConfFileSha256(fPath string, content []byte) (ok bool, shaStr string) {
	h := sha256.New()
	h.Write(content)
	shaStr = fmt.Sprintf("%x", h.Sum(nil))

	content1, _ := os.ReadFile(fPath)
	h1 := sha256.New()
	h1.Write(content1)
	shaStr1 := fmt.Sprintf("%x", h1.Sum(nil))

	if shaStr == shaStr1 {
		ok = false
		return
	}
	ok = true
	os.WriteFile(fPath, content, os.ModePerm)
	return
}

func (u *Uploader) DisableSaveSha256() {
	u.doNotSaveSha = true
}
