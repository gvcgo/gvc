package clis

func (that *Cli) SetVersionInfo(gitTag, gitHash, gitTime string) {
	that.gitHash = gitHash
	that.gitTag = gitTag
	that.gitTime = gitTime
}
