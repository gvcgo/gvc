## GVC

For more info, please see [gvc wiki](https://github.com/moqsien/gvc/wiki).


--------

## 关于GVC能做什么的说明

总的来说，GVC是一个命令行工具。主要用于帮助管理你的开发环境。

其中，很重要的一方面就是对各种编程语言的管理。GVC支持Go,Java,Python,NodeJS,Flutter(Dart),Julia语言的一键安装，多版本切换管理(非docker实现，更轻量)，环境变量自动配置，中国大陆下载加速等。此外，还支持C/C++(Cygwin/Msys2), Rust, Vlang, Zig, Typst等语言的一键安装和环境变量配置。通过几个简单的命令，就可以立即得到一个可用的开发环境。

对于Go和Flutter，还有一些额外的命令支持。例如，GVC可以用于无脚本跨平台编译和打包你的go项目，这对于发布跨平台可执行文件非常有用。GVC可以通过几个命令来配置单独使用VSCode作为Flutter开发的IDE，无需另外安装Android Studio，不用非常麻烦的手动安装各种SDK。

另一个重要的方面，GVC还有很强的git支持。在中国大陆，由于某些未知原因，github访问速度非常慢，甚至时常失联。GVC提供了一些可用的git加速方案，集成了lazygit(一个强大好用基于TUI的git客户端)。你可以通过GVC来设置ssh协议代理，并控制ssh代理的开启和关闭，这比git自身的代理配置方便很多。另外，GVC还提供了一些自带的git命令以及git命令的简单组合，这样你可以在没有安装git的系统上使用这个自带命令。总之，有了GVC，你可以快速访问github，这对与github重度用户应该很有用。

还有一个非常大的功能，就是NeoBox。它是一个基于Xray-core和Sing-box的免费梯子客户端。提供了非常好用的免费梯子和各种强大的功能。具体可以参考NeoBox的wiki。

GVC提供的另外一个好用功能是，你可以配置一个远程github或者gitee仓库作为远程存储。这个远程存储可以用于保存你的各种配置文件，比如GVC自身的配置文件，VSCode相关的配置文件，你的浏览器书签和密码列表，ChatGPT/讯飞星火相关的配置文件，Asciinema的配置文件，.ssh文件等等。有了这些配置文件，你可以方便地在任何机器上通过GVC恢复你熟悉的开发环境。当然，这些配置文件中涉及密码或token的隐私部分会进行加密之后再保存，密码是由你自己设置的。

最后，GVC还提供了诸如ChatGPT终端问答机器人、项目代码统计、Asciinema终端录制、使用Github/Gitee作为Markdown图床等等实用功能。

GVC的所有功能的描述都能在帮助信息中找到，相信你在逐步使用过程中会发现GVC的便捷和强大。
