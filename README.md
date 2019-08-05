# Jrebel 激活服务器 By GoLang
- [实现原理和实现代码参考Java版本:https://gitee.com/gsls200808/JrebelLicenseServerforJava](https://gitee.com/gsls200808/JrebelLicenseServerforJava)

>Java开发者刚刚学GO,练手....

## Win编译
```shell
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build ./
```
- GOOS
    - linux
    - windows
    - darwin : 苹果系统
- GOARCH:
    - amd64 : 64位
    - 386:  : 32位