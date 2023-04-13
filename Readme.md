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
## Something nice about [gvc](https://github.com/moqsien/gvc).
---------
GVC is a nice tool designed for managing your development environment on multi-platforms and -machines.
It will help you to create a dev environment for Go, Python, Java, NodeJS or Rust, Cygwin, etc.
You do never need to worry about env variables or where to dowload files, GVC will do everything for you.
You can even install VSCode as well as Neovim through GVC.
Moreover, GVC will sync config files to WebDAV if you have configured one, which will help recreate your dev environment on another PC.

Therefore, If you are planning to get a new PC for development, probably the only thing you need to download is GVC!

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
    <td><font color="LightBlue">Intall, Uninstall, SwitchVersion, SetEnv, SearchPackage</font></td>
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
    <td bgcolor="PaleVioletRed">gvc cygwin help; Only for Windows; git,bash, clang, gcc, etc.</td>
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
    <td><font color="LightBlue">Auto get free Vmess start listen at localhost:2019</font></td>
    <td bgcolor="PaleVioletRed">HelpInfo: gvc xray help; EnterTheOperationShell: gvc xray</td>
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

### Subcommands Help Docs (in Zh-CN)
[github docs](https://github.com/moqsien/gvc/blob/main/docs/commands/command_list_github.md)

[gitee docs](https://gitee.com/moqsien/gvc_tools/blob/main/docs/commands/command_list_gitee.md)

## thanks to
---------
- [xray-core](https://github.com/XTLS/Xray-core)
- [pyenv](https://github.com/pyenv/pyenv)
- [pyenv-win](https://github.com/pyenv-win/pyenv-win)
- [g](https://github.com/voidint/g)
- [gvm](https://github.com/andrewkroh/gvm)
