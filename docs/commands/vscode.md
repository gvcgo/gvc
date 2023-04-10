## gvc vscode -h
```shell
NAME:
   g vscode - VSCode and extensions installation.

USAGE:
   g vscode command [command options] [arguments...]

COMMANDS:
   install, i, ins               Automatically install vscode.
   install-extensions, ie, iext  Automatically install extensions for vscode.
   help, h                       Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

------------

### gvc vscode install
- 从官网获取并通过国内cnd加速下载vscode，然后安装，windows下会自动创建桌面快捷方式；
- Linux下需要手动创建快捷方式或者固定到任务栏

### gvc vscode iext
- 自动安装vscode插件
- 如果有从webdav同步插件id列表，则按照id列表来安装
- 否则，使用gvc-config.yml默认的常用插件id列表
