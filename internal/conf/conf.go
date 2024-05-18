package conf

const (
	ConfFileName string = "vcollector_conf.json"
)

type Config struct {
	Proxy       string
	GithubToken string
	GithubRepo  string
}
