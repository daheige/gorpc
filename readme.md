# go grpc and http gateway项目实战

    1.通过pb文件生成客户端代码，同时支持grpc和http gw
    2.支持metrics prometheus，validator验证

# 基于gmicro框架
    
    github.com/daheige/gmicro

# 运行
    
    代码生成
    sh bin/go-generate.sh

    获取go package
    go mod tidy    
    go run cmd/rpc/main.go

    浏览器访问： http://localhost:1338/v1/say/23
    pprof: http://localhost:2358/debug/pprof/
    metrics: http://localhost:2358/metrics
    
    db,redis使用参考： github.com/daheige/goapp

# layout参考
    
    https://github.com/golang-standards/project-layout
