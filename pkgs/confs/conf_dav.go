package confs

import "github.com/moqsien/gvc/pkgs/utils"

type Filemap map[string]string

type DavConf struct {
	DefaultWebdavConfigPath string             `koanf:"default_webdav_config_path"`
	DefaultWebdavLocalDir   string             `koanf:"default_webdav_local_dir"`
	DefaultWebdavRemoteDir  string             `koanf:"default_webdav_remote_dir"`
	DefaultWebdavHost       string             `koanf:"default_webdav_host"`
	FilesToSync             map[string]Filemap `koanf:"files_to_sync"`
}

func NewDavConf() (r *DavConf) {
	r = &DavConf{
		DefaultWebdavConfigPath: GVCWebdavConfigPath,
		DefaultWebdavLocalDir:   GVCBackupDir,
		DefaultWebdavRemoteDir:  "/gvc_backups",
		DefaultWebdavHost:       "https://dav.jianguoyun.com/dav/",
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
			CodeUserSettingsBackupFileName: `$appdata$Code\User\settings.json`,
			CodeKeybindingsBackupFileName:  `$appdata$Code\User\keybindings.json`,
			CodeExtensionsBackupFileName:   "",
			NVimInitBackupFileName:         `$home$\AppData\Local\nvim\init.vim`,
			"chrome-bookmark":              "",
		},
		utils.Linux: {
			CodeUserSettingsBackupFileName: "$home$.config/Code/User/settings.json",
			CodeKeybindingsBackupFileName:  "$home$.config/Code/User/keybindings.json",
			CodeExtensionsBackupFileName:   "",
			NVimInitBackupFileName:         "$home$.config/nvim/init.vim",
			"chrome-bookmark":              "",
		},
		utils.MacOS: {
			CodeUserSettingsBackupFileName: "$home$Library/Application Support/Code/User/settings.json",
			CodeKeybindingsBackupFileName:  "$home$Library/Application Support/Code/User/keybindings.json",
			CodeExtensionsBackupFileName:   "",
			NVimInitBackupFileName:         "$home$.config/nvim/init.vim",
			"chrome-bookmark":              "",
		},
	}
}
