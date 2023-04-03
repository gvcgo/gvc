## [En](https://github.com/moqsien/gvc)
---------

## 关于[gvc](https://github.com/moqsien/gvc)的一些美好的事情
---------
GVC是一个全平台、多机器的一键管理多语言开发环境的辅助开发工具。
目前支持MacOS、Linux、Windows三大平台。
使用GVC能够轻松帮你一键搭建Go、Python、Java、Nodejs、Rust、Cygwin等开发环境，你可以轻松管理某个开发语言的多个版本，也不用自己操心任何环境变量。
此外，它还能轻松帮你一键搞定VSCode+Neovim安装和配置。
同时，GVC能把你的gvc配置，VSCode和Neovim配置同步到网盘，实现在其他机器上一键重建你熟悉的开发环境。你只需要配置一个任何支持WebDAV的网盘就行。
而且，GVC已经默认把很多加速方案进行了集成，比如Go的GOPROXY，Python的Pip以及本身安装包换成国内源，NPM添加国内源，Rust下载添加国内镜像等等。
重要的是，GVC是高度可配置的，你可以在gvc的主文件夹的backup目录下找到配置文件gvc-config.yml，然后修改比如加速镜像地址之类的，这样你就可以使用离你最近的镜像源，比如你在南方，可以使用中国科大或者浙大的镜像，你在北方可以使用清华镜像源等等。
除了Rust需要自己选择安装路径(由官方installer提供)之外，其他语言都默认安装在gvc的主目录中，当你不想要这些时，同样可以一键卸载所有，真是"强迫症"和"洁癖"患者的福音。

总之，GVC能帮助你搞定那些无聊的开发环境配置操作，当你想要尝试某个语言的新版本或者要在新的机器上做开发时，你无需再到处找下载资源，无需手动配置环境变量，你只需下载gvc即可。

# GVC 是一个跨平台多机器开发环境配置管理工具，让你轻松使用vscode进行多语言开发。
目前，gvc拥有以下功能或特点：
- Go编译器自动安装和添加环境变量，多版本轻松切换；
- Java JDK自动安装和添加环境变量，版本切换(jdk17 or jdk19)；
- Rust编译器自动安装和加速；
- Nodejs自动安装和添加环境变量，多版本轻松切换，安装包加速；
- Python自动安装，使用国内源解决下载慢问题，编译安装过程可能需要等待一段时间，同时自动配置环境变量和pip加速源；
- Cygwin自动安装和配置，包括了git, gcc, gfortran, clang, cmake, bash, wget等，解决Windows下c/c++开发以及git问题；
- VSCode自动安装，一键安装插件(需要配置，也可以使用默认配置)，一键备份和同步插件信息、用户设置、快捷键设置到webdav网盘；
- Neovim自动安装和配置，默认与vscode-neovim插件配合，有默认配置可以使用；
- Homebrew自动安装和加速；
- Vlang自动安装；
- Github下载加速网站使用默认浏览器打开；
- Hosts文件更新，加速github访问，对国内用户友好；
- 一键启动Xray代理，免费Vmess VPN，虽然速度不快，但事用于访问Google，Github等进行查资料没问题(本人不提供任何翻墙服务，也不通过这些收取服务费，请自行斟酌是否符合当地法律法规，依法使用)；
- 所有上述需要下载的地方，如果在国内较慢的，一般都有加速；
- 下载源可配置，如果你有更快的下载源，可以在gvc-config.yml中配置并注意保存；
- WebDAV网盘同步配置信息，可以一键将本地的包括gvc-config.yml在内的必要配置同步到网盘，在新机器上只需要使用这些配置就能重新搭建一样的开发环境；
- MacOS、Windows、Linux(暂未测试)全平台支持

gvc将要提供的功能或特点：
- Flutter自动安装；

### 下载和安装
下载文件，解压，双击或者在命令行运行(不带子任何命令和参数)，即可安装到默认文件夹。
- [注意] Windows下安装，请确保系统自带或者已安装PowerShell。
- [github release](https://github.com/moqsien/gvc/releases)
- [gitee release](https://gitee.com/moqsien/gvc/releases/tag/v2)

## gvc具体功能展示
---------
### gvc -h
```shell
moqsien@iMac gvc % gvc -h
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

```shell
moqsien@iMac gvc % gvc github -h
NAME:
   gvc github - Github download acceleration websites.

USAGE:
   gvc github [command options] [arguments...]

OPTIONS:
   --help, -h  show help
```