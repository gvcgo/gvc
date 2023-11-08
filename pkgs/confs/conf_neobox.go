package confs

import (
	"os"
	"path/filepath"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
	neoconf "github.com/moqsien/neobox/pkgs/conf"
)

type NeoboxConf struct {
	NeoConf *neoconf.NeoConf `koanf:"neobox_conf"`
	path    string
}

func NewNeoboxConf() (r *NeoboxConf) {
	r = &NeoboxConf{
		path:    ProxyFilesDir,
		NeoConf: &neoconf.NeoConf{},
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

func (that *NeoboxConf) Reset() {
	if that.NeoConf == nil {
		that.NeoConf = &neoconf.NeoConf{}
	}
	that.NeoConf.WorkDir = that.path
	that.NeoConf.LogDir = filepath.Join(that.path, "neobox_logs")
	that.NeoConf.GeoInfoDir = filepath.Join(that.path, "neobox_geoinfo")
	that.NeoConf.SocketDir = filepath.Join(that.path, "neobox_socket")
	that.NeoConf.DatabaseDir = GVCBackupDir // save sqlite db file to BackupDir

	that.NeoConf.DownloadUrl = "https://gitlab.com/moqsien/neobox_related/-/raw/main/conf.txt"

	// ping related
	that.NeoConf.MaxPingers = 120
	that.NeoConf.MaxPingAvgRTT = 600
	that.NeoConf.MaxPingPackLoss = 10

	// geoinfo files related
	that.NeoConf.GeoInfoSumUrl = "https://gitlab.com/moqsien/gvc_resources/-/raw/main/files_info.json?ref_type=heads&inline=false"
	// TODO: reverse proxy
	that.NeoConf.GeoInfoUrls = map[string]string{
		"geoip.dat":   "https://gitlab.com/moqsien/neobox_related/-/raw/main/geoip.dat",
		"geosite.dat": "https://gitlab.com/moqsien/neobox_related/-/raw/main/geosite.dat",
		"geoip.db":    "https://gitlab.com/moqsien/neobox_related/-/raw/main/geoip.db",
		"geosite.db":  "https://gitlab.com/moqsien/neobox_related/-/raw/main/geosite.db",
	}

	// verifier related
	that.NeoConf.InboundPort = 2023
	that.NeoConf.VerificationPortRange = &neoconf.PortRange{
		Min: 9035,
		Max: 9095,
	}
	that.NeoConf.VerificationTimeout = 3
	that.NeoConf.MaxToSaveRTT = 2000 // in milliseconds
	that.NeoConf.VerificationUrl = "https://www.google.com"
	that.NeoConf.VerificationCron = "@every 2h"

	// location related
	that.NeoConf.CountryAbbrevsUrl = "https://gitlab.com/moqsien/neobox_related/-/raw/main/country_names.json?ref_type=heads&inline=false"
	that.NeoConf.IPLocationQueryUrl = "https://www.fkcoder.com/ip?ip=%s"

	// keeper related
	that.NeoConf.KeeperCron = "@every 3m"

	// cloudflare/wireguard related
	that.NeoConf.CloudflareConf = &neoconf.CloudflareConf{
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
}
