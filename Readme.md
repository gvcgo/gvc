## [中文](https://github.com/moqsien/gvc/blob/main/docs/Readme_CN.md)
---------

### What's supported?
<figure class="third">
<img src="https://pkg.go.dev/static/shared/logo/go-white.svg" width="10%">
<img src="https://www.oracle.com/a/ocom/img/rc30v1-java-se.png" width="20%">
<img src="https://maven.apache.org/images/maven-logo-black-on-white.png" width="20%">
<img src="https://gradle.org/icon/favicon.ico" width="10%">
<img src="https://www.python.org/static/img/python-logo.png" width="25%">
<img src="https://nodejs.org/static/images/favicons/favicon.png" width="8%">
<img src="https://www.rust-lang.org/static/images/rust-logo-blk.svg" width="10%">
<img src="https://vlang.io/img/v-logo.png" width="10%">
<img src="https://www.cygwin.com/favicon.ico" width="10%">
<img src="https://storage.googleapis.com/cms-storage-bucket/ec64036b4eacc9f3fd73.svg" width="25%">
<img src="https://cn.julialang.org/assets/infra/logo_cn.png" width="18%">
<img src="https://code.visualstudio.com/favicon.ico" width="8%">
<img src="https://neovim.io/favicon.ico" width="8%">
<img src="https://brew.sh//assets/img/homebrew.svg" width="8%">
<img src="https://github.githubassets.com/favicons/favicon.svg" width="10%">
</figure>

---------
## Something nice about [gvc](https://github.com/moqsien/gvc).
---------
GVC is a nice tool designed for managing your development environment on multi-platforms and -machines.
It will help you to create a dev environment for Go, Python, Java, NodeJS or Rust, Cygwin, etc.
You do never need to worry about env variables or where to dowload files, GVC will do everything for you.
You can even install VSCode as well as Neovim through GVC.
Moreover, GVC will sync config files to WebDAV if you have configured one, which will help recreate your dev environment on another PC.

Therefore, If you are planning to get a new PC for development, probably the only thing you need to download GVC!

### Features
<table>
  <tbody>
  <tr>
    <th>Language/Tool</th>
    <th>Functions</th>
    <th>Note</th>
  </tr>
  <tr>
    <td><font color="Gree"> Go</font></td>
    <td><font color="LightBlue">Intall, Uninstall, SwitchVersion, SetEnv, SearchPackage</font></td>
    <td bgcolor="PaleVioletRed">gvc go help</td>
  </tr>
  <tr>
    <td><font color="Gree">Java</font></td>
    <td><font color="LightBlue">Install, Uninstall, SwitchVersion, SetEnv</font></td>
    <td bgcolor="LavenderBlush">gvc java help</td>
  </tr>
  <tr>
    <td><font color="Gree">Maven</font></td>
    <td><font color="LightBlue">Install, Uninstall, SwitchVersion, SetEnv, SetRepo</font></td>
    <td bgcolor="PaleVioletRed">gvc maven help</td>
  </tr>
  <tr>
    <td><font color="Gree">Gradle</font></td>
    <td><font color="LightBlue">Install, Uninstall, SwitchVersion, SetEnv, SetRepo</font></td>
    <td bgcolor="LavenderBlush">gvc gradle help</td>
  </tr>
  <tr>
    <td><font color="Gree">Python</font></td>
    <td><font color="LightBlue">Install, Uninstall, SwitchVersion, SetEnv, UpdatePyenv</font></td>
    <td bgcolor="PaleVioletRed">gvc py help</td>
  </tr>
  <tr>
    <td><font color="Gree">NodeJS</font></td>
    <td><font color="LightBlue">Install, Uninstall, SwitchVersion, SetEnv</font></td>
    <td bgcolor="LavenderBlush">gvc node help</td>
  </tr>
  <tr>
    <td><font color="Gree">Rust</font></td>
    <td><font color="LightBlue">Install, SetEnv</font></td>
    <td bgcolor="PaleVioletRed">gvc rust help</td>
  </tr>
  <tr>
    <td><font color="Gree">Vlang</font></td>
    <td><font color="LightBlue">Install, SetEnv</font></td>
    <td bgcolor="LavenderBlush">gvc vlang help</td>
  </tr>
  <tr>
    <td><font color="Gree">Cygwin</font></td>
    <td><font color="LightBlue">Install, InstallPackage</font></td>
    <td bgcolor="PaleVioletRed">gvc cygwin help; Only for Windows; git,bash, clang, gcc, etc.</td>
  </tr>
  <tr>
    <td><font color="Gree">Flutter</font></td>
    <td><font color="LightBlue">Install, Uninstall, SwitchVersion, SetEnv</font></td>
    <td bgcolor="LavenderBlush">gvc flutter help</td>
  </tr>
  <tr>
    <td><font color="Gree">Julia</font></td>
    <td><font color="LightBlue">Install, Uninstall, SwitchVersion, SetEnv</font></td>
    <td bgcolor="PaleVioletRed">gvc julia help</td>
  </tr>
  <tr>
    <td><font color="Gree">VSCode</font></td>
    <td><font color="LightBlue">Install, InstallExts,SetEnv, SyncSettingsToWebdav</font></td>
    <td bgcolor="LavenderBlush">gvc vscode help</td>
  </tr>
  <tr>
    <td><font color="Gree">NeoVim</font></td>
    <td><font color="LightBlue">Install, SyncInitFileToWebdav</font></td>
    <td bgcolor="PaleVioletRed">gvc nvim help</td>
  </tr>
  <tr>
    <td><font color="Gree">Homebrew</font></td>
    <td><font color="LightBlue">Install, SetEnv</font></td>
    <td bgcolor="LavenderBlush">gvc homebrew help</td>
  </tr>
  <tr>
    <td><font color="Gree">Hosts File</font></td>
    <td><font color="LightBlue">AutoModifyHostsFile</font></td>
    <td bgcolor="PaleVioletRed">gvc host help; Need root.</td>
  </tr>
  <tr>
    <td><font color="Gree">GVC Config</font></td>
    <td><font color="LightBlue">SetWebdavInfo, ResetDefaultConf, SyncSettingsToWebdav</font></td>
    <td bgcolor="LavenderBlush">gvc config help</td>
  </tr>
  <tr>
    <td><font color="Gree">Xray-Core</font></td>
    <td><font color="LightBlue">Auto get free Vmess start listen at localhost:2019</font></td>
    <td bgcolor="PaleVioletRed">HelpInfo: gvc xray help; EnterTheOperationShell: gvc xray</td>
  </tr>
  <tr>
    <td><font color="Gree">Github</font></td>
    <td><font color="LightBlue">Open github download acceleration website.</font></td>
    <td bgcolor="LavenderBlush">gvc github 1; gvc github 2</td>
  </tr>
