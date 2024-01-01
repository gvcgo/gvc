package confs

import (
	"os"
	"path/filepath"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/goutils/pkgs/koanfer"
	"github.com/moqsien/gvc/pkgs/utils"
	neoconf "github.com/moqsien/neobox/pkgs/conf"
)

type NeoboxConf struct {
	NeoConfPath string           `koanf:"neobox_conf_path"` // neobox config file path
	NeoConf     *neoconf.NeoConf `koanf:"neobox_conf"`
	nconf       *neoconf.NeoConf
	path        string
}

func NewNeoboxConf() (r *NeoboxConf) {
	r = &NeoboxConf{
		NeoConfPath: filepath.Join(ProxyFilesDir, "neobox_conf.json"),
		path:        ProxyFilesDir,
		NeoConf:     &neoconf.NeoConf{},
		nconf:       &neoconf.NeoConf{},
	}
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
		k.Load(that.nconf)
	}
	return that.nconf
}

func (that *NeoboxConf) Reset() {
	if that.nconf == nil {
		that.nconf = &neoconf.NeoConf{}
	}
	that.nconf.WorkDir = that.path
	that.nconf.LogDir = filepath.Join(that.path, "neobox_logs")
	that.nconf.GeoInfoDir = filepath.Join(that.path, "neobox_geoinfo")
	that.nconf.SocketDir = filepath.Join(that.path, "neobox_socket")
	that.nconf.DatabaseDir = GVCBackupDir // save sqlite db file to BackupDir

	that.nconf.DownloadUrl = "https://gitlab.com/moqsien/neobox_related/-/raw/main/conf.txt"

	// ping related
	that.nconf.MaxPingers = 120
	that.nconf.MaxPingAvgRTT = 600
	that.nconf.MaxPingPackLoss = 10

	// geoinfo files related
	that.nconf.GeoInfoSumUrl = "https://gitlab.com/moqsien/gvc_resources/-/raw/main/files_info.json?ref_type=heads&inline=false"
	that.nconf.GeoInfoUrls = map[string]string{
		"geoip.dat":   "https://gitlab.com/moqsien/neobox_related/-/raw/main/geoip.dat",
		"geosite.dat": "https://gitlab.com/moqsien/neobox_related/-/raw/main/geosite.dat",
		"geoip.db":    "https://gitlab.com/moqsien/neobox_related/-/raw/main/geoip.db",
		"geosite.db":  "https://gitlab.com/moqsien/neobox_related/-/raw/main/geosite.db",
	}

	// verifier related
	that.nconf.InboundPort = 2023
	that.nconf.VerificationPortRange = &neoconf.PortRange{
		Min: 9035,
		Max: 9095,
	}
	that.nconf.VerificationTimeout = 3
	that.nconf.MaxToSaveRTT = 2000 // in milliseconds
	that.nconf.VerificationUrl = "https://www.google.com"
	that.nconf.VerificationCron = "@every 2h"

	// location related
	that.nconf.CountryAbbrevsUrl = "https://gitlab.com/moqsien/neobox_related/-/raw/main/country_names.json?ref_type=heads&inline=false"
	that.nconf.IPLocationQueryUrl = "https://www.fkcoder.com/ip?ip=%s"

	// keeper related
	that.nconf.KeeperCron = "@every 3m"

	// cloudflare/wireguard related
	that.nconf.CloudflareConf = &neoconf.CloudflareConf{
		CloudflareIPV4URL:       "https://www.cloudflare.com/ips-v4",
		PortList:                []int{443, 8443, 2053, 2096, 2087, 2083},
		MaxPingCount:            4,
		MaxGoroutines:           300,
		MaxRTT:                  500,
		MaxLossRate:             0.0,
		MaxSaveToDB:             1000,
		WireGuardConfDir:        GVCBackupDir, // save wireguard configurations to BackupDir
		CloudflareDomainFileUrl: "https://gitlab.com/moqsien/neobox_related/-/raw/main/cloudflare_domains.txt?ref_type=heads&inline=false",
	}
	// save to config file.
	k, _ := koanfer.NewKoanfer(that.NeoConfPath)
	k.Save(that.nconf)
}
