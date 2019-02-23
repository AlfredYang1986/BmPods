# BmPods

## 使用方案
- 下载所以的依赖

```
go get github.com/alfredyang1986/blackmirror

go get github.com/alfredyang1986/BmServiceDef

go get github.com/alfredyang1986/BmPods
```

- 下载配置文件

```
go get github.com/alfredyang1986/BmServiceDeploy
```

- 配置环境变量BM_HOME

```
export BM_HOME=/path/to/you/BmServiceDeply
```

- 执行main.go
```
go run main.go
```

## 发布
- 下载docker文件
```
go get github.com/alfredyang1986/BmServiceDeploy
```

- 编译docker file
```
docker built . 
```

- docker 发布
