## gvc go help
```shell
NAME:
   g go - Go version management.

USAGE:
   g go command [command options] [arguments...]

COMMANDS:
   remote, r                   Show remote versions.
   use, u                      Download and use version.
   local, l                    Show installed versions.
   remove-unused, ru           Remove unused versions.
   remove-version, rm          Remove a version.
   add-envs, env, e, ae        Add envs for go.
   search-package, sp, search  Search for third-party packages.
   help, h                     Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

----------

### gvc go remote help
```shell
NAME:
   g go remote - Show remote versions.

USAGE:
   g go remote [command options] [arguments...]

OPTIONS:
   --show-all, -a, --all  Show all remote versions. (default: false)
   --help, -h             show help
```
- -a选项会展示所有可以安装的版本，否则只展示当前官方推荐的稳定版本。

### gvc go env
- 自动配置诸如GOPROXY，GOROOT，GOBIN，GOPATH等环境变量，GOPATH默认为~/data/projects/go(想修改的话，可以在gvc的配置文件中修改)

### gvc go local
- 显示本地已安装的go版本有哪些

### gvc go use xxx
- 下载并使用版本为xxx的go编译器

### gvc go rm xxx
- 卸载版本为xxx的go编译器

### gvc go ru
- 一键卸载不用的go编译器

### gvc go search xxx

```shell
NAME:
   g go search-package - Search for third-party packages.

USAGE:
   g go search-package [command options] [arguments...]

OPTIONS:
   --package-name value, -n value, --name value  Name of the package.
   --order-by-time, -o, --ou                     Order by update time. (default: false)
   --help, -h                                    show help
```
- 关键字xxx搜索可用的第三方包
- -o选项表示是否按更新时间排序，否则按引用量排序。
  