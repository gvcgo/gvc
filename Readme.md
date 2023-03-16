## [中文](https://github.com/moqsien/gvc/blob/main/docs/Readme_CN.md)
---------

## Something nice about [gvc](https://github.com/moqsien/gvc).
---------
- If only you don't need to work overtime every single day.
- If only the code you maintained won't become a crap.
- If only you really enjoy coding.
- If only you will never be forced to accept silly dev requirements.
- If only you could build your dev environments easily.
- If only you could get it work and enjoy studying as a newbie.
- If only you could work more efficiently without meaningless "involutions"(内卷).
### Therefore, gvc was born!
At Present, gvc has following features:
- Automatically installation of golang compilers, version switching and env variables setup;
- Automatically installation of jdk, version switching and env variables setup;
- Automatically installation of the latest rust compiler, env variables setup;
- Automatically installation of nodejs, version switching and env variables setup;
- Automatically installation of vscode and extensions. You can also backup extensions info, user settings and keybindings;
- Automatically installation of neovim with default init file from gvc. And It was made to work with vscode by default;
- Hosts file management to accelerate github visit in China；
- Nearly all downloadings are accelerated for users in China；
- Highly configurable. You can configure your fast downloading url in gvc-config.yml;
- Sync config files to webdav if you have setup one;
- Supported Platform: MacOS, Windows, Linux(untested at present);

Fetures on the way：
- HomeBrew auto-installation and acceleration for users in China；
- Python acceleration for users in China；
- git.exe auto-installation under Windows；
- Flutter auto-installation；

## gvc Help Info
---------
### gvc -h
```shell
moqsien@iMac gvc % gvc -h
NAME:
   gvc - gvc <Command> <SubCommand>...

USAGE:
   gvc [global options] command [command options] [arguments...]

DESCRIPTION:
   A productive tool to manage your development environment.

COMMANDS:
   host, h, hosts        Manage system hosts file.
   go, g                 Go version control.
   vscode, vsc, vs, v    VSCode management.
   config, conf, cnf, c  GVC config file management.
   nvim, neovim, nv, n   GVC neovim management.
   java, jdk, j          GVC jdk management.
   rust, rustc, ru, r    GVC rust management.
   nodejs, node, no      Nodejs version control.
   help, h               Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```
### gvc go -h
```shell
moqsien@iMac-Pro gvc % gvc go help
NAME:
   gvc go - Go version control.

USAGE:
   gvc go command [command options] [arguments...]

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

### gvc vs -h
```shell
moqsien@iMac-Pro gvc % gvc vs -h
NAME:
   gvc vscode - VSCode management.

USAGE:
   gvc vscode command [command options] [arguments...]

COMMANDS:
   install, i, ins               Automatically install vscode.
   install-extensions, ie, iext  Automatically install extensions for vscode.
   sync-extensions, se, sext     Push local installed vscode extensions info to remote webdav.
   get-settings, gs, gset        Get vscode settings(keybindings include) info from remote webdav.
   push-settings, ps, pset       Push vscode settings(keybindings include) info to remote webdav.
   help, h                       Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc nv -h
```shell
moqsien@iMac-Pro gvc % gvc nv -h
NAME:
   gvc nvim - GVC neovim management.

USAGE:
   gvc nvim command [command options] [arguments...]

COMMANDS:
   install, ins, i  Install neovim.
   help, h          Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc cnf -h
```shell
moqsien@iMac-Pro gvc % gvc cnf -h
NAME:
   gvc config - GVC config file management.

USAGE:
   gvc config command [command options] [arguments...]

COMMANDS:
   webdav, dav, w   Setup webdav account info to backup local settings for gvc, vscode, neovim etc.
   pull, pl         Pull settings to local backup dir from your remote webdav.
   push, ph         Push settings from local backup dir to your remote webdav.
   show, sh, s      Show path to conf files.
   reset, rs, r     Reset config file to default values.
   download, dl, d  Download example config files from gitee when backup dir is empty.
   help, h          Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc host -h
```shell
moqsien@iMac-Pro gvc % gvc host -h
NAME:
   gvc host - Manage system hosts file.

USAGE:
   gvc host command [command options] [arguments...]

COMMANDS:
   fetch, f      Fetch github hosts info.
   fetchall, fa  Get all github hosts info with no ping filters.
   show, s       Show hosts file path.
   help, h       Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc java -h
```shell
moqsien@iMac gvc % gvc java -h
NAME:
   gvc java - GVC jdk management.

USAGE:
   gvc java command [command options] [arguments...]

COMMANDS:
   use, u   Download and use jdk.
   show, s  Show available versions.
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc rust -h
```shell
moqsien@iMac gvc % gvc rust -h
NAME:
   gvc rust - GVC rust management.

USAGE:
   gvc rust command [command options] [arguments...]

COMMANDS:
   install, ins, i  Install the latest rust compiler tools.
   help, h          Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### gvc node -h
```shell
moqsien@iMac gvc % gvc node -h
NAME:
   gvc nodejs - Nodejs version control.

USAGE:
   gvc nodejs command [command options] [arguments...]

COMMANDS:
   remote, r           Show remote versions.
   use, u              Download and use version.
   local, l            Show installed versions.
   remove-unused, ru   Remove unused versions.
   remove-version, rm  Remove a version.
   help, h             Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### Download & Install
Download files, unarchive, then double clik or just run with no subcommand or argument, gvc will install itself to default dir.
- [github release](https://github.com/moqsien/gvc/releases)
- [gitee release](https://gitee.com/moqsien/gvc/releases/tag/v2)
