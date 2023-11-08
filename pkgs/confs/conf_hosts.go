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
		"https://raw.githubusercontent.com/JohyC/Hosts/main/MicrosoftHosts.txt",
		"https://raw.githubusercontent.com/JohyC/Hosts/main/EpicHosts.txt",
		"https://raw.githubusercontent.com/JohyC/Hosts/main/SteamDomains.txt",
		"https://raw.githubusercontent.com/JohyC/Hosts/main/hosts.txt",
		"https://raw.githubusercontent.com/ineo6/hosts/master/next-hosts",
		"https://raw.githubusercontent.com/sengshinlee/hosts/main/hosts",
	}
	that.HostFilters = []string{
		"github",
	}
	that.ReqTimeout = 30
	that.MaxAvgRtt = 400
	that.PingCount = 10
	that.WorkerNum = 100
}
