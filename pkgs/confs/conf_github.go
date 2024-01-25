package confs

type GithubConf struct {
	WinGitUrl string `koanf,json:"win_git_url"`
}

func NewGithubConf() (ghc *GithubConf) {
	ghc = &GithubConf{}
	return
}

func (that *GithubConf) Reset() {
	that.WinGitUrl = "https://github.com/git-for-windows/git/releases/latest/"
}
