package confs

type HostsConf struct {
	SourceUrls  []string `koanf:"source_urls"`
	HostFilters []string `koanf:"host_filters"`
	ReqTimeout  int      `koanf:"req_timeout"`
	MaxAvgRtt   int      `koanf:"max_avg_rtt"`
	PingCount   int      `koanf:"ping_count"`
	WorkerNum   int      `koanf:"worker_num"`
}

func NewHostsConf() *HostsConf {
	return &HostsConf{}
}

func (that *HostsConf) Reset() {
	that.SourceUrls = []string{
		"https://www.foul.trade:3000/Johy/Hosts/raw/branch/main/hosts.txt",
		"https://gitlab.com/ineo6/hosts/-/raw/master/next-hosts",
		"https://raw.hellogithub.com/hosts",
	}
	that.HostFilters = []string{
		"github",
	}
	that.ReqTimeout = 30
	that.MaxAvgRtt = 400
	that.PingCount = 10
	that.WorkerNum = 100
}
