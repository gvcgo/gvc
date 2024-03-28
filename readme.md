### GVC是什么？

GVC是一个实用工具集合。它集成了一些好用的命令行工具, 支持MacOS/Windows/Linux。

### 子命令简介

```bash
homedir % g -h

asciinema   Records your terminal in asciinema cast form.
browser     Handles data from browser.
cloc        Counts lines of code.
git         Git related CLIs.
gopher      Some useful comand for gophers.
gpt         ChatGPT or FlyTek spark bot.
repo        Uses remote github/gitee repo as OSS.
```

**asciinema**: 终端session录制功能，支持编辑和上传，也支持转换为gif(通过version-manager安装agg后支持)后上传到github/gitee，对于写文档非常有用。

**browser**: 浏览器数据导出，数据一般存放在$HOME/.gvc/browser_data/目录下。

**cloc**: 项目代码行数统计，统计项目中使用的各种代码的类型、行数，注释行数，空行数。

**git**: 系统hosts文件一键更新，加速github访问(需要管理员权限，会自动备份旧的hosts文件)。为git ssh协议适配本地代理，加速github访问，可以一键切换有无代理模式。

**gopher**: go build命令增强；一键重命名go package；一键安装常用的go项目，例如grpc-go-gen、goctl、gf、dlv、gopls等，可以选择安装。

**gpt**: 一个基于TUI的ChatGPT/讯飞星火客户端。

**repo**: 1. vscode、asiinema、gpt、.ssh等配置文件的一键备份和还原，支持github/gitee仓库，敏感信息会自动加密；2、图片一键上传到github/gitee仓库，然后生成markdown可以引用的图片地址。

### 如何安装？

```bash
go install github.com/gvcgo/gvc/cmd/g@latest
```
