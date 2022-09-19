# pprof (This is a community driven project)


pprof middleware for Hertz framework, inspired by [pprof](https://github.com/gin-contrib/pprof).
This project would not have been possible without the support from the CloudWeGo community and previous work done by the gin community.

- Package pprof serves via its HTTP server runtime profiling data in the format expected by the pprof visualization tool.


## Install
```shell
go get github.com/hertz-contrib/pprof
```

## Usage
### Example

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

### change default path prefix

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

### custom router group

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


### Use the pprof tool

Then use the pprof tool to look at the heap profile:

```bash
go tool pprof http://localhost:8888/debug/pprof/heap
```

Or to look at a 30-second CPU profile:

```bash
go tool pprof http://localhost:8888/debug/pprof/profile
```

Or to look at the goroutine blocking profile, after calling runtime.SetBlockProfileRate in your program:

```bash
go tool pprof http://localhost:8888/debug/pprof/block
```

Or to collect a 5-second execution trace:

```bash
wget http://localhost:8888/debug/pprof/trace?seconds=5
```


## License
This project is under the Apache License 2.0. See the LICENSE file for the full license text.