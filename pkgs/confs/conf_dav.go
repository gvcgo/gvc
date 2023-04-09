package confs

type DavConf struct {
	DefaultWebdavConfigPath string `koanf:"default_webdav_config_path"`
	DefaultWebdavLocalDir   string `koanf:"default_webdav_local_dir"`
	DefaultWebdavRemoteDir  string `koanf:"default_webdav_remote_dir"`
	DefaultWebdavHost       string `koanf:"default_webdav_host"`
	DefaultFilesUrl         string `koanf:"default_files_url"`
}

func NewDavConf() (r *DavConf) {
	r = &DavConf{
		DefaultWebdavConfigPath: GVCWebdavConfigPath,
		DefaultWebdavLocalDir:   GVCBackupDir,
		DefaultWebdavRemoteDir:  "/gvc_backups",
		DefaultWebdavHost:       "https://dav.jianguoyun.com/dav/",
		DefaultFilesUrl:         "https://gitee.com/moqsien/gvc/releases/download/v1/misc-all.zip",
	}
	return
}

func (that *DavConf) Reset() {
	that.DefaultWebdavConfigPath = GVCWebdavConfigPath
	that.DefaultWebdavLocalDir = GVCBackupDir
	that.DefaultWebdavRemoteDir = "/gvc_backups"
	that.DefaultWebdavHost = "https://dav.jianguoyun.com/dav/"
	that.DefaultFilesUrl = "https://gitee.com/moqsien/gvc/releases/download/v1/misc-all.zip"
}
