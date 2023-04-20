# pprof (本项目由CloudWeGo社区贡献者协作开发并维护)

[English](README.md) | 中文

pprof 是为 Hertz 框架开发的中间件，参考了 [Gin](https://github.com/gin-gonic/gin) 中 [pprof](https://github.com/gin-contrib/pprof) 的实现。
本项目的完成得益于 CloudWeGo 社区的工作以及 Gin 社区所做的相关前置工作。

fgprof 部分参考了 [fgprof](https://github.com/felixge/fgprof)的实现。
如果使用了fgprof，请升级到Go 1.19或更高版本。在旧版本的Go中，对于具有大量goroutine（>1-10k）的应用程序，fgprof可能会导致显著的STW延迟。有关更多详细信息，请参见 [CL 387415](https://go-review.googlesource.com/c/go/+/387415)。

## 安装
```shell
go get github.com/hertz-contrib/pprof
```

## 使用
### pprof 代码实例1

```go
package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/pprof"
)

func main() {
	h := server.Default()

	pprof.Register(h)

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
	})

	h.Spin()
}
```

### pprof 代码实例2: 自定义前缀

```go
package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/pprof"
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

### pprof 代码实例3: 自定义路由组

```go
package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/pprof"
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

### fgprof 代码实例1
```go
package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/pprof"
)

	h := server.Default()

	pprof.FgprofRegister(h)

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
	})

	h.Spin()
```

### fgprof 代码实例2: 自定义前缀

```go
package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/pprof"
)

func main() {
	h := server.Default()

	// default is "debug/pprof"
	pprof.FgprofRegister(h, "dev/fgprof")

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
	})

	h.Spin()
}

```

### fgprof 代码实例3：自定义路由组

```go
package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/pprof"
)

func main() {
	h := server.Default()

	pprof.FgprofRegister(h)

	adminGroup := h.Group("/admin")

	adminGroup.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"ping": "pong"})
	})

	pprof.FgprofRouteRegister(adminGroup, "fgprof")

	h.Spin()
}
```
---

### 如何使用 pprof

使用 `pprof tool` 工具查看堆栈采样信息：

```bash
go tool pprof http://localhost:8888/debug/pprof/heap
```

使用 `pprof tool` 工具查看 30s 的 CPU 采样信息：

```bash
go tool pprof http://localhost:8888/debug/pprof/profile
```

使用 `pprof tool` 工具查看 go 协程阻塞信息：

```bash
go tool pprof http://localhost:8888/debug/pprof/block
```

使用 `pprof tool` 工具查看 5s 内的执行 trace：

```bash
wget http://localhost:8888/debug/pprof/trace?seconds=5
```

### 如何使用 fgprof

使用 `pprof tool` 工具查看 3s 内的采样信息：

```
go tool pprof --http=:6061 http://localhost:8888/debug/fgprof?seconds=3
```

## License
This project is under the Apache License 2.0. See the LICENSE file for the full license text.