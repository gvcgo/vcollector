package iconf

type FileItems struct {
	Windows []string `toml:"windows"`
	Linux   []string `toml:"linux"`
	MacOS   []string `toml:"darwin"`
}

type (
	DirPath  []string
	DirItems struct {
		Windows []DirPath `toml:"windows"` // <symbolLinkPath>/<filepath.Join(List)>, ...
		Linux   []DirPath `toml:"linux"`
		MacOS   []DirPath `toml:"darwin"`
	}
)

type AdditionalEnv struct {
	Name    string
	Value   []DirPath // <symbolLinkPath>/<filepath.Join(Value)>
	Version string    // major>8 or major<=8(for JDK)
}

type AdditionalEnvList []AdditionalEnv

type BinaryRename struct {
	NameFlag string `toml:"name_flag"`
	RenameTo string `toml:"rename_to"`
}

/*
Installation configs
*/
type InstallerConfig struct {
	FlagFiles       *FileItems        `toml:"flag_files"`
	FlagDirExcepted bool              `toml:"flag_dir_excepted"`
	BinaryDirs      *DirItems         `toml:"binary_dirs"`
	BinaryRename    *BinaryRename     `toml:"binary_rename"`
	AdditionalEnvs  AdditionalEnvList `toml:"additional_envs"`
}
