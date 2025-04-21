# Jrebel Active Server

## Download 下载 (Win,Linux,Mac)

- [https://github.com/yu-xiaoyao/jrebel-license-active-server/releases](https://github.com/yu-xiaoyao/jrebel-license-active-server/releases)

## Test Server
> 请手动添加 `http://`.

- 117.50.194.13:12345

Example:

```shell
117.50.194.13:12345/524f1d03-d1d8-5e94-a099-042736d40bd9
```

## 编译

- GOOS
  - linux
  - windows
  - darwin : 苹果系统
- GOARCH:
  - amd64 : 64位
  - 386:  : 32位

### Win编译
**Cmd**:
```shell
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build ./
```

### Mac编译

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./
```

### 运行

默认端口: 12345

```shell
# 自定义端口
./license-active-server --port=5555
# index page show as https
./license-active-server --port=5555 --exportSchema=https --exportHost=jrebel.domain.com
```


## docker

### Pull Image
```shell
docker pull yuxiaoyao520/jrebel-license-active-server:latest
```

### Run
```shell
docker run -p 12345:12345 --name jrebel-license-active-server yuxiaoyao520/jrebel-license-active-server:latest
```

### docker-compose

#### Simple
**docker-compose.yml**
```yaml
services:
  jrebel-license-active-server:
    image: yuxiaoyao520/jrebel-license-active-server:latest
    container_name: jrebel-license-active-server
    ports:
      - "12345:12345"
```

#### Add Run Args
**docker-compose.yml**
```yaml
services:
  jrebel-license-active-server:
    image: yuxiaoyao520/jrebel-license-active-server:latest
    container_name: jrebel-license-active-server
    command: ./jrebel-license-active-server --port=5555 --exportSchema=https --exportHost=jrebel.domain.com
    ports:
      - "5555:5555"
```

### Custom Build Image

```shell
docker build -t jrebel-license-active-server .
```

### Test Run

```shell
docker run --rm -p 12345:12345 jrebel-license-active-server:latest
```

