## gvc hosts -h
```shell
NAME:
   g host - Sytem hosts file management(need admistrator or root).

USAGE:
   g host command [command options] [arguments...]

COMMANDS:
   fetch, f      Fetch github hosts info.
   fetchall, fa  Get all github hosts info with no ping filters.
   show, s       Show hosts file path.
   help, h       Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

--------

### gvc hosts fa
- 不进行响应时间测试，直接使用所有获取到的host，响应快，推荐

### gvc hosts f
- 进行响应时间测试，只使用符合响应时间配置的host，较慢，也有可能因为权限问题而失败，不推荐

### gvc hosts show
- 显示系统的hosts文件的存放路径
