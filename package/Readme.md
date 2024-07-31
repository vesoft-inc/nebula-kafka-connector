# package

目前工具不会开源，有些项目需要工具的源码来编译，或者使用 SDK 的源码。
写一个统一的 shell 用来打包。

## 主要功能

* 通用逻辑：

  * 去掉 .git，.github，.gitignore
  * 统一包名，如 ${product}-${version}.tar.gz

* 可以指定 excluding 目录
* (暂时不做，以后可以加) 在每个文件加上指定的 license 信息
