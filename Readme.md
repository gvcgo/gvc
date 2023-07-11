## [中文](https://github.com/moqsien/gvc/blob/main/docs/Readme_CN.md)
---------

![logo](https://github.com/moqsien/gvc/blob/main/docs/logo.png)

### What's supported?

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
  </tbody>
</table>

### Download & Install
- Download the latest [release](https://github.com/moqsien/gvc/releases).
- Unzip, double click the executable file, or run executable file in Terminal/PowerShell.
- Open a new Terminal/PowerShell, then the command **g** is available. Help info will be displayed using **'g help'**.

---------
## Something nice about [gvc](https://github.com/moqsien/gvc).


---------
### Features
- Management for different programming languages, including go, java, python, node, julia, c/cpp, etc. It provides functionalities like auto-installation, multi-versions-management, envs-setup, etc.
- VSCode auto-installation/upgrade, and configs/extensions info synchronization to WebDAV.
- NeoVim auto-installation, and configs synchronization to WebDAV.
- hosts file modifications for github visit acceleration.
- Homebrew installation and setup.
- Browser data management. Bookmarks, plugins, local password synchronization to WebDAV.
- Free VPNs.

## Help Info
---------
GVC is a command line tool. 
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

**How to use free vpns?**
- Open **neobox shell**.
```bash
PS C:\Users\moqsien> g ns
>>>
```

- Show available command in **neobox shell**.
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
>>>
```
**See** [neobox docs](https://github.com/moqsien/neobox).

---------
## Demo
- gvc installation
[![asciicast](https://asciinema.org/a/gzILqlGMUpoNiyrfRdvnJ3eMz.svg)](https://asciinema.org/a/gzILqlGMUpoNiyrfRdvnJ3eMz)

- go version management.
[![asciicast](https://asciinema.org/a/2YY5Tk7YhJeKiccxjaWjA6mqE.svg)](https://asciinema.org/a/2YY5Tk7YhJeKiccxjaWjA6mqE)

- neobox free vpns.
[![asciicast](https://asciinema.org/a/Wy5S9kxZU1tL68Xz8Wlj8EJjc.svg)](https://asciinema.org/a/Wy5S9kxZU1tL68Xz8Wlj8EJjc)

- vscode installation, and vscode extensions installation
[![asciicast](https://asciinema.org/a/eTecjmXlHSpVZPC14OaDQk2e3.svg)](https://asciinema.org/a/eTecjmXlHSpVZPC14OaDQk2e3)

---------
## special statement
gvc provides no paid services, so, users should make use of it within the limits permitted by law in his/her country.

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

## buy me a coffee
[wechat](https://github.com/moqsien/moqsien/blob/main/imgs/wechat.jpeg)
[alipay](https://github.com/moqsien/moqsien/blob/main/imgs/alipay.jpeg)
