package confs

import (
	"path/filepath"

	"github.com/moqsien/gvc/pkgs/utils"
)

/*
gvc related
*/
var (
	GVCWorkDir          = filepath.Join(utils.GetHomeDir(), ".gvc/")
	GVCWebdavConfigPath = filepath.Join(GVCWorkDir, "webdav.yml")
	GVCBackupDir        = filepath.Join(GVCWorkDir, "backup")
	DefaultConfigPath   = filepath.Join(GVCWorkDir, "config.yml")
	RealConfigPath      = filepath.Join(GVCBackupDir, "gvc-config.yml")
)
