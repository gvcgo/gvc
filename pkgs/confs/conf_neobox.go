package confs

import (
	"os"
	"path/filepath"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/koanfer"
	"github.com/moqsien/gvc/pkgs/utils"
	neoconf "github.com/moqsien/neobox/pkgs/conf"
)

type NBConf struct {
	Conf *neoconf.NeoConf `koanf:"neobox_conf"`
}

type NeoboxConf struct {
	NeoConfPath string `koanf:"neobox_conf_path"` // neobox config file path
	nconf       *NBConf
	path        string
}

func NewNeoboxConf() (r *NeoboxConf) {
	r = &NeoboxConf{
		path: ProxyFilesDir,
		nconf: &NBConf{
			Conf: &neoconf.NeoConf{},
		},
	}
	r.NeoConfPath = filepath.Join(r.path, "neobox_conf.json")
	r.setup()
	return
}

func (that *NeoboxConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *NeoboxConf) GetNeoConf() *neoconf.NeoConf {
	if ok, _ := utils.PathIsExist(that.NeoConfPath); ok {
		k, _ := koanfer.NewKoanfer(that.NeoConfPath)
		k.Load(that.nconf.Conf)
	}
	return that.nconf.Conf
}

func (that *NeoboxConf) Reset() {
	if that.nconf == nil {
		that.nconf = &NBConf{Conf: &neoconf.NeoConf{}}
	}
	that.nconf.Conf.WorkDir = that.path
	that.nconf.Conf.LogDir = filepath.Join(that.path, "neobox_logs")
	that.nconf.Conf.GeoInfoDir = filepath.Join(that.path, "neobox_geoinfo")
	that.nconf.Conf.SocketDir = filepath.Join(that.path, "neobox_socket")
	// TODO: Encrypted.
	that.nconf.Conf.DatabaseDir = GVCBackupDir // save sqlite db file to BackupDir
	// new added
	that.nconf.Conf.HistoryMaxLines = 300
	that.nconf.Conf.HistoryFileName = neoconf.HistoryFileName
	that.nconf.Conf.ShellSocketName = neoconf.ShellSocketName
	// TODO: change to github.
	that.nconf.Conf.DownloadUrl = "https://gitlab.com/moqsien/neobox_related/-/raw/main/conf.txt"

	// ping related
	that.nconf.Conf.MaxPingers = 120
	that.nconf.Conf.MaxPingAvgRTT = 600
	that.nconf.Conf.MaxPingPackLoss = 10

	// geoinfo files related
	that.nconf.Conf.GeoInfoSumUrl = "https://gitlab.com/moqsien/gvc_resources/-/raw/main/files_info.json?ref_type=heads&inline=false"
	// TODO: change to github.
	that.nconf.Conf.GeoInfoUrls = map[string]string{
		"geoip.dat":   "https://gitlab.com/moqsien/neobox_related/-/raw/main/geoip.dat",
		"geosite.dat": "https://gitlab.com/moqsien/neobox_related/-/raw/main/geosite.dat",
		"geoip.db":    "https://gitlab.com/moqsien/neobox_related/-/raw/main/geoip.db",
		"geosite.db":  "https://gitlab.com/moqsien/neobox_related/-/raw/main/geosite.db",
	}

	// verifier related
	that.nconf.Conf.InboundPort = 2023
	that.nconf.Conf.VerificationPortRange = &neoconf.PortRange{
		Min: 9035,
		Max: 9095,
	}
	that.nconf.Conf.VerificationTimeout = 3
	that.nconf.Conf.MaxToSaveRTT = 2000 // in milliseconds
	that.nconf.Conf.VerificationUrl = "https://www.google.com"
	that.nconf.Conf.VerificationCron = "@every 2h"

	// location related
	that.nconf.Conf.CountryAbbrevsUrl = "https://gitlab.com/moqsien/neobox_related/-/raw/main/country_names.json?ref_type=heads&inline=false"
	that.nconf.Conf.IPLocationQueryUrl = "https://www.fkcoder.com/ip?ip=%s"
	that.nconf.Conf.IPLocationQueryUrl2 = "http://ip-api.com/json/%s"

	// keeper related
	that.nconf.Conf.KeeperCron = "@every 3m"

	// cloudflare/wireguard related
	that.nconf.Conf.CloudflareConf = &neoconf.CloudflareConf{
		CloudflareIPV4URL: "https://www.cloudflare.com/ips-v4",
		PortList:          []int{443, 8443, 2053, 2096, 2087, 2083},
		MaxPingCount:      4,
		MaxGoroutines:     300,
		MaxRTT:            500,
		MaxLossRate:       0.0,
		MaxSaveToDB:       1000,
		WireGuardConfDir:  GVCBackupDir, // save wireguard configurations to BackupDir
		// TODO: change to github.
		CloudflareDomainFileUrl: "https://gitlab.com/moqsien/neobox_related/-/raw/main/cloudflare_domains.txt?ref_type=heads&inline=false",
	}
	// save to config file.
	k, _ := koanfer.NewKoanfer(that.NeoConfPath)
	k.Save(that.nconf.Conf)
}
