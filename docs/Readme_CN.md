## [En](https://github.com/moqsien/gvc)
---------

### gvc支持哪些语言或应用？

<table>
  <thead>
    <tr>
      <th>语言/应用</th>
      <th>备注</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><a href="https://go.dev/">go官网</a></td>
      <td><a href="https://golang.google.cn/">go国内下载</a>|<a href="https://mirrors.aliyun.com/golang/">阿里镜像源</a></td>
    </tr>
    <tr>
      <td><a href="https://www.oracle.com/java/technologies/downloads/">java/jdk</a></td>
      <td><a href="https://www.injdk.cn/">java/jdk国内下载</a></td>
    </tr>
    <tr>
      <td><a href="https://maven.apache.org/download.cgi">maven官网</a></td>
      <td><a href="https://dlcdn.apache.org/maven/">maven下载</a></td>
    </tr>
    <tr>
      <td><a href="https://gradle.org/install/">gradle官网</a></td>
      <td><a href="https://gradle.org/releases/">gradle下载</a></td>
    </tr>
    <tr>
      <td><a href="https://www.python.org/downloads/">python官网</a></td>
      <td><a href="https://github.com/pyenv/pyenv">pyenv</a>|<a href="https://github.com/pyenv-win/pyenv-win">pyenv-win</a></td>
    </tr>
    <tr>
      <td><a href="https://nodejs.org/en/download">nodejs官网</a></td>
      <td><a href="https://nodejs.org/dist/index.json">nodejs版本信息</a></td>
    </tr>
    <tr>
      <td><a href="https://www.rust-lang.org/tools/install">rust官网</a></td>
      <td><a href="https://www.rust-lang.org/zh-CN/tools/install">rust中文官网</a></td>
    </tr>
    <tr>
      <td><a href="https://www.cygwin.com/">Cygwin官网</a></td>
      <td><a href="https://www.cygwin.com/install.html">Cygwin安装</a></td>
    </tr>
    <tr>
      <td><a href="https://www.msys2.org/">Msys2官网</a></td>
      <td><a href="https://mirrors.tuna.tsinghua.edu.cn/help/msys2/">Msys2清华镜像源</a></td>
    </tr>
    <tr>
      <td><a href="https://github.com/microsoft/vcpkg-tool">vcpkg-tool</a></td>
      <td><a href="https://github.com/microsoft/vcpkg">vcpkg</a></td>
    </tr>
    <tr>
      <td><a href="https://julialang.org/">julia官网</a></td>
      <td><a href="https://cn.julialang.org/">julia中文社区</a></td>
    </tr>
    <tr>
      <td><a href="https://vlang.io/">vlang官网</a></td>
      <td><a href="https://github.com/vlang/v">vlang github</a></td>
    </tr>
    <tr>
      <td><a href="https://typst.app/docs/">typst官网</a></td>
      <td><a href="https://github.com/typst/">typst github</a></td>
    </tr>
    <tr>
      <td><a href="https://flutter.dev/">flutter官网</a></td>
      <td><a href="https://mirrors.nju.edu.cn/flutter/flutter_infra_release/releases/">flutter南大镜像源</a></td>
    </tr>
    <tr>
      <td><a href="https://code.visualstudio.com/download">vscode官网</a></td>
      <td><a href="https://blog.csdn.net/feinifi/article/details/127697851">vscode国内CDN加速</a></td>
    </tr>
    <tr>
      <td><a href="https://neovim.io/">NeoVim官网</a></td>
      <td><a href="https://github.com/neovim">NeoVim github</a></td>
    </tr>
    <tr>
      <td><a href="https://brew.sh/">Homebrew官网</a></td>
      <td><a href="https://gitee.com/moqsien/gvc/raw/master/homebrew.sh">Homebrew安装脚本</a></td>
    </tr>
    <tr>
      <td><a href="https://github.com/moqsien/hackbrowser">浏览器数据管理</a></td>
      <td><a href="https://github.com/moonD4rk/HackBrowserData">浏览器数据管理github</a></td>
    </tr>
    <tr>
      <td><a href="https://github.com/moqsien/neobox">neobox免费梯子</a></td>
      <td><a href="https://github.com/SagerNet/sing-box">sing-box</a>|<a href="https://github.com/XTLS/Xray-core">xray-core</a>|<a href="https://github.com/moqsien/wgcf">wgcf</a></td>
    </tr>
    <tr>
      <td><a href="https://gitlab.com/ineo6/hosts/-/raw/master/next-hosts">github hosts加速</a></td>
      <td><a href="https://github.com/jianboy/github-host/raw/master/hosts">github hosts加速</a></td>
    </tr>
  </tbody>
</table>

---------
## 下载和安装
- 在[release](https://github.com/moqsien/gvc/releases)中下载最新版;
- 解压之后，双击可执行文件，或者在Terminal/PowerShell中执行该文件;
- 在新的Termnial/PowerShell中可以使用g或者gvc命令了，例如：g help;

---------
## 主要功能
- 各种语言编译器/解释器自动安装，多版本管理，自动配置国内资源加速；
- vscode自动安装/更新，配置/插件信息同步到webdav网盘
- neovim自动安装，配置信息同步到webdav网盘
- github hosts加速
- homebrew一键安装
- 浏览器数据管理(各种常见浏览器的书签、插件、本地密码)，加密后同步到webdav网盘
- 可用性较高的免费梯子

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
