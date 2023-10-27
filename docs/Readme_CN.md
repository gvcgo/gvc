## [En](https://github.com/moqsien/gvc)
<!-- ![logo](https://github.com/moqsien/gvc/blob/main/docs/logo.png)-->
<img src="https://github.com/moqsien/gvc/blob/main/docs/logo.png" width="30%">

- [En](#en)
- [gvc是什么？](#gvc是什么)
- [gvc支持哪些语言或应用？](#gvc支持哪些语言或应用)
- [下载和安装](#下载和安装)
- [功能概览](#功能概览)
  - [gpt子命令](#gpt子命令)
  - [go子命令](#go子命令)
  - [proto子命令](#proto子命令)
  - [python子命令](#python子命令)
  - [java子命令](#java子命令)
  - [cpp子命令](#cpp子命令)
  - [vlang子命令](#vlang子命令)
  - [vscode子命令](#vscode子命令)
  - [hosts子命令](#hosts子命令)
  - [github子命令](#github子命令)
  - [git-XXX 子命令](#git-xxx-子命令)
  - [browser子命令](#browser子命令)
  - [asciinema子命令](#asciinema子命令)
  - [cloc子命令](#cloc子命令)
  - [config子命令](#config子命令)
  - [neobox-shell子命令](#neobox-shell子命令)
    - [neobox-shell内部命令](#neobox-shell内部命令)
- [特别申明](#特别申明)
- [Demo](#demo)
- [感谢](#感谢)
- [送我一杯咖啡~~~](#送我一杯咖啡)

---------
gvc交流QQ群：

<img src="https://github.com/moqsien/neobox/blob/main/docs/gvc_qq_group.jpg" width="30%">

---------
## gvc是什么？
最开始，gvc是go version controller的缩写。当时只是想做一个好用一点，界面相对美观一点的go多版本管理工具。

原因在于[gvm](https://github.com/andrewkroh/gvm)安装都需要自行设置环境变量，而且长时间不更新(不过，我也没使用过它)。
[g](https://github.com/voidint/g)可以正常使用，也能满足需求，同时也支持Mac/Win/Linux，是个不错的选择。
但是g显示版本是一行一条，当版本很多时，上下翻找体验确实不好。
再者，g只负责安装某个版本，不负责配置环境变量。也就是说，如果你的电脑上之前没有安装过go编译器，你还得自行配置GOPROXY, GOPATH,GOBIN等等。
另外，g默认从[go.dev](https://go.dev)下载，对于中国大陆的用户不太友好，想要使用[golang.google.cn](https://golang.google.cn/)加速还得自行设置环境变量。
综合考虑，所以自行写了gvc。旨在一键搞定所有，方便好用。

值得提一嘴的是，在Windows下搞定自动配置环境变量，是件比较麻烦的事情。因为Windows通过setx设置环境变量之后，首先容易被接下来的setx命令覆盖，其次是如果不重新打开一个Powershell，
是无法读取到刚刚通过setx设置的环境变量的。所以，**Windows用户在使用过程中，如果通过gvc安装了软件，记得重新打开Powershell**。

后来，觉得其他语言也可以有类似的功能。因为作者平时使用的语言主要就有go/python/typescript等等，还会看看rust，vlang，c/cpp之类的。所以，为什么gvc不可以支持一下这些语言呢？
说干就干，一顿操作下来，gvc最终支持的语言有go/java/python/nodejs/flutter(dart)/julia/rust/cpp/vlang/typst。

所以，gvc可以说是General version controller的缩写。

不过，光有语言的编译器/解释器，总觉得还缺点什么？
是的，缺IDE/Editor。所以，gvc支持了一键安装VSCode和Neovim。二者可以通过插件进行搭配，轻松实现多语言开发。

但是，很多小伙伴觉得VSCode配置麻烦，使用起来不如JetBrains系列方便。没问题，安排！
gvc增加了WebDAV协议的网盘同步功能，它能把你的本地配置一键同步到WebDAV网盘(例如坚果云)中。这样的话，你只需要配置一次就好了。
在其他任何地方，你都可以使用WebDAV中的配置信息，通过gvc恢复你熟悉的VSCode配置。这些包括你常用的VSCode插件，你熟悉的Keybindings配置，你的settings.json等。

至此，gvc已经帮你解决了各种编程语言以及IDE的麻烦。然而gvc的功能还远未介绍完毕。

考虑到在家学习时，经常需要上google和github。所以gvc有了免费的梯子。通过作者之前写的一个交互式shell库，实现了命令行客户端。筛选免费梯子的效率高过很多客户端。
既然要使用github，也不能完全指望免费梯子，毕竟不太稳定。所以，gvc必须有一键修饰hosts文件，加速github访问的功能。另外，gvc还有github下载加速，只需要提供某个github项目的主页，即可选择加速下载源码或者最新的release了，无需代理或者找加速网站。

此外，gvc还提供了代码统计(Count Lines of Code)， asciinema终端录制和上传，浏览器数据网盘同步等等功能。

目前为止，**gvc可以说已经成为一个跨平台，多机器开发环境管理的脚手架工具**。

---------
## gvc支持哪些语言或应用？

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
    <tr>
      <td><a href="https://asciinema.org/">asciinema终端录频</a></td>
      <td><a href="https://github.com/moqsien/asciinema">全平台支持</a></td>
    </tr>
    <tr>
      <td><a href="https://github.com/hhatto/gocloc">项目代码统计功能</a></td>
      <td><a href="https://github.com/hhatto/gocloc">支持各种语言</a></td>
    </tr>
    <tr>
      <td><a href="https://github.com/protocolbuffers/protobuf">protobuf</a></td>
      <td><a href="https://google.golang.org/protobuf/cmd/protoc-gen-go">protoc-go-gen</a></td>
    </tr>
    <tr>
      <td><a href="https://docs.docker.com/desktop/install/windows-install/">docker-for-windows</a></td>
      <td><a href="https://docs.docker.com/desktop/install/mac-install/">docker-for-MacOS</a></td>
    </tr>
    <tr>
      <td><a href="https://github.com/go-git/go-git">可以使用本地代理的git命令</a></td>
      <td><a href="https://github.com/go-git/go-git">go-git，比git配置全局命令更方便</a></td>
    </tr>
    <tr>
      <td><a href="https://openai.com/">openai</a></td>
      <td><a href="https://github.com/sashabaranov/go-openai">go-openai</a></td>
    </tr>
    <tr>
      <td><a href="https://xinghuo.xfyun.cn/">讯飞星火</a></td>
      <td><a href="https://xinghuo.xfyun.cn/sparkapi">讯飞星火</a></td>
    </tr>
  </tbody>
</table>

---------
## 下载和安装
- 在[release](https://github.com/moqsien/gvc/releases)中下载最新版;
- 解压之后，双击可执行文件，或者在Terminal/PowerShell中执行该文件;
- 在新的Termnial/PowerShell中可以使用g或者gvc命令了，例如：g help;

- 或者使用go install安装
```bash
go install -tags "with_wireguard with_shadowsocksr with_utls with_gvisor with_grpc with_ech with_dhcp" github.com/moqsien/gvc@latest
```

---------
## 功能概览
安装成功之后，**打开一个新的终端或者PowerShell**，可以执行g help命令，就能看到gvc的帮助信息。
例如，Windows下，在PowerShell中，就能得到类似于如下的信息：
```bash
$moqsien> g help

NAME:
   g.exe - gvc <Command> <SubCommand>...

USAGE:
   g.exe [global options] command [command options] [arguments...]

DESCRIPTION:
   A productive tool to manage your development environment.

COMMANDS:
   go, g                                            Go version management.
   proto, protobuf, protoc, pt                      Protoc installation.
   python, py                                       Python version management.
   java, jdk, j                                     Java jdk version management.
   maven, mav, ma                                   Maven version management.
   gradle, gra, gr                                  Gradle version management.
   nodejs, node, no                                 NodeJS version management.
   flutter, flu, fl                                 Flutter version management. // 无需安装Android Studio，轻松配置VSCode进行Flutter安卓开发。
   julia, jul, ju                                   Julia version management.
   rust, rustc, ru, r                               Rust installation.
   cpp                                              C/C++ management.
   typst, ty                                        Typst installation.
   vlang, vl                                        Vlang installation.
   vscode, vsc, vs, v                               VSCode and extensions installation.
   nvim, neovim, nv, n                              Neovim installation.
   gpt-spark, gpt, gspark                           ChatGPT/Spark bot.
   neobox-shell, shell, box, ns                     Start a neobox shell.
   neobox-runner, nbrunner, nbr                     Start a neobox client.
   neobox-keeper, nbkeeper, nbk                     Start a neobox keeper.
   browser, br                                      Browser data management.
   homebrew, brew, hb                               Homebrew installation or update.
   gsudo, winsudo, gs, ws                           Gsudo for windows.
   hosts, h, host                                   Sytem hosts file management(need admistrator or root).
   git-set-proxy, gsproxy, gsp                      Set default proxy for git [default: http://localhost:2023]. // 设置git-*子命令的默认代理
   git-clone, gclone, gclo                          Git Clone using a proxy.
   git-pull, gpull, gpul                            Git Pull using a proxy.
   git-push, gpush, gpus                            Git Push using a proxy.
   git-commit-push, gcpush, gcp                     Git commit and push to remote using a proxy.
   git-add-tag-push, gaddtag, gatag, gat            Git add a new tag and push to remote using a proxy.
   git-del-tag-push, gdeltag, gdtag, gdt            Git delete a tag and push to remote using a proxy.
   git-show-tag-latest, gshowtaglatest, gstag, gst  Git show the latest tag of a local repository.
   win-git-install, wgit, wgi                       Install git for windows.
   github, gh                                       Github download speedup.
   cloc, cl                                         Count lines of code.
   asciinema, ascii, asc                            Asciinema terminal recorder.
   docker, dck, dock                                Docker installation.
   config, conf, cnf, c                             Config file management for gvc.
   ssh-files, sshf, ssh                             Backup your ssh files.
   version, ver, vsi                                Show gvc version info.
   check, checklatest, checkupdate                  Check and download the latest version of gvc.
   show, sho, sh                                    Show [gvc] installation path and config file path.
   uninstall, unins, delete, del                    [Caution] Remove gvc and softwares installed by gvc!
   help, h                                          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

```bash
$moqsien> g version

      Name: GVC
      Version: v1.5.7(c7d768d9)
      UpdateAt: Tue Sep 26 13:14:49 2023 +0800
      Homepage: https://github.com/moqsien/gvc
      Email: moqsien2022@gmail.com
```

### gpt子命令
```bash
$moqsien> g gpt help

NAME:
   g.exe gpt-spark - ChatGPT/Spark bot.

USAGE:
   g.exe gpt-spark [command options] [arguments...]

OPTIONS:
   --help, -h  show help
```
ChatGPT和讯飞星火大模型客户端。基于TUI。详细信息请查看[gogpt](https://github.com/moqsien/gogpt)。

### go子命令
```bash
$moqsien> g go help

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
   renameTo, rnt, rto                       Rename a local go module[gvc go rto NEW_MODULE_NAME].
   list-distributions, list-dist, dist, ld  List the platforms supported by go compilers.
   help, h                                  Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
go子命令可以安装、删除、切换版本。还能一键配置好诸如GOPATH、GOPROXY、GOBIN之类的环境变量。
还能通过search-package来搜索第三方库。
build子命令还提供了对于go build的增强，可以跨平台编译并压缩打包，无需编写任何脚本。
renameTo子命令能重命名一个本地go项目。

### proto子命令
```bash
$moqsien> g proto help

NAME:
   g.exe proto - Protoc installation.

USAGE:
   g.exe proto command [command options] [arguments...]

COMMANDS:
   install, ins, i                  Install protoc.
   install-go-plugin, igo, ig       Install protoc-gen-go.
   install-grpc-plugin, igrpc, igr  Install protoc-gen-go-grpc.
   help, h                          Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
一键安装protoc, protoc-gen-go以及protoc-gen-go-grpc。

### python子命令
```bash
$moqsien> g python help

NAME:
   g.exe python - Python version management.

USAGE:
   g.exe python command [command options] [arguments...]

COMMANDS:
   remote, r           Show remote versions.
   use, u              Download and use a version.
   local, l            Show installed versions.
   remove-version, rm  Remove a version.
   update, up          Install or update pyenv.
   path, pth           Show pyenv versions path.
   rmfix, rfix         Automatically remove python.exe generated by Windows system in ~/AppData/Local/Microsoft/WindowsApps .
   help, h             Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
python多版本管理采用了现成的pyenv和pyenv-win脚本，对其输出做了优化，另外针对大陆下载慢的情况进行了改进。安装python时，会自动配置好pip的国内加速源。
Windows11有时候经常在~/AppData/Local/Microsoft/WindowsApps中自动生成python.exe，干扰了pyenv-win的python的环境变量，这里使用rmfix提供一键删除。

### java子命令
```bash
$moqsien> g java help

NAME:
   g.exe java - Java jdk version management.

USAGE:
   g.exe java command [command options] [arguments...]

COMMANDS:
   use, u                  Download and use jdk. <Command> {gvc jdk use [-z] xxx}
   remote, r               Show available versions.  <Command> {gvc jdk remote [-z]}
   local, l                Show installed versions.
   remove, rm              Remove an installed version.
   remove-unused, rmu, ru  Remove unused versions.
   help, h                 Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
大陆用户请使用-z参数，这将从国内资源下载jdk，支持的版本也比较全。否则只能从官网下载最新版。

### cpp子命令
```bash
$moqsien> g cpp help

NAME:
   g.exe cpp - C/C++ management.

USAGE:
   g.exe cpp command [command options] [arguments...]

COMMANDS:
   install-msys2, insm, im                Install the latest msys2.
   uninstall-msys2, unim, um, remove, rm  Uninstall msys2.
   install-cygwin, insc, ic               Install Cygwin.
   install-vcpkg, insv, iv                Install vcpkg.
   help, h                                Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
主要用于windows下，安装Cygwin或者Msys2。会默认安装git之类的工具。
另外vcpkg是微软开源的cpp包管理器，与python的pip类似。

### vlang子命令
```bash
$moqsien> g vlang help

NAME:
   g.exe vlang - Vlang installation.

USAGE:
   g.exe vlang command [command options] [arguments...]

COMMANDS:
   install, ins, i             Install Vlang.
   install-analyzer, insa, ia  Install v-analyzer and related extension for vscode.
   setenv, env, se, e          Set env for Vlang.
   help, h                     Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
vlang是一个性能强悍语法简单的新兴语言，此命令可以一键安装/更新到vlang的最新版。
install-analyzer会下载安装vlang的语法解析器以及相应的VSCode插件，如果系统已经安装了VSCode的话。
这样VSCode就可以有vlang的语法高亮和自动补全之类的功能了。

### vscode子命令
```bash
$moqsien> g vscode help

NAME:
   g.exe vscode - VSCode and extensions installation.

USAGE:
   g.exe vscode command [command options] [arguments...]

COMMANDS:
   install, i, ins                    Automatically install vscode.
   install-extensions, ie, iext       Automatically install extensions for vscode.
   use-msys2-cygwin-git, use-git, ug  Repair and make use of git.exe from Msys2/Cygwin.
   help, h                            Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
一键安装VSCode，支持国内CDN加速，解决中国大陆下载慢的问题。
根据WebDAV中保存的VSCode插件信息，一键自动安装插件。
Cygwin或者Msys2中的git不能直接被VSCode识别，通过use-msys2-cygwin-git命令可以一键解决此问题。

### hosts子命令
```bash
$moqsien> g hosts help

NAME:
   g.exe hosts - Sytem hosts file management(need admistrator or root).

USAGE:
   g.exe hosts command [command options] [arguments...]

COMMANDS:
   fetch, f  Fetch github hosts info.
   show, s   Show hosts file path.
   help, h   Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
一键更新hosts文件，加速github、microsoft、steam等的访问。

### github子命令
```bash
$moqsien> g github help

NAME:
   g.exe github - Github download speedup.

USAGE:
   g.exe github command [command options] [arguments...]

COMMANDS:
   download, dl, d        Download files from github project.
   openbrowser, open, ob  Open acceleration website in browser.
   help, h                Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
github下载加速。只需要提供github项目主页连接即可。
使用g github download -code "https://github.com/moqsien/gvc"会下载gvc的源码压缩文件。
使用g github download "https://github.com/moqsien/gvc"会下载gvc的最新release。

### git-XXX 子命令
```bash
$moqsien> g help
...
git-clone, gclone, gclo                Git Clone using a proxy.
git-pull, gpull, gpul                  Git Pull using a proxy.
git-push, gpush, gpus                  Git Push using a proxy.
git-commit-push, gcpush, gcp           Git commit and push to remote using a proxy.
git-add-tag-push, gaddtag, gatag, gat  Git add a new tag and push to remote using a proxy.
git-del-tag-push, gdeltag, gdtag, gdt  Git delete a tag and push to remote using a proxy.
...

```

使用代理加速的git相关命令。还有诸如commit+push，tag+push等组合命令。
如果不特别指定代理，则默认使用neobox提供的免费代理"http://localhost:2023"。组合命令还可以不使用代理，这样在一般情况下可以少敲一点命令。

### browser子命令
```bash
$moqsien> g browser help

NAME:
   g.exe browser - Browser data management.

USAGE:
   g.exe browser command [command options] [arguments...]

COMMANDS:
   show-info, show, sh  Show supported browsers and data restore dir.
   push, psh, pu        Push browser Bookmarks/Password/ExtensionInfo to webdav.
   save, sa, s          Save browser Bookmarks/Password/ExtensionInfo to local dir.
   pull, pul, pl        Pull browser data from webdav to local dir.
   help, h              Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
一键自动管理浏览器数据。支持多款基于chromium的浏览器以及firefox。可以将浏览器数据，例如书签、插件列表、本地账号密码(加密处理)上传的自己的WebDAV网盘。

### asciinema子命令
```bash
$moqsien> g asciinema help

NAME:
   g.exe asciinema - Asciinema terminal recorder.

USAGE:
   g.exe asciinema command [command options] [arguments...]

COMMANDS:
   record, rec, r  Record terminal operations.
   play, pl, p     Play local asciinema file.
   auth, au, a     Bind local install-id to your asciinem.org account.
   upload, up, u   Upload local asciinema file to asciinema.org.
   help, h         Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
asciinema.org终端录制功能。可以录制、播放。并支持授权之后进行上传。写文档的利器。

### cloc子命令
```bash
$moqsien> g cloc help

NAME:
   g.exe cloc - Count lines of code.

USAGE:
   g.exe cloc [command options] [arguments...]

OPTIONS:
   --by-file, --bf                                    Report results for every encountered source file. (default: false)
   --debug, --de, -d                                  Dump debug log for developer. (default: false)
   --skip-duplicated, --skipdup, --sd                 Skip duplicated files. (default: false)
   --show-lang, --shlang, --sl                        Print about all languages and extensions. (default: false)
   --sort-tag value, --sort value, --st value         Sort based on a certain column["name", "files", "blank", "comment", "code"]. (default: "name")
   --output-type value, --output value, --ot value    Output type [values: default,cloc-xml,sloccount,json]. (default: "default")
   --exclude-ext value, --excl value, --ee value      Exclude file name extensions (separated commas).
   --include-lang value, --langs value, --il value    Include language name (separated commas).
   --match value, --mat value, -m value               Include file name (regex).
   --not-match value, --nmat value, --nm value        Exclude file name (regex).
   --match-dir value, --matd value, --md value        Include dir name (regex).
   --not-match-dir value, --nmatd value, --nmd value  Exclude dir name (regex).
   --help, -h                                         show help
```
代码统计功能。支持对单个文件或者某个项目进行统计。支持正则表达式排除文件或文件夹等等。
例如通过cloc功能的统计，gvc自身代码以及其他独立出去的库的代码，总计已经将近19k。

### config子命令
```bash
$moqsien> g config help

NAME:
   g.exe config - Config file management for gvc.

USAGE:
   g.exe config command [command options] [arguments...]

COMMANDS:
   webdav, dav, w  Setup webdav account info.
   pull, pl        Pull settings from remote webdav and apply them to applications.
   push, ph        Gather settings from applications and sync them to remote webdav.
   reset, rs, r    Reset the gvc config file to default values.
   help, h         Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
这是配置文件统一管理的子命令。
webdav子命令用于配置自己的webdav信息，包括账户、密码、host以及用于加密本地浏览器数据的密码等等。
pull和push分别用于拉取和推送这些需要保存的信息。

### neobox-shell子命令
```bash
PS C:\Users\moqsien> g neobox-shell help
NAME:
   g.exe neobox-shell - Start a neobox shell.

USAGE:
   g.exe neobox-shell [command options] [arguments...]

OPTIONS:
   --help, -h  show help
```
免费梯子的交互式shell启动命令。

#### neobox-shell内部命令
```bash
>>> help

Commands:
  add            Add proxies to neobox mannually.
  added          Add edgetunnel proxies to neobox.
  cfip           Test speed for cloudflare IPv4s.
  clear          clear the screen
  dedge          Download rawList for a specified edgeTunnel proxy [dedge proxy_index].
  domain         Download selected domains file for edgeTunnels.
  exit           exit the program
  filter         Start filtering proxies by verifier manually.
  gc             Start GC manually.
  geoinfo        Install/Update geoip&geosite for neobox client.
  graw           Manually dowload rawUri list(conf.txt from gitlab) for neobox client.
  guuid          Generate UUIDs.
  help           display help
  parse          Parse rawUri of a proxy to xray-core/sing-box outbound string [xray-core by default].
  pingd          Ping selected domains for edgeTunnels.
  qcode          Generate QRCode for a chosen proxy. [qcode proxy_index]
  remove         Remove a manually added proxy [manually or edgetunnel].
  restart        Restart the running neobox client with a chosen proxy. [restart proxy_index]
  setkey         Setup rawlist encrytion key for neobox. [With no args will set key to default value]
  setping        Setup ping without root for Linux.
  show           Show neobox info.
  start          Start a neobox client/keeper.
  stop           Stop neobox client.
  sys-proxy      To enable or disable System Proxy.
  wireguard      Register wireguard account and update licenseKey to Warp+ [if a licenseKey is specified].
```
在交互式shell中，可以控制免费梯子的启停，筛选等等。
注意，**一定要去neobox项目认真阅读文档**，否则你可能无法使用免费梯子。具体请查看[neobox](https://github.com/moqsien/neobox)。

---------
## 特别申明
本项目不提供任何收费服务，请任何使用者自觉遵守本国法律法规。

---------
## Demo
- gvc 安装
[![asciicast](https://asciinema.org/a/597749.svg)](https://asciinema.org/a/597749)

- go自动安装和多版本管理.
[![asciicast](https://asciinema.org/a/597750.svg)](https://asciinema.org/a/597750)

- [neobox](https://github.com/moqsien/neobox)免费vpn.
[![asciicast](https://asciinema.org/a/597753.svg)](https://asciinema.org/a/597753)

- vscode安装/更新，以及vscode插件信息同步，自动根据同步的信息安装插件
[![asciicast](https://asciinema.org/a/597755.svg)](https://asciinema.org/a/597755)

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
- [bubbles](https://github.com/charmbracelet/bubbles)
- [lipgloss](https://github.com/charmbracelet/lipgloss)
- [pterm](https://github.com/pterm/pterm)
- [goutils](https://github.com/moqsien/goutils)
- [asciinema](https://github.com/securisec/asciinema)
- [PowerSession-rs](https://github.com/Watfaq/PowerSession-rs)
- [conpty-go](https://github.com/qsocket/conpty-go)
- [gocloc](https://github.com/hhatto/gocloc)
- [protobuf](https://github.com/protocolbuffers/protobuf)
- [docker](https://docs.docker.com/desktop/)
- [go-git](https://github.com/go-git/go-git)
- [gogpt](https://github.com/moqsien/gogpt)

## 送我一杯咖啡~~~
[wechat](https://github.com/moqsien/moqsien/blob/main/imgs/wechat.jpeg)

[alipay](https://github.com/moqsien/moqsien/blob/main/imgs/alipay.jpeg)
