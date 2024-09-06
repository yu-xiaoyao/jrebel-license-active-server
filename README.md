# Jrebel 激活服务器 By GoLang


## 直接下载(Win,Linux,Mac)

- [https://github.com/yu-xiaoyao/jrebel-license-active-server/releases](https://github.com/yu-xiaoyao/jrebel-license-active-server/releases)


## Test Server
- http://117.50.194.13:12345

Example:
```shell
http://117.50.194.13:12345/524f1d03-d1d8-5e94-a099-042736d40bd9
```

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
## Mac编译
```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./
```

## docker
```shell
docker run --rm --env GOPROXY=https://goproxy.cn -v "$PWD":/root -w /root/src/project/main golang:latest go build ./ -v 
```

## 运行
默认端口: 12345
```shell
# 自定义端口
./license-active-server --port=5555
```