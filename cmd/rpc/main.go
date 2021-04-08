package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/daheige/gmicro"
	"github.com/daheige/gorpc/api/clients/go/pb"
	"github.com/daheige/gorpc/internal/services"
	"google.golang.org/grpc"
)

var sharePort int
var shutdownFunc func()

func init() {
	sharePort = 8081

	shutdownFunc = func() {
		fmt.Println("Server shutting down")
	}
}

// http://localhost:8081/v1/say/1
/**
% go run services.go
2021/02/26 22:58:52 Starting http services and grpc services listening on 8081
2021/02/26 22:59:17 exec begin
2021/02/26 22:59:17 client_ip: 127.0.0.1
2021/02/26 22:59:17 req data:  id:1
2021/02/26 22:59:17 exec end,cost time: 0 ms
*/

func main() {
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
		gmicro.WithShutdownTimeout(6*time.Second),
		gmicro.WithHandlerFromEndpoint(pb.RegisterGreeterServiceHandlerFromEndpoint),
		gmicro.WithLogger(gmicro.LoggerFunc(log.Printf)),
		gmicro.WithRequestAccess(true),
		gmicro.WithPrometheus(true),
		gmicro.WithGRPCServerOption(grpc.ConnectionTimeout(10*time.Second)),
		gmicro.WithGRPCNetwork("tcp"), // grpc services start network
		gmicro.WithStaticAccess(true), // enable static file access,if use http gw
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

	newRoute2 := gmicro.Route{
		Method:  "GET",
		Pattern: gmicro.PathPattern("info"),
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	}

	s.AddRoute(newRoute2)

	// you can start grpc services and http gateway on one port
	log.Fatalln(s.StartGRPCAndHTTPServer(sharePort))

	// you can also specify ports for grpc and http gw separately
	// log.Fatalln(s.Start(sharePort, 50051))

	// you can start services without http gateway
	// log.Fatalln(s.StartGRPCWithoutGateway(50051))
}
