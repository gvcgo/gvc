## [中文](https://github.com/moqsien/gvc/blob/main/docs/Readme_CN.md)
---------

### What's supported?

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
## Gallery (In powershell for example.)

### Help info. 
![](https://github.com/moqsien/gvc/blob/main/docs/ghelp.png)

### Subcommands available for go management. 
![](https://github.com/moqsien/gvc/blob/main/docs/goRemote.png)

### Show the versions available for go compilers. 
![](https://github.com/moqsien/gvc/blob/main/docs/goRemoteShow.png)

### Installation of go compiler.
![](https://github.com/moqsien/gvc/blob/main/docs/goInstall.png)
![](https://github.com/moqsien/gvc/blob/main/docs/goInstallationFinished.png)

### Search third party libraries written in go using keyword.
![](https://github.com/moqsien/gvc/blob/main/docs/goSearch.png)

### Show available versions of compiler/interpreter for other languages.
![](https://github.com/moqsien/gvc/blob/main/docs/pyNodeFlutterJulia.png)

### "go build" enhancement from gvc.
![beforeBuild](https://github.com/moqsien/gvc/blob/main/docs/beforebuild.png)
![chooseOptionsForBuild](https://github.com/moqsien/gvc/blob/main/docs/gobuild.png)
![chooseWethertoCompress](https://github.com/moqsien/gvc/blob/main/docs/compressOrNot.png)
![startBuild](https://github.com/moqsien/gvc/blob/main/docs/compiling.png)

### How to use NeoBox?
![](https://github.com/moqsien/gvc/blob/main/docs/neobox.png)

---------
## Something nice about [gvc](https://github.com/moqsien/gvc).

![logo](https://github.com/moqsien/gvc/blob/main/docs/logo.png)
---------
GVC is a nice tool designed for managing your development environment on multi-platforms and -machines.
It will help you to create a dev environment for Go, Python, Java, NodeJS or Rust, Cygwin, etc.
You do never need to worry about env variables or where to dowload files, GVC will do everything for you.
You can even install VSCode as well as Neovim through GVC.
Moreover, GVC will sync config files to WebDAV if you have configured one, which will help recreate your dev environment on another PC.

Therefore, If you are planning to get a new PC for development, probably the only thing you need to download is GVC!

### Note that: 

All download urls are customed for speeding-up in China by default. However, you can also <b>custom your own download urls</b> in gvc's config file at your convenience. The config file is likely to be <b>${Your-home-dir}/.gvc/backup/gvc-config.json</b>. Once you have found the file, you will know how to modify it.

Of course, you can also use <b>"g config show"</b> to show the config file path after the installation of gvc.

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
    <td><font color="LightBlue">Intall, Uninstall, SwitchVersion, SetEnv, SearchPackage, MorePorwerfulGoCompilationAndArchive</font></td>
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
    <td><font color="Gree">Typst</font></td>
    <td><font color="LightBlue">Install, SetEnv</font></td>
    <td bgcolor="LavenderBlush">gvc typst help</td>
  </tr>
  <tr>
    <td><font color="Gree">Cygwin</font></td>
    <td><font color="LightBlue">Install, InstallPackage</font></td>
    <td bgcolor="PaleVioletRed">gvc cpp ic help; Only for Windows; git,bash, clang, gcc, etc.</td>
  </tr>
  <tr>
    <td><font color="Gree">Msys2</font></td>
    <td><font color="LightBlue">Install</font></td>
    <td bgcolor="PaleVioletRed">gvc cpp im help; Only for Windows.</td>
  </tr>
  <tr>
    <td><font color="Gree">vcpkg</font></td>
    <td><font color="LightBlue">Install</font></td>
    <td bgcolor="PaleVioletRed">gvc cpp iv help; Cpp package management.</td>
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
    <td><font color="LightBlue">Auto get free Vmess/Vless/Trojan/Shadowsocks start listen at localhost:2019</font></td>
    <td bgcolor="PaleVioletRed">HelpInfo: gvc x help; EnterTheOperationShell: gvc x; <a href="https://github.com/moqsien/xtray">xtray docs</a></td>
  </tr>
  <tr>
    <td><font color="Gree">Browser</font></td>
    <td><font color="LightBlue">Auto backup browser data, such as bookmarks/password/extensionInfo. Your private info data, such as password is encrypted automatically.</font></td>
    <td bgcolor="PaleVioletRed">HelpInfo: gvc browser help</td>
  </tr>
  <tr>
    <td><font color="Gree">Github</font></td>
    <td><font color="LightBlue">Open github download acceleration website.</font></td>
    <td bgcolor="LavenderBlush">gvc github 1; gvc github 2</td>
  </tr>
</table>

### Download & Install
Download files, unarchive, then double clik or just run with no subcommand or argument, gvc will install itself to default dir.

- [Note] Make sure you have PowerShell installed under windows.
- [github release](https://github.com/moqsien/gvc/releases)
- [gitee release](https://gitee.com/moqsien/gvc_tools/releases)

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
   neobox-runner, nbrunner, nbr   Start a neobox client. # this is used by neobox-shell, users do not need care about.
   neobox-keeper, nbkeeper, nbk   Start a neobox keeper. # this is used by neobox-shell, users do not need care about.
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

### NeoBox Shell help info
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

### Subcommands Help Docs (in Zh-CN)
[github docs](https://github.com/moqsien/gvc/blob/main/docs/commands/command_list_github.md)

[gitee docs](https://gitee.com/moqsien/gvc_tools/blob/main/docs/commands/command_list_gitee.md)

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
