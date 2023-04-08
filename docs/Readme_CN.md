## [En](https://github.com/moqsien/gvc)
---------
### gvc支持哪些语言和工具?
<figure class="third">
<img src="https://golang.google.cn/images/favicon-gopher.png" width="10%"><img src="https://www.oracle.com/a/ocom/img/rc30v1-java-se.png" width="20%"><img src="https://maven.apache.org/images/maven-logo-black-on-white.png" width="20%"><img src="https://gradle.org/icon/favicon.ico" width="10%"><img src="https://www.python.org/static/img/python-logo.png" width="25%"><img src="https://nodejs.org/static/images/favicons/favicon.png" width="8%"><img src="https://www.rust-lang.org/static/images/rust-logo-blk.svg" width="10%"><img src="https://vlang.io/img/v-logo.png" width="10%"><img src="https://www.cygwin.com/favicon.ico" width="10%"><img src="https://storage.googleapis.com/cms-storage-bucket/ec64036b4eacc9f3fd73.svg" width="25%"><img src="https://cn.julialang.org/assets/infra/logo_cn.png" width="18%"><img src="https://code.visualstudio.com/favicon.ico" width="8%"><img src="https://neovim.io/favicon.ico" width="8%"><img src="https://brew.sh//assets/img/homebrew.svg" width="8%"><img src="https://github.githubassets.com/favicons/favicon.svg" width="10%">
</figure>

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
<table>
  <tbody>
  <tr>
    <th>语言/工具</th>
    <th>功能</th>
    <th>备注</th>
  </tr>
  <tr>
    <td><font color="Gree"> Go</font></td>
    <td><font color="LightBlue">自动安装, 卸载, 版本切换, 配置环境变量, 关键词搜索第三方包</font></td>
    <td bgcolor="PaleVioletRed">gvc go help</td>
  </tr>
  <tr>
    <td><font color="Gree">Java</font></td>
    <td><font color="LightBlue">自动安装, 卸载, 版本切换, 配置环境变量</font></td>
    <td bgcolor="LavenderBlush">gvc java help</td>
  </tr>
  <tr>
    <td><font color="Gree">Maven</font></td>
    <td><font color="LightBlue">自动安装, 卸载, 版本切换, 配置环境变量, 配置公有仓库国内镜像</font></td>
    <td bgcolor="PaleVioletRed">gvc maven help</td>
  </tr>
  <tr>
    <td><font color="Gree">Gradle</font></td>
    <td><font color="LightBlue">自动安装, 卸载, 版本切换, 配置环境变量, 配置公有仓库国内镜像</font></td>
    <td bgcolor="LavenderBlush">gvc gradle help</td>
  </tr>
  <tr>
    <td><font color="Gree">Python</font></td>
    <td><font color="LightBlue">自动安装, 卸载, 版本切换, 配置环境变量(含pip加速), 自动更新Pyenv</font></td>
    <td bgcolor="PaleVioletRed">gvc py help</td>
  </tr>
  <tr>
    <td><font color="Gree">NodeJS</font></td>
    <td><font color="LightBlue">自动安装, 卸载, 版本切换, 配置环境变量(含npm加速)</font></td>
    <td bgcolor="LavenderBlush">gvc node help</td>
  </tr>
  <tr>
    <td><font color="Gree">Rust</font></td>
    <td><font color="LightBlue">自动安装, 配置环境变量(国内加速)</font></td>
    <td bgcolor="PaleVioletRed">gvc rust help</td>
  </tr>
  <tr>
    <td><font color="Gree">Vlang</font></td>
    <td><font color="LightBlue">自动安装, 配置环境变量</font></td>
    <td bgcolor="LavenderBlush">gvc vlang help</td>
  </tr>
  <tr>
    <td><font color="Gree">Cygwin</font></td>
    <td><font color="LightBlue">自动安装, 配置环境变量, 国内源加速，自动添加Cygwin支持的软件工具</font></td>
    <td bgcolor="PaleVioletRed">gvc cygwin help; 仅用于Windows; git,bash, clang, gcc等将被默认安装.</td>
  </tr>
  <tr>
    <td><font color="Gree">Flutter</font></td>
    <td><font color="LightBlue">自动安装, 卸载, 切换版本, 配置环境变量(含国内源加速)</font></td>
    <td bgcolor="LavenderBlush">gvc flutter help</td>
  </tr>
  <tr>
    <td><font color="Gree">Julia</font></td>
    <td><font color="LightBlue">自动安装, 卸载, 切换版本, 配置环境变量(含国内源加速)</font></td>
    <td bgcolor="PaleVioletRed">gvc julia help</td>
  </tr>
  <tr>
    <td><font color="Gree">VSCode</font></td>
    <td><font color="LightBlue">自动安装, 自动安装插件(如果已配置),配置环境变量, VSCode相关配置同步到Webdav网盘(例如坚果云盘等)</font></td>
    <td bgcolor="LavenderBlush">gvc vscode help</td>
  </tr>
  <tr>
    <td><font color="Gree">NeoVim</font></td>
    <td><font color="LightBlue">自动安装, 环境变量配置，init配置文件同步到Webdav网盘</font></td>
    <td bgcolor="PaleVioletRed">gvc nvim help</td>
  </tr>
  <tr>
    <td><font color="Gree">Homebrew</font></td>
    <td><font color="LightBlue">自动安装, 配置环境变量(国内加速)</font></td>
    <td bgcolor="LavenderBlush">gvc homebrew help</td>
  </tr>
  <tr>
    <td><font color="Gree">Hosts File</font></td>
    <td><font color="LightBlue">自动修改系统Hosts文件，加速github和vscode插件市场访问</font></td>
    <td bgcolor="PaleVioletRed">gvc host help; 需要root或管理员权限.</td>
  </tr>
  <tr>
    <td><font color="Gree">GVC Config</font></td>
    <td><font color="LightBlue">配置Webdav信息, 恢复gvc默认配置, 同步配置文件到Webdav网盘</font></td>
    <td bgcolor="LavenderBlush">gvc config help; gvc自身的配置</td>
  </tr>
  <tr>
    <td><font color="Gree">Xray-Core</font></td>
    <td><font color="LightBlue">一键开启免费VPN——localhost:2019</font></td>
    <td bgcolor="PaleVioletRed">帮助信息: gvc xray help; 进入xray操作shell: gvc xray(可以在shell内控制启停等)</td>
  </tr>
  <tr>
    <td><font color="Gree">Github</font></td>
    <td><font color="LightBlue">使用默认浏览器打开github文件下载国内加速网站.</font></td>
    <td bgcolor="LavenderBlush">gvc github 1; gvc github 2</td>
  </tr>
</table>


gvc将要提供的功能或特点：
- Flutter自动安装；

### 下载和安装
下载文件，解压，双击或者在命令行运行(不带子任何命令和参数)，即可安装到默认文件夹。
- [注意] Windows下安装，请确保系统自带或者已安装PowerShell。
- [github release](https://github.com/moqsien/gvc/releases)
- [gitee release](https://gitee.com/moqsien/gvc_tools/releases)

## gvc具体功能展示
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
moqsien@iMac ~ % gvc java -h
NAME:
   gvc java - GVC jdk management.

USAGE:
   gvc java command [command options] [arguments...]

COMMANDS:
   use, u                  Download and use jdk.
   show, s                 Show available versions.
   local, l                Show installed versions.
   remove, rm              Remove an installed version.
   remove-unused, rmu, ru  Remove unused versions.
   help, h                 Shows a list of commands or help for one command

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

## 感谢
---------
- [xray-core](https://github.com/XTLS/Xray-core)
- [pyenv](https://github.com/pyenv/pyenv)
- [pyenv-win](https://github.com/pyenv-win/pyenv-win)
- [g](https://github.com/voidint/g)
- [gvm](https://github.com/andrewkroh/gvm)
