## [中文](https://github.com/moqsien/gvc/blob/main/docs/Readme_CN.md)
<!-- ![logo](https://github.com/moqsien/gvc/blob/main/docs/logo.png) -->
<img src="https://github.com/moqsien/gvc/blob/main/docs/logo.png" width="30%">

- [中文](#中文)
- [What's GVC?](#whats-gvc)
- [What's supported?](#whats-supported)
- [Download \& Install](#download--install)
- [Main features](#main-features)
  - [Subcommand: go](#subcommand-go)
  - [Subcommand: proto](#subcommand-proto)
  - [Subcommand: python](#subcommand-python)
  - [Subcommand: java](#subcommand-java)
  - [Subcommand: cpp](#subcommand-cpp)
  - [Subcommand: vlang](#subcommand-vlang)
  - [Subcommand: vscode](#subcommand-vscode)
  - [Subcommand: hosts](#subcommand-hosts)
  - [Subcommand: github](#subcommand-github)
  - [Subcommand: browser](#subcommand-browser)
  - [Subcommand: asciinema](#subcommand-asciinema)
  - [Subcommand: cloc](#subcommand-cloc)
  - [Subcommand: config](#subcommand-config)
  - [Subcommand: neobox-shell](#subcommand-neobox-shell)
    - [Supported command in neobox-shell](#supported-command-in-neobox-shell)
- [special statement](#special-statement)
- [Demo](#demo)
- [thanks to](#thanks-to)
- [buy me a coffee](#buy-me-a-coffee)

## What's GVC?
At the very beginning, GVC is just the abbreviation for Go-Version-Controller, which means, it provides auto-installation, environment variables handling, as well as multi-versions management only for Go compilers.

As we know, we already have [gvm](https://github.com/andrewkroh/gvm) or [g](https://github.com/voidint/g) with the similar features implemented. So, why do we need a new one?

The reason to create GVC is for more convenience and a better UI(maybe TUI more pricisely).

However, this never becomes the end of the story. After the version-management for Go has been implemented, an idea for managing other languages flashes across my mind. Therefore, GVC starts to support version-control and auto-installation also for **Java/Python/NodeJS/Flutter/Julia/Protoc/Rust/Cpp/Vlang/Typst**.

At this point, GVC becomes **General-Version-Controller**.

And the story still continues.

Auto-installation for Visual Studio Code(**VSCode**) and **NeoVim** is adopted.The **WebDAV** support is also introduced for saving config files from VSCode/NeoVim to user's netdisk(eg. jianguoyun.com). So, you can rebuild your Development Environment using these files on any machine.

Besides, GVC also supports [asciinema](https://asciinema.org/) terminal recording, browser data management, counting lines of code, etc.

Finally, GVC becomes something just like a **Scaffolding Tool for local development environments management**.

---------
## What's supported?

<table>
  <thead>
    <tr>
      <th>Lang/App</th>
      <th>additions</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><a href="https://go.dev/">go homepage</a></td>
      <td><a href="https://golang.google.cn/">go homepage cn</a>|<a href="https://mirrors.aliyun.com/golang/">aliyun mirror</a></td>
    </tr>
    <tr>
      <td><a href="https://www.oracle.com/java/technologies/downloads/">java/jdk homepage</a></td>
      <td><a href="https://www.injdk.cn/">java/jdk cn</a></td>
    </tr>
    <tr>
      <td><a href="https://maven.apache.org/download.cgi">maven homepage</a></td>
      <td><a href="https://dlcdn.apache.org/maven/">maven downloads</a></td>
    </tr>
    <tr>
      <td><a href="https://gradle.org/install/">gradle homepage</a></td>
      <td><a href="https://gradle.org/releases/">gradle releases</a></td>
    </tr>
    <tr>
      <td><a href="https://www.python.org/downloads/">python homepage</a></td>
      <td><a href="https://github.com/pyenv/pyenv">pyenv</a>|<a href="https://github.com/pyenv-win/pyenv-win">pyenv-win</a></td>
    </tr>
    <tr>
      <td><a href="https://nodejs.org/en/download">nodejs homepage</a></td>
      <td><a href="https://nodejs.org/dist/index.json">nodejs versions</a></td>
    </tr>
    <tr>
      <td><a href="https://www.rust-lang.org/tools/install">rust homepage</a></td>
      <td><a href="https://www.rust-lang.org/zh-CN/tools/install">rust homepage cn</a></td>
    </tr>
    <tr>
      <td><a href="https://www.cygwin.com/">Cygwin homepage</a></td>
      <td><a href="https://www.cygwin.com/install.html">Cygwin installation</a></td>
    </tr>
    <tr>
      <td><a href="https://www.msys2.org/">Msys2 homepage</a></td>
      <td><a href="https://mirrors.tuna.tsinghua.edu.cn/help/msys2/">Msys2 tsinghua mirror</a></td>
    </tr>
    <tr>
      <td><a href="https://github.com/microsoft/vcpkg-tool">vcpkg-tool</a></td>
      <td><a href="https://github.com/microsoft/vcpkg">vcpkg</a></td>
    </tr>
    <tr>
      <td><a href="https://julialang.org/">julia homepage</a></td>
      <td><a href="https://cn.julialang.org/">julia community cn</a></td>
    </tr>
    <tr>
      <td><a href="https://vlang.io/">vlang homepage</a></td>
      <td><a href="https://github.com/vlang/v">vlang github</a></td>
    </tr>
    <tr>
      <td><a href="https://typst.app/docs/">typst homepage</a></td>
      <td><a href="https://github.com/typst/">typst github</a></td>
    </tr>
    <tr>
      <td><a href="https://flutter.dev/">flutter homepage</a></td>
      <td><a href="https://mirrors.nju.edu.cn/flutter/flutter_infra_release/releases/">flutter nju mirror</a></td>
    </tr>
    <tr>
      <td><a href="https://code.visualstudio.com/download">vscode homepage</a></td>
      <td><a href="https://blog.csdn.net/feinifi/article/details/127697851">vscode cdn acceleration</a></td>
    </tr>
    <tr>
      <td><a href="https://neovim.io/">NeoVim homepage</a></td>
      <td><a href="https://github.com/neovim">NeoVim github</a></td>
    </tr>
    <tr>
      <td><a href="https://brew.sh/">Homebrew homepage</a></td>
      <td><a href="https://gitee.com/moqsien/gvc/raw/master/homebrew.sh">Homebrew shell script</a></td>
    </tr>
    <tr>
      <td><a href="https://github.com/moqsien/hackbrowser">Browser data management</a></td>
      <td><a href="https://github.com/moonD4rk/HackBrowserData">Browser data management github</a></td>
    </tr>
    <tr>
      <td><a href="https://github.com/moqsien/neobox">neobox vpns</a></td>
      <td><a href="https://github.com/SagerNet/sing-box">sing-box</a>|<a href="https://github.com/XTLS/Xray-core">xray-core</a>|<a href="https://github.com/moqsien/wgcf">wgcf</a></td>
    </tr>
    <tr>
      <td><a href="https://gitlab.com/ineo6/hosts/-/raw/master/next-hosts">github hosts file</a></td>
      <td><a href="https://github.com/jianboy/github-host/raw/master/hosts">github hosts file</a></td>
    </tr>
    <tr>
      <td><a href="https://asciinema.org/">asciinema</a></td>
      <td><a href="https://github.com/moqsien/asciinema">asciinema for full-platform</a></td>
    </tr>
    <tr>
      <td><a href="https://github.com/hhatto/gocloc">count lines of code(cloc)</a></td>
      <td><a href="https://github.com/hhatto/gocloc">cloc</a></td>
    </tr>
    <tr>
      <td><a href="https://github.com/protocolbuffers/protobuf">protobuf</a></td>
      <td><a href="https://google.golang.org/protobuf/cmd/protoc-gen-go">protoc-go-gen</a></td>
    </tr>
  </tbody>
</table>

## Download & Install
- Download the latest [release](https://github.com/moqsien/gvc/releases).
- Unzip, double click the executable file, or run executable file in Terminal/PowerShell.
- Open a new Terminal/PowerShell, then the command **g** is available. Help info will be displayed using **'g help'**.

---------
## Main features
GVC is a command-line tool, use "g help" or "gvc help", to see help info.

```bash
$moqsien> g help

NAME:
   g.exe - gvc <Command> <SubCommand>...

USAGE:
   g.exe [global options] command [command options] [arguments...]

DESCRIPTION:
   A productive tool to manage your development environment.

COMMANDS:
   go, g                            Go version management.
   proto, protobuf, protoc, pt      Protoc installation.
   python, py                       Python version management.
   java, jdk, j                     Java jdk version management.
   maven, mav, ma                   Maven version management.
   gradle, gra, gr                  Gradle version management.
   nodejs, node, no                 NodeJS version management.
   flutter, flu, fl                 Flutter version management.
   julia, jul, ju                   Julia version management.
   rust, rustc, ru, r               Rust installation.
   cpp                              C/C++ management.
   typst, ty                        Typst installation.
   vlang, vl                        Vlang installation.
   vscode, vsc, vs, v               VSCode and extensions installation.
   nvim, neovim, nv, n              Neovim installation.
   neobox-shell, shell, box, ns     Start a neobox shell.
   neobox-runner, nbrunner, nbr     Start a neobox client.
   neobox-keeper, nbkeeper, nbk     Start a neobox keeper.
   browser, br                      Browser data management.
   homebrew, brew, hb               Homebrew installation or update.
   gsudo, winsudo, gs, ws           Gsudo for windows.
   hosts, h, host                   System hosts file management(need admistrator or root).
   github, gh                       Github download speedup.
   cloc, cl                         Count lines of code.
   asciinema, ascii, asc            Asciinema terminal recorder.
   config, conf, cnf, c             Config file management for gvc.
   version, ver, vsi                Show gvc version info.
   check, checklatest, checkupdate  Check and download the latest version of gvc.
   show, sho, sh                    Show [gvc] installation path and config file path.
   uninstall, unins, delete, del    [Caution] Remove gvc and softwares installed by gvc!
   help, h                          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

```bash
$moqsien> g version

 ██████  ██    ██  ██████
██           ██    ██ ██
██   ███  ██    ██ ██
██    ██    ██  ██  ██
 ██████    ████    ██████

┌────────────────────────────────────────────────────────────────────┐
|                                                                    |
|                                                                    |
|     Version:     v1.4.2(f684b2a1a57c560228add15590783d428d92b480)  |
|     UpdateAt:    Wed Aug 23 17:33:08 2023 +0800                    |
|     Homepage:    https://github.com/moqsien/gvc                    |
|     Email:       moqsien@foxmail.com                               |
|                                                                    |
|                                                                    |
└────────────────────────────────────────────────────────────────────┘
```

### Subcommand: go
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
In this subcommand, you can show/install/remove/change go compiler versions, search third-party packages written in Go, and build Go source code for multi-platforms without prepare any scripts.
You can also rename a local go module using subcommand: **g go rto NEW_MODULE_NAME**.

### Subcommand: proto
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
In this subcommand, you can auto-install protoc, protoc-gen-go and protoc-gen-go-grpc.

### Subcommand: python
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
   help, h             Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
This subcommand benifits a lot from pyenv/pyenv-win. 

### Subcommand: java
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
In this subcommand, the option "-z" if for users in China.

### Subcommand: cpp
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
This subcommand is for Windows users. It will install Msys2 or Cygwin, just to your preference.
You can also install the Cpp-package-manager [vcpkg](https://github.com/microsoft/vcpkg) maintained by Microsoft.

### Subcommand: vlang
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
Vlang is a fantastic new language with high performance. This subcommand will install or update to the latest version of Vlang. You can also install [v-analyser](https://github.com/v-analyzer/v-analyzer), which brings vlang the code completion/IntelliSense/go to definition features for VSCode and other editors. If you have VSCode installed, this subcommand will automatically install related extensions and config the settings for you.

### Subcommand: vscode
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
This subcommand will install VSCode for you. You can also install extensions using the extension-info-files saved to WebDAV by GVC. An adapter for git tools in Msys2/Cygwin to VSCode usage is available in this subcommand. You can easily make use of the git tool from either Msys2 or Cygwin in VSCode.

### Subcommand: hosts
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
This subcommand will automatically update the hosts file. Main purpose of this subcommand is to speedup visits to github, microsoft, steam, etc.

### Subcommand: github
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
This subcommand speedups github downloadings in China.

### Subcommand: browser
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
This subcommand handles browser data, save data to WebDAV. 

### Subcommand: asciinema
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
This subcommand provides terminal recording features for both **Powershell** and **Unix-Like Shells**.

### Subcommand: cloc
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
This subcommand provides CLOC(Count Lines of Code) features.

### Subcommand: config
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
The subcommand webdav will interactively direct you to configure your WebDAV Account info and secrets for encrytion of browser data.
The subcommand push and pull will interact with remote WebDAV.
The subcommand reset will reset gvc-config-files to default values.

### Subcommand: neobox-shell
```bash
PS C:\Users\moqsien> g neobox-shell help
NAME:
   g.exe neobox-shell - Start a neobox shell.

USAGE:
   g.exe neobox-shell [command options] [arguments...]

OPTIONS:
   --help, -h  show help
```
This subcommand will start the neobox-shell.
[neobox](https://github.com/moqsien/neobox) provides some available free VPNs for user.

#### Supported command in neobox-shell
```bash
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
```
Note: You should read docs(in Chinese) for neobox, for details, see [neobox](https://github.com/moqsien/neobox).

---------
## special statement
gvc provides no paid services, so, users should make use of it within the limits permitted by law in his/her country.

---------
## Demo
- gvc installation
[![asciicast](https://asciinema.org/a/597749.svg)](https://asciinema.org/a/597749)

- go version management.
[![asciicast](https://asciinema.org/a/597750.svg)](https://asciinema.org/a/597750)

- [neobox](https://github.com/moqsien/neobox) free vpns.
[![asciicast](https://asciinema.org/a/597753.svg)](https://asciinema.org/a/597753)

- vscode installation, and vscode extensions installation
[![asciicast](https://asciinema.org/a/597755.svg)](https://asciinema.org/a/597755)

## thanks to
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
- [asciinema](https://github.com/securisec/asciinema)
- [PowerSession-rs](https://github.com/Watfaq/PowerSession-rs)
- [conpty-go](https://github.com/qsocket/conpty-go)
- [gocloc](https://github.com/hhatto/gocloc)

## buy me a coffee
[wechat](https://github.com/moqsien/moqsien/blob/main/imgs/wechat.jpeg)
[alipay](https://github.com/moqsien/moqsien/blob/main/imgs/alipay.jpeg)
