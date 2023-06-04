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
<img src="https://julialang.org/assets/infra/logo.svg" class="fl">
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
## 演示(以windows下的Powershell为例)

### 帮助信息
![](https://github.com/moqsien/gvc/blob/main/docs/ghelp.png)

### Go语言相关帮助信息
![](https://github.com/moqsien/gvc/blob/main/docs/goRemote.png)

### 显示目前可以获取的go编译器版本
![](https://github.com/moqsien/gvc/blob/main/docs/goRemoteShow.png)

### 自动安装go编译器并配置环境变量
![](https://github.com/moqsien/gvc/blob/main/docs/goInstall.png)
![](https://github.com/moqsien/gvc/blob/main/docs/goInstallationFinished.png)

### 通过关键字搜索go的第三方库
![](https://github.com/moqsien/gvc/blob/main/docs/goSearch.png)

### gvc的go build增强版
![beforeBuild](https://github.com/moqsien/gvc/blob/main/docs/beforebuild.png)
![chooseOptionsForBuild](https://github.com/moqsien/gvc/blob/main/docs/gobuild.png)
![chooseWethertoCompress](https://github.com/moqsien/gvc/blob/main/docs/compressOrNot.png)
![startBuild](https://github.com/moqsien/gvc/blob/main/docs/compiling.png)

### 其他语言的编译器或解释器的可获取的版本
![](https://github.com/moqsien/gvc/blob/main/docs/pyNodeFlutterJulia.png)

### 如何使用gvc中的NeoBox？
![](https://github.com/moqsien/gvc/blob/main/docs/neobox.png)


---------

![logo](https://github.com/moqsien/gvc/blob/main/docs/logo.png)
## 关于[gvc](https://github.com/moqsien/gvc)的一些美好的事情
---------
GVC是一个全平台、多机器的一键管理多语言开发环境的辅助开发工具。
目前支持MacOS、Linux、Windows三大平台。
使用GVC能够轻松帮你一键搭建Go、Python、Java、Nodejs、Flutter、Julia、Rust、C/C++等开发环境，你可以轻松管理某个开发语言的多个版本，也不用自己操心任何环境变量。

此外，它还能轻松帮你一键搞定VSCode+Neovim安装和配置。
同时，GVC能把你的gvc配置，VSCode和Neovim配置同步到网盘，实现在其他机器上一键重建你熟悉的开发环境。你只需要配置一个任何支持WebDAV的网盘就行。
而且，GVC已经默认把很多加速方案进行了集成，比如Go的GOPROXY，Python的Pip以及本身安装包换成国内源，NPM添加国内源，Rust下载添加国内镜像等等。

重要的是，GVC是高度可配置的，你可以在gvc的主文件夹的backup目录下找到配置文件gvc-config.json，然后修改比如加速镜像地址之类的，这样你就可以使用离你最近的镜像源，比如你在南方，可以使用中国科大或者浙大的镜像，你在北方可以使用清华镜像源等等。

GVC还能管理你的浏览器数据，很多常见的基于Chromium的浏览器以及Firefox浏览器的书签、插件信息、本地密码(可以设置加密保护)，都能一键导出，并同步到你自己的网盘。

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
    <td><font color="LightBlue">自动安装, 卸载, 版本切换, 配置环境变量, 关键词搜索第三方包, 更强大方便的go编译打包功能</font></td>
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
    <td bgcolor="PaleVioletRed">gvc cpp ic help; 仅用于Windows; git,bash, clang, gcc等将被默认安装.</td>
  </tr>
  <tr>
    <td><font color="Gree">Msys2</font></td>
    <td><font color="LightBlue">自动安装, 配置环境变量, 国内源加速</font></td>
    <td bgcolor="PaleVioletRed">gvc cpp im help; 仅用于Windows.</td>
  </tr>
  <tr>
    <td><font color="Gree">vcpkg</font></td>
    <td><font color="LightBlue">自动安装, 配置环境变量, 国内源加速</font></td>
    <td bgcolor="PaleVioletRed">gvc cpp iv help; C++包管理器，支持全平台.</td>
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
    <td bgcolor="PaleVioletRed">帮助信息: gvc x help; 进入xray操作shell: gvc x(可以在shell内控制启停等); <a href="https://github.com/moqsien/xtray">xtray docs</a></td>
  </tr>
  <tr>
    <td><font color="Gree">Browser</font></td>
    <td><font color="LightBlue">一键自动备份浏览器数据到webdav网盘，包括书签(html格式)、密码(json格式，可以自动加密内容之后上传，安全有保障)、插件信息；支持Chromium系列浏览器和Firefox浏览器等</font></td>
    <td bgcolor="PaleVioletRed">帮助信息: gvc browser help</td>
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
## gvc Help Info
---------
### gvc -h(use "g -h" for short)
```shell
NAME:
   g.exe - gvc <Command> <SubCommand>...

USAGE:
   g.exe [global options] command [command options] [arguments...]

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
   cpp                            C/C++ management.
   typst, ty                      Typst installation.
   vlang, vl                      Vlang installation.
   vscode, vsc, vs, v             VSCode and extensions installation.
   nvim, neovim, nv, n            Neovim installation.
   neobox-shell, shell, box, ns   Start a neobox shell.
   neobox-runner, nbrunner, nbr   Start a neobox client. # 此命令由NeoBox的Shell调用，用户无需关心。只需要在交互式Shell中操作即可。
   neobox-keeper, nbkeeper, nbk   Start a neobox keeper. # 此命令由NeoBox的Shell调用，用户无需关心。只需要在交互式Shell中操作即可。
   browser, br                    Browser data management.
   homebrew, brew, hb             Homebrew installation or update.
   hosts, h, host                 Sytem hosts file management(need admistrator or root).
   github, gh                     Open github download acceleration websites.
   config, conf, cnf, c           Config file management for gvc.
   version, ver, vsi              Show gvc version info.
   show, sho, sh                  Show [gvc] installation path and config file path.
   uninstall, unins, delete, del  [Caution] Remove gvc and softwares installed by gvc!
   help, h                        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

### NeoBox Shell帮助信息
### gvc ns
```bash
>>> help

Commands:
  add           Add proxies to neobox mannually.
  clear         clear the screen
  exit          exit the program
  export        Export vpn history list.
  filter        Filter vpns by verifier.
  gc            Start GC manually.
  geoinfo       Install/Update geoip&geosite for sing-box.
  help          display help
  parse         Parse raw proxy URIs to human readable ones.
  pingunix      Setup ping without root for Unix/Linux.
  restart       Restart the running sing-box client with a chosen vpn. [restart vpn_index]
  show          Show neobox info.
  start         Start an sing-box client/keeper.
  stop          Stop the running sing-box client/keeper.
```

### 子命令帮助文档 (中文)
[github docs](https://github.com/moqsien/gvc/blob/main/docs/commands/command_list_github.md)

[gitee docs](https://gitee.com/moqsien/gvc_tools/blob/main/docs/commands/command_list_gitee.md)

## 特别申明
本项目不提供任何收费服务，请任何使用者自觉遵守本国法律法规。

## 感谢
---------
- [xray-core](https://github.com/XTLS/Xray-core)
- [sing-box](https://github.com/SagerNet/sing-box)
- [pyenv](https://github.com/pyenv/pyenv)
- [pyenv-win](https://github.com/pyenv-win/pyenv-win)
- [g](https://github.com/voidint/g)
- [gvm](https://github.com/andrewkroh/gvm)
- [neobox](https://github.com/moqsien/neobox)
- [HackBrowserData](https://github.com/moonD4rk/HackBrowserData)
- [cygwin](https://github.com/cygwin/cygwin)
- [msys2](https://github.com/orgs/msys2/repositories)
- [vcpkg-tool](https://github.com/microsoft/vcpkg-tool)
- [gf](https://github.com/gogf/gf)
- [cli](https://github.com/urfave/cli)
- [pterm](https://github.com/pterm/pterm)
- [goutils](https://github.com/moqsien/goutils)

## 送我一杯咖啡~~~
[wechat](https://github.com/moqsien/moqsien/blob/main/imgs/wechat.jpeg)
[alipay](https://github.com/moqsien/moqsien/blob/main/imgs/alipay.jpeg)
