# BmPods

## 1. 环境搭建
### 1.1 配置好本机的GOROOT和GOPATH
如我的机器:
```angular2html
GOROOT=/usr/local/Cellar/go/1.10.3/libexec
GOPATH=$HOME/.go
```

### 1.2 配置编译器的GOROOT和GOPATH
不同编译器略有差别, 自行去Google百度


## 2. 运行项目
使用命令`go get`, 依赖将自动安装到`$GOPATH`环境下  

### 2.1 本项目依次安装以下依赖:
```angular2html
go get github.com/alfredyang1986/blackmirror
go get github.com/alfredyang1986/BmPods
go get github.com/alfredyang1986/BmServiceDef
// 上面语句如果报异常`can't load package`,则进入下载的$GO_PATH/src对应的目录下,运行`go install`

// go get github.com/go-mgo/mgo
// go get github.com/go-redis/redis
// go get github.com/go-yaml/yaml
// go get github.com/aliyun/aliyun-oss-go-sdk/oss
// go get github.com/hashicorp/go-uuid
```

### 2.2 下载配置文件
```
go get github.com/alfredyang1986/BmServiceDeploy
```

### 2.3 配置环境变量BM_HOME
```
export BM_HOME=/path/to/you/BmServiceDeply/deploy-config/
```

### 2.4 执行main.go
```
go run main.go
```

## 3. 发布
### 3.1 下载docker文件
```
go get github.com/alfredyang1986/BmServiceDeploy
```

### 3.2 编译docker file
```
docker built . 
```

### 3.3 docker 发布
