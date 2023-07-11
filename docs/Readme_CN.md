## [En](https://github.com/moqsien/gvc)
---------
![logo](https://github.com/moqsien/gvc/blob/main/docs/logo.png)

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
- 各种语言的编译器/解释器的自动安装，多版本管理，环境变量配置，自动配置国内资源加速等；
- vscode自动安装/更新，配置/插件信息同步到webdav网盘
- neovim自动安装，配置信息同步到webdav网盘
- github hosts加速
- homebrew一键安装
- 浏览器数据管理(各种常见浏览器的书签、插件、本地密码)，数据加密后同步到webdav网盘
- 可用性较高的免费梯子

有了gvc，你基本可以不需要再关心去哪里下载编译器/解释器，如何配置环境变量等等。<br/>
有了gvc，你可以把自己本地的开发环境的配置同步到任何支持webdav的网盘，然后在新的机器上一键重建自己熟悉的开发环境。<br/>
有了gvc，你可以轻松尝试某个语言的不同版本，进行多版本来回切换，目前支持多版本管理的有go，java，python，nodejs，julia，flutter(dart)。<br/>
有了gvc，你可以轻松重建vscode，包括配置和习惯使用的插件，这些都能保存在webdav。<br/>
有了gvc，你可以拥有可用性较高的免费梯子，平时google查资料，github浏览可以无忧。<br/>
有了gvc，你的常见浏览器的数据可以同步到webdav网盘，在没有代理或者没有google账号的情况下，你也可以把浏览器同步到webdav网盘，从而实现书签、密码、插件信息的备份，而且密码数据会自动加密，提升安全性。<br/>
总之，gvc存在的目标就是为了搞定那些开发过程中的各种环境配置之类的琐事，尤其是如果你喜欢尝鲜各种操作系统或者新的机器，那么gvc能帮你节省很多时间。<br/>
gvc天生比较适合爱折腾，持续学习，喜欢尝试新事物的程序猿。反之，绕行即可。<br/>
<br/>
gvc是个命令行工具，别问为啥没有GUI，因为确实没必要！

---------
## 具体功能
安装成功之后，**打开一个新的终端或者PowerShell**，可以执行g help命令，就能看到gvc的帮助信息。
例如，Windows下，在PowerShell中，就能得到类似于如下的信息：
```bash
PS C:\Users\moqsien> g help
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
   neobox-runner, nbrunner, nbr   Start a neobox client.
   neobox-keeper, nbkeeper, nbk   Start a neobox keeper.
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

以上为gvc的所有一级子命令的帮助信息。
如果要查看某个一级子命令的二级子命令的帮助信息，例如, 如果要查看go子命令的帮助信息，则可以使用g go help命令。
Windows下，执行g go help，就能得到类似于如下的信息：
```bash
PS C:\Users\moqsien> g go help
NAME:
   g.exe go - Go version management.

USAGE:
   g.exe go command [command options] [arguments...]

COMMANDS:
   remote, r                                Show remote versions.
   use, u                                   Download and use version.
   local, l                                 Show installed versions.
   remove-unused, ru                        Remove unused versions.
   remove-version, rm                       Remove a version.
   add-envs, env, e, ae                     Add envs for go.
   search-package, sp, search               Search for third-party packages.
   build, bui, b                            Compiles go code for multi-platforms [with <-ldflags "-s -w"> builtin].
   list-distributions, list-dist, dist, ld  List the platforms supported by go compilers.
   help, h                                  Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
以此类推，如果有三级子命令，也可以用相应的help命令常看帮助文档。<br/>
**如何使用免费梯子？**
- 使用g neobox-shell(或者使用简写命令，例如g ns等)，打开neobox的shell；
进入neobox的shell之后，有">>>"提示符，例如：
```bash
PS C:\Users\moqsien> g ns
>>>
```
- 在shell中，输入help，回车，查看shell的帮助，看看shell提供了哪些可用的命令；
neobox shell提供的命令类似如下：
```bash
PS C:\Users\moqsien> g ns
>>> help

Commands:
  add            Add proxies to neobox mannually.
  cfips          download/update valid cloudflare ips.
  clear          clear the screen
  exit           exit the program
  export         Export vpn history list.
  filter         Filter vpns by verifier.
  gc             Start GC manually.
  geoinfo        Install/Update geoip&geosite for sing-box.
  help           display help
  parse          Parse raw proxy URIs to human readable ones.
  pingunix       Setup ping without root for Unix-like OS.
  restart        Restart the running sing-box client with a chosen vpn. [restart vpn_index]
  setkey         Setup rawlist encrytion key for neobox. [With no args will set key to default value]
  show           Show neobox info.
  start          Start an sing-box client/keeper.
  stop           Stop the running sing-box client/keeper.
  system         enable current vpn as system proxy. [disable when an arg is provided]
  wireguard      register wireguard account and update licenseKey to warp plus [if a licenseKey is specified].


>>>
```
具体使用方法，可以查看[neobox](https://github.com/moqsien/neobox)文档。文档中，包括了如何获取aes-key(**必须**，否则无法使用)等注意事项。

**注意**：neobox-runner和neobox-keeper两个子命令，用户无需关心，它们仅仅是给neobox-shell使用的。neobox-shell中的命令已经提供了交互式操作，用于控制在后台运行的neobox。

---------
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