</table>

Fetures on the way：
- Flutter auto-installation；

### Download & Install
Download files, unarchive, then double clik or just run with no subcommand or argument, gvc will install itself to default dir.

- [Note] Make sure you have PowerShell installed under windows.
- [github release](https://github.com/moqsien/gvc/releases)
- [gitee release](https://gitee.com/moqsien/gvc_tools/releases)

## gvc Help Info
---------
### gvc -h
```shell
moqsien@iMac ~ % gvc help
NAME:
   gvc - gvc <Command> <SubCommand>...

USAGE:
   gvc [global options] command [command options] [arguments...]

DESCRIPTION:
   A productive tool to manage your development environment.

COMMANDS:
   uninstall, unins, delete, del  [Caution] Delete gvc and softwares installed by gvc!
   show, sho, sh                  Show [gvc] install path.
   host, h, hosts                 Manage system hosts file.
   go, g                          Go version control.
   vscode, vsc, vs, v             VSCode management.
   config, conf, cnf, c           GVC config file management.
   nvim, neovim, nv, n            GVC neovim management.
   java, jdk, j                   GVC jdk management.
   gradle, gra, gr                Gradle management.
   maven, mav, ma                 Maven management.
   rust, rustc, ru, r             GVC rust management.
   nodejs, node, no               Nodejs version control.
   python, py                     Python version management.
   cygwin, cygw, cyg, cy          Cygwin management.
   xray, ray, xry, x              Start Xray Client.
   github, gh                     Github download acceleration websites.
   homebrew, brew, hb             Homebrew management.
   vlang, vl                      Vlang management.
   help, h                        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

### gvc go -h
```shell
moqsien@iMac-Pro gvc % gvc go help
NAME:
   gvc go - Go version control.

USAGE:
   gvc go command [command options] [arguments...]

COMMANDS:
   remote, r                   Show remote versions.
   use, u                      Download and use version.
   local, l                    Show installed versions.
   remove-unused, ru           Remove unused versions.
   remove-version, rm          Remove a version.
   add-envs, env, e, ae        Add envs for go.
   search-package, sp, search  Search for third-party packages.
   help, h                     Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc vs -h
```shell
moqsien@iMac-Pro gvc % gvc vs -h
NAME:
   gvc vscode - VSCode management.

USAGE:
   gvc vscode command [command options] [arguments...]

COMMANDS:
   install, i, ins               Automatically install vscode.
   install-extensions, ie, iext  Automatically install extensions for vscode.
   sync-extensions, se, sext     Push local installed vscode extensions info to remote webdav.
   get-settings, gs, gset        Get vscode settings(keybindings include) info from remote webdav.
   push-settings, ps, pset       Push vscode settings(keybindings include) info to remote webdav.
   help, h                       Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc nv -h
```shell
moqsien@iMac-Pro gvc % gvc nv -h
NAME:
   gvc nvim - GVC neovim management.

USAGE:
   gvc nvim command [command options] [arguments...]

COMMANDS:
   install, ins, i  Install neovim.
   help, h          Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc cnf -h
```shell
moqsien@iMac-Pro gvc % gvc cnf -h
NAME:
   gvc config - GVC config file management.

USAGE:
   gvc config command [command options] [arguments...]

COMMANDS:
   webdav, dav, w   Setup webdav account info to backup local settings for gvc, vscode, neovim etc.
   pull, pl         Pull settings to local backup dir from your remote webdav.
   push, ph         Push settings from local backup dir to your remote webdav.
   show, sh, s      Show path to conf files.
   reset, rs, r     Reset config file to default values.
   download, dl, d  Download example config files from gitee when backup dir is empty.
   help, h          Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc host -h
```shell
moqsien@iMac-Pro gvc % gvc host -h
NAME:
   gvc host - Manage system hosts file.

USAGE:
   gvc host command [command options] [arguments...]

COMMANDS:
   fetch, f      Fetch github hosts info.
   fetchall, fa  Get all github hosts info with no ping filters.
   show, s       Show hosts file path.
   help, h       Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc java -h
```shell
moqsien@iMac gvc % gvc java -h
NAME:
   gvc java - GVC jdk management.

USAGE:
   gvc java command [command options] [arguments...]

COMMANDS:
   use, u   Download and use jdk.
   show, s  Show available versions.
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc rust -h
```shell
moqsien@iMac gvc % gvc rust -h
NAME:
   gvc rust - GVC rust management.

USAGE:
   gvc rust command [command options] [arguments...]

COMMANDS:
   install, ins, i  Install the latest rust compiler tools.
   help, h          Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc node -h
```shell
moqsien@iMac gvc % gvc node -h
NAME:
   gvc nodejs - Nodejs version control.

USAGE:
   gvc nodejs command [command options] [arguments...]

COMMANDS:
   remote, r           Show remote versions.
   use, u              Download and use version.
   local, l            Show installed versions.
   remove-unused, ru   Remove unused versions.
   remove-version, rm  Remove a version.
   help, h             Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc py -h
```shell
moqsien@iMac gvc % gvc py -h  
NAME:
   gvc python - Python version management.

USAGE:
   gvc python command [command options] [arguments...]

COMMANDS:
   remote, r           Show remote versions.
   use, u              Download and use a version.
   local, l            Show installed versions.
   remove-version, rm  Remove a version.
   update, up          Install or update pyenv.
   help, h             Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc cy -h
```shell
moqsien@iMac gvc % gvc cy -h
NAME:
   gvc cygwin - Cygwin management.

USAGE:
   gvc cygwin command [command options] [arguments...]

COMMANDS:
   install, ins, i   Install Cygwin.
   package, pack, p  Install packages for Cygwin.
   help, h           Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc xray -h
```shell
moqsien@iMac gvc % gvc xray -h
NAME:
   gvc xray - Start Xray Client.

USAGE:
   gvc xray [command options] [arguments...]

OPTIONS:
   --start, --st, -s  Start Xray Client. (default: false)
   --help, -h         show help
```

### gvc brew -h
```shell
moqsien@iMac gvc % gvc brew -h
NAME:
   gvc homebrew - Homebrew management.

USAGE:
   gvc homebrew command [command options] [arguments...]

COMMANDS:
   install, ins, i     Install Homebrew.
   setenv, env, se, e  Set env to accelerate Homebrew in China.
   help, h             Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc vlang -h
```shell
moqsien@iMac gvc % gvc vlang -h
NAME:
   gvc vlang - Vlang management.

USAGE:
   gvc vlang command [command options] [arguments...]

COMMANDS:
   install, ins, i     Install Vlang.
   setenv, env, se, e  Set env for Vlang.
   help, h             Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc github -h
```shell
moqsien@iMac gvc % gvc github -h
NAME:
   gvc github - Github download acceleration websites.

USAGE:
   gvc github [command options] [arguments...]

OPTIONS:
   --help, -h  show help
```

### gvc maven -h
```shell
moqsien@iMac ~ % gvc maven -h
NAME:
   gvc maven - Maven management.

USAGE:
   gvc maven command [command options] [arguments...]

COMMANDS:
   use, u                  Download and use maven.
   show, s                 Show available versions.
   local, l                Show installed versions.
   set, se                 Set mirrors and local repository path.
   remove, rm              Remove an installed version.
   remove-unused, rmu, ru  Remove unused versions.
   help, h                 Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc gradle -h
```shell
moqsien@iMac ~ % gvc gradle -h
NAME:
   gvc gradle - Gradle management.

USAGE:
   gvc gradle command [command options] [arguments...]

COMMANDS:
   use, u                  Download and use gradle.
   show, s                 Show available versions.
   local, l                Show installed versions.
   set, se                 Set aliyun repository.
   remove, rm              Remove an installed version.
   remove-unused, rmu, ru  Remove unused versions.
   help, h                 Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

## thanks to
---------
- [xray-core](https://github.com/XTLS/Xray-core)
- [pyenv](https://github.com/pyenv/pyenv)
- [pyenv-win](https://github.com/pyenv-win/pyenv-win)
- [g](https://github.com/voidint/g)
- [gvm](https://github.com/andrewkroh/gvm)
