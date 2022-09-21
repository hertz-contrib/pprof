# pprof (本项目由CloudWeGo社区贡献者协作开发并维护)

[English](README.md) | 中文

pprof 是为 Hertz 框架开发的中间件，参考了 [Gin](https://github.com/gin-gonic/gin) 中 [pprof](https://github.com/gin-contrib/pprof) 的实现。
本项目的完成得益于 CloudWeGo 社区的工作以及 Gin 社区所做的相关前置工作。


## 安装
```shell
go get github.com/hertz-contrib/pprof
```

## 使用
### 代码实例1

```go
func main() {
    h := server.Default()
    
    pprof.Register(h)
    
    h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
    ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
    })
    
    h.Spin()
}
```

### 代码实例2: 自定义前缀

```go
import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/common/utils"
    "github.com/cloudwego/hertz/pkg/protocol/consts"
    "hertz-contrib-beiye/pprof"
)

func main() {
    h := server.Default()
    
    // default is "debug/pprof"
    pprof.Register(h, "dev/pprof")
    
    h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
        ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
    })
    
    h.Spin()
}

```

### 代码实例3: 自定义路由组

```go
import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    "hertz-contrib-beiye/pprof"
    "net/http"
)

func main() {
    h := server.Default()
    
    pprof.Register(h)

    adminGroup := h.Group("/admin")
    adminGroup.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
    ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
    })
    
    pprof.RouteRegister(adminGroup, "pprof")
    
    h.Spin()
}

```


### 如何使用

使用 `pprof tool` 工具查看堆栈采样信息：

```bash
go tool pprof http://localhost:8888/debug/pprof/heap
```

使用 `pprof tool` 工具查看30s的CPU采样信息：

```bash
go tool pprof http://localhost:8888/debug/pprof/profile
```

使用 `pprof tool` 工具查看go 协程阻塞信息：

```bash
go tool pprof http://localhost:8888/debug/pprof/block
```

使用 `pprof tool` 工具查看5s内的执行trace：

```bash
wget http://localhost:8888/debug/pprof/trace?seconds=5
```


## License
This project is under the Apache License 2.0. See the LICENSE file for the full license text.