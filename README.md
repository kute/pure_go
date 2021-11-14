
安装库时在项目根目录执行go get -u -t -v library

1、环境设置：http://www.vangoleo.com/2019/11/02/go/go-hello-world-02/
2、go module：实现go文件不必非在GOPATH/src目录下也可以引用到
3、go标准库文档：https://studygolang.com/pkgdoc
4、go中文网：https://studygolang.com/
5、gvm: https://github.com/moovweb/gvm
6、go get
7、https://pkg.go.dev/

```
-d 只下载不安装
-f 只有在你包含了 -u 参数的时候才有效，不让 -u 去验证 import 中的每一个都已经获取了，这对于本地 fork 的包特别有用
-fix 在获取源码之后先运行 fix，然后再去做其他的事情
-t 同时也下载需要为运行测试所需要的包
-u 强制使用网络去更新包和它的依赖包
-v 显示执行的命令
-insecure	允许使用不安全的 HTTP 方式进行下载操作
```

