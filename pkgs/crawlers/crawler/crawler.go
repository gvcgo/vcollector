package crawler

type Crawler interface {
	Start()
	GetSDKName() string
}
