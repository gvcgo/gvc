package vctrl

type Package struct {
	Url       string
	FileName  string
	Kind      string
	OS        string
	Arch      string
	Size      string
	Checksum  string
	CheckType string
}

type GoVersion struct {
	Version  string     // version
	Packages []*Package // packages info
}
