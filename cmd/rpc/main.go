package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/daheige/gmicro"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/daheige/gorpc/api/clients/go/pb"
	"github.com/daheige/gorpc/config"
	"github.com/daheige/gorpc/internal/interceptor"
	"github.com/daheige/gorpc/internal/services"
	"github.com/daheige/tigago/gpprof"
	"github.com/daheige/tigago/logger"
	"github.com/daheige/tigago/monitor"
)

var (
	configDir string
)

func init() {
	flag.StringVar(&configDir, "config_dir", "./", "config dir")
	flag.Parse()

	// init config.
	err := config.InitConfig(configDir)
	if err != nil {
		log.Fatalf("init config err: %v", err)
	}

	// 日志文件设置
	if config.AppServerConf.LogDir == "" {
		config.AppServerConf.LogDir = "./logs"
	}

	// 添加prometheus性能监控指标
	prometheus.MustRegister(monitor.CpuTemp)
	prometheus.MustRegister(monitor.HdFailures)

	// 性能监控的端口port+1000,只能在内网访问
	httpMux := gpprof.New()

	// 添加prometheus metrics处理器
	httpMux.Handle("/metrics", promhttp.Handler())
	gpprof.Run(httpMux, config.AppServerConf.PProfPort)
}

func main() {
	defer config.CloseAllDatabase()

	log.Println("rpc start...")
	log.Println("server pid: ", os.Getppid())

	// add the /test endpoint
	route := gmicro.Route{
		Method:  "GET",
		Pattern: gmicro.PathPattern("test"),
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Write([]byte("Hello!"))
		},
	}

	// test Option func
	s := gmicro.NewService(
		gmicro.WithRouteOpt(route),
		gmicro.WithShutdownFunc(shutdownFunc),
		gmicro.WithPreShutdownDelay(2*time.Second),
		gmicro.WithShutdownTimeout(5*time.Second),
		gmicro.WithHandlerFromEndpoint(pb.RegisterGreeterServiceHandlerFromEndpoint),
		gmicro.WithLogger(gmicro.LoggerFunc(log.Printf)),
		// gmicro.WithLogger(gmicro.LoggerFunc(gRPCPrintf)), // 定义grpc logger printf
		// gmicro.WithRequestAccess(true),
		gmicro.WithPrometheus(true),
		gmicro.WithGRPCServerOption(grpc.ConnectionTimeout(10*time.Second)),
		gmicro.WithUnaryInterceptor(interceptor.AccessLog), // 自定义访问日志记录
		gmicro.WithGRPCNetwork("tcp"),
	)

	// register grpc service
	pb.RegisterGreeterServiceServer(s.GRPCServer, &services.GreeterService{})

	newRoute := gmicro.Route{
		Method:  "GET",
		Pattern: gmicro.PathPattern("health"),
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	}

	s.AddRoute(newRoute)

	// log.Fatalln(s.StartGRPCAndHTTPServer(config.AppServerConf.GRPCPort))

	// run grpc and http gateway
	log.Fatalln(s.Start(config.AppServerConf.HttpPort, config.AppServerConf.GRPCPort))
}

func shutdownFunc() {
	log.Println("server will shutdown")
	logger.Info(context.Background(), "server will shutdown", nil)
}

// gmicro logger printf打印日志函数
func gRPCPrintf(format string, v ...interface{}) {
	logger.Info(context.Background(), fmt.Sprintf(format, v...), nil)
}
