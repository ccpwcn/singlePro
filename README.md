# singlePro
Go语言1.16及以上版本原生支持嵌入静态资源使用示例
> Go语言从1.16版本（2021年2月16日发布）起，在语言层面原生支持静态资源嵌入，实现了某些动态语言的特性（动态变更执行代码和数据I/O），并且为版本发布和运维部署提供了很大便利，这是一个很重要的特性。以前我们都是借助市面上的一些第三方工具实现的，现在Go语言自身具备这样的特性（通过embed标准库），这是一个非常好的消息。

# 验证方法
1. 保持`static/index.html`文件内容中的第9行的h2文本为`embed mode`
2. 编译输出一个可执行文件并执行它，不带任何参数，此时会加载编译时已经内置进去的static目录中的内容
4. 在浏览器中访问 http://localhost:9100 即可访问index.html，此时你会在页面上看到`embed mode`
5. 修改`index.html`，将`embed mode`改为`live mode`，重新运行编译好的程序，带上live参数，此时你会发现网页内容已经变更为`live mode`了
6. 再重新运行可执行程序，不带任何参数，你会在页面上得到`embed mode`。

# 原理探究
其实也很简单：  
- 将index.html内容改为`embed mode`，编译完之后再改回`live mode`，此时，`embed mode`的index.html已经在可执行程序中了
- 再执行可执行程序时：
  - 有`live`命令行参数，加载磁盘文件，于是看到`live mode`
  - 没有`live`命令行参数，加载内嵌到程序中的文件，于是看到`embed mode`
- 这是Go 1.16中的一个很重要的功能，基于此，我觉得以下场景会十分有用：
  - 编译时自动打进去一个version.txt，集成产品版本，以前的时候是使用`-ldflags "-X 'package_name/global.AppVersion=0.0.1'"`的方式，现在这种，更优雅了
  - 在一些特定场所，程序需要输出自己的源码（好像是Quine），那么就更方便了
  - 集成一个LICENSE到软件中，此时的LICENSE是固化在软件实体中的，它能够实时地跟着软件版本走，对于一些商业化和版权相关的领域，很有用
  - 先将一段代码编码输出到外部资源中，然后再在程序中动态加载和解码这段代码并执行它，简单说：A程序将一段代码编码导出到swap.gob，B程序编译时将swap.gob集成进来然后在程序运行时执行swap.gob中的代码。这使得Go语言多多少少有了点动态脚本语言的特性，它能够实现B程序不必做出修改就改变程序代码的行为（扩展插件思想？），所以，这是一个令人兴奋的特性。
  - 集成一个文件的内容到一个go程序变量，成为这个变量的值，这对于要实现开发期动态发布期静态的需求，实现很方便了
  - 直觉是这个`embed`特性应该还有其他用途，以后想到了再补充
  
总得来说，我觉得这一次这个embed特性是十分有意义的，代码可以动态交换（它控制软件流程和逻辑），文件也可以动态交换（它为软件提供数据），所以，代码逻辑可以动态替换，数据也可以动态替换的Go语言，未来必将更加大放异彩。

# 注意事项
- `main.go`文件中的第11行`//go:embed static`看着是个注释，但是不可缺少，Go编译器是依据它来决定是否要将static目录编译进可执行的二进制程序中的。
- 据我个人理解，这种将静态资源集成到可执行的二进制程序中的做法，有以下特点：
  - 对于一些不太想公开的资源，在一定程度上可以做到加密（虽然想想办法也能解开，但是增加了解开的门槛）
  - 正常情况下，一个外部的静态资源要加载进来，操作虽然烦琐，但是可以封装成函数，到处使用，也不算太麻烦，所以它更多的考虑是性能问题，毕竟，封装的函数每次执行时都会触发一次磁盘I/O
  - 所以，静态资源编译进可执行程序，从理论上讲可以大幅提升这些资源加载和使用的效率（直接走内存？我还没研究源码，不太确定） 
