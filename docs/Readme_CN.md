## [En](https://github.com/moqsien/gvc)
---------

<!-- <style>
   .cropped {
     width: 120px;
     height: 60px;
     overflow: hidden;
     border: 1px solid black;
   }
   .go {
      width: 80px; height: 60px;
      overflow: hidden;
      border: 1px solid black;
   }

   .s {
      width: 70px; height: 60px;
      overflow: hidden;
      border: 1px solid black;
   }

   .fl {
     width: 250px;
     height: 50px;
     overflow: hidden;
     border: 1px solid black;
   }
</style> -->
### gvc支持哪些语言或应用？
 <table>
 <tr>
<td>
<img align="left" src="https://golang.google.cn/images/favicon-gopher.png" class="go">
</td>

<td>
<img src="https://nodejs.org/static/images/favicons/favicon.png" class="cropped">
</td>

<td>
<img align="left" src="https://maven.apache.org/images/maven-logo-black-on-white.png" class="cropped">
</td>

<td>
<img src="https://gradle.org/icon/favicon.ico" class="s">
</td>

<td>
<img align="left" src="https://www.oracle.com/a/ocom/img/rc30v1-java-se.png" class="cropped">
</td>

<td>
<img src="https://www.python.org/static/img/python-logo.png" >
</td>
</tr>

<tr>
<td>
<img src="https://www.rust-lang.org/static/images/rust-logo-blk.svg" class="s">
</td>

<td>
<img src="https://vlang.io/img/v-logo.png" class="s">
</td>

<td>
<img src="https://www.cygwin.com/favicon.ico" class="s">
</td>

<td>
<img src="https://code.visualstudio.com/favicon.ico" class="go">
</td>

<td>
<img src="https://docs.flutter.dev/assets/images/shared/brand/flutter/logo/flutter-lockup.png" class="fl">
</td>

<td>
<img src="https://cn.julialang.org/assets/infra/logo_cn.png" class="fl">
</td>
</tr>
<tr>
<td>
<img src="https://brew.sh//assets/img/homebrew.svg" class="s">
</td>

<td>
<img src="https://neovim.io/favicon.ico" class="s">
</td>

<td>
<img src="https://github.githubassets.com/favicons/favicon.svg" class="s">
</td>
</tr>
</table>


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
    <td><font color="Gree">Typst</font></td>
    <td><font color="LightBlue">自动安装, 配置环境变量</font></td>
    <td bgcolor="LavenderBlush">gvc typst help</td>
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

### 下载和安装
下载文件，解压，双击或者在命令行运行(不带子任何命令和参数)，即可安装到默认文件夹。
- [注意] Windows下安装，请确保系统自带或者已安装PowerShell。
- [github release](https://github.com/moqsien/gvc/releases)
- [gitee release](https://gitee.com/moqsien/gvc_tools/releases)

## gvc具体功能展示
---------
### gvc -h
```shell
NAME:
   gvc - gvc <Command> <SubCommand>...

USAGE:
   gvc [global options] command [command options] [arguments...]

DESCRIPTION:
   A productive tool to manage your development environment.

COMMANDS:
   go, g                          Go version management.
   python, py                     Python version management.
   java, jdk, j                   Java jdk version management.
   maven, mav, ma                 Maven version management.
   gradle, gra, gr                Gradle version management.
   nodejs, node, no               NodeJS version management.
   flutter, flu, fl               Flutter version management.
   julia, jul, ju                 Julia version management.
   rust, rustc, ru, r             Rust installation.
   typst, ty                      Typst installation.
   vlang, vl                      Vlang installation.
   cygwin, cygw, cyg, cy          Cygwin installation.
   vscode, vsc, vs, v             VSCode and extensions installation.
   nvim, neovim, nv, n            Neovim installation.
   xray, ray, xry, x              Start Xray Shell for free VPN.
   homebrew, brew, hb             Homebrew installation or update.
   host, h, hosts                 Sytem hosts file management(need admistrator or root).
   github, gh                     Open github download acceleration websites.
   config, conf, cnf, c           Config file management for gvc.
   show, sho, sh                  Show [gvc] installation path and config file path.
   uninstall, unins, delete, del  [Caution] Delete gvc and softwares installed by gvc!
   help, h                        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

### 子命令帮助文档 (中文)
[github docs](https://github.com/moqsien/gvc/blob/main/docs/commands/command_list_github.md)

[gitee docs](https://gitee.com/moqsien/gvc_tools/blob/main/docs/commands/command_list_gitee.md)

## 感谢
---------
- [xray-core](https://github.com/XTLS/Xray-core)
- [pyenv](https://github.com/pyenv/pyenv)
- [pyenv-win](https://github.com/pyenv-win/pyenv-win)
- [g](https://github.com/voidint/g)
- [gvm](https://github.com/andrewkroh/gvm)
