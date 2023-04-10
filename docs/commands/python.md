## gvc py -h
```shell
NAME:
   g python - Python version management.

USAGE:
   g python command [command options] [arguments...]

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

------------

### gvc py path
- 显示pyenv或者pyenv-win(Windows下)的安装目录

### gvc py update
- 更新pyenv或pyenv-win(注意pyenv-win强制更新后，会覆盖掉之前已经安装的所有python版本，需要重新安装)

### 其他子命令功能基本与go和jdk类似
