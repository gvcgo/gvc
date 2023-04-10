package confs

import "github.com/moqsien/gvc/pkgs/utils"

type Filemap map[string]string

type DavConf struct {
	DefaultWebdavConfigPath string             `koanf:"default_webdav_config_path"`
	DefaultWebdavLocalDir   string             `koanf:"default_webdav_local_dir"`
	DefaultWebdavRemoteDir  string             `koanf:"default_webdav_remote_dir"`
	DefaultWebdavHost       string             `koanf:"default_webdav_host"`
	DefaultFilesUrl         string             `koanf:"default_files_url"`
	FilesToSync             map[string]Filemap `koanf:"files_to_sync"`
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
	that.FilesToSync = map[string]Filemap{
		utils.Windows: {
			CodeUserSettingsBackupFileName: CodeUserSettingsFilePathForWin,
			CodeKeybindingsBackupFileName:  CodeKeybindingsFilePathForWin,
			NVimInitBackupFileName:         NVimWinInitPath,
			"chrome-bookmark":              "",
		},
		utils.Linux: {
			CodeUserSettingsBackupFileName: CodeUserSettingsFilePathForLinux,
			CodeKeybindingsBackupFileName:  CodeKeybindingsFilePathForLinux,
			NVimInitBackupFileName:         NVimUnixInitPath,
			"chrome-bookmark":              "",
		},
		utils.MacOS: {
			CodeUserSettingsBackupFileName: CodeUserSettingsFilePathForMac,
			CodeKeybindingsBackupFileName:  CodeKeybindingsFilePathForMac,
			NVimInitBackupFileName:         NVimUnixInitPath,
			"chrome-bookmark":              "",
		},
	}
}
