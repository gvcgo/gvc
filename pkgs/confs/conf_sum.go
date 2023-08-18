package confs

type SumConf struct {
	SumFileUrls string `koanf:"sum_file_urls"`
}

func NewSumConf() (r *SumConf) {
	r = &SumConf{}
	return
}

func (that *SumConf) Reset() {
	that.SumFileUrls = "https://gitlab.com/moqsien/gvc_resources/-/raw/main/files_info.json"
}
