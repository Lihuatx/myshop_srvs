package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/hashicorp/consul/api"
	"myshop_srvs/goods_srv/global"
	"myshop_srvs/goods_srv/handler"
	"myshop_srvs/goods_srv/initialize"
	"myshop_srvs/goods_srv/proto"
	"myshop_srvs/goods_srv/utils"
)

// main 函数是程序的入口点，用于初始化服务并启动 goods 服务的 gRPC 服务器。
//
// 该函数首先定义了两个命令行参数：ip 和 port，分别用于指定服务器的 IP 地址和端口号。
//
// 接下来，调用 initialize 包中的函数来初始化日志记录器、配置文件、数据库和 Elasticsearch 连接。
//
// 然后，解析命令行参数，并根据参数值打印出 IP 地址和端口号。如果未指定端口号，则调用 utils 包中的 GetFreePort 函数来获取一个空闲端口。
//
// 创建一个 gRPC 服务器实例，并注册 GoodsServer 服务。接着，监听指定的 IP 地址和端口，等待客户端连接。
//
// 注册服务健康检查，以便 Consul 可以监控服务的健康状态。
//
// 配置 Consul 客户端，并创建一个服务注册对象。该对象包含服务的名称、ID、端口、标签和地址等信息，以及一个健康检查对象。
//
// 使用 Consul 客户端将服务注册到 Consul 代理中。如果注册失败，则程序将崩溃。
//
// 在一个 goroutine 中启动 gRPC 服务器，以便主线程可以继续执行其他操作。
//
// 设置一个通道来接收终止信号（如 SIGINT 或 SIGTERM），并在接收到信号时注销服务并关闭服务器。
func main() {
	// 定义命令行参数
	IP := flag.String("ip", "0.0.0.0", "指定服务器的 IP 地址")
	Port := flag.Int("port", 0, "指定服务器的端口号")

	// 初始化
	initialize.InitLogger()           // 初始化日志记录器
	initialize.InitConfig()           // 初始化配置文件
	initialize.InitDB()               // 初始化数据库连接
	initialize.InitEs()               // 初始化 Elasticsearch 连接
	zap.S().Info(global.ServerConfig) // 打印服务器配置信息

	// 解析命令行参数
	flag.Parse()
	zap.S().Info("ip: ", *IP) // 打印 IP 地址
	if *Port == 0 {
		*Port, _ = utils.GetFreePort() // 获取一个空闲端口
	}
	zap.S().Info("port: ", *Port) // 打印端口号

	// 创建 gRPC 服务器实例
	server := grpc.NewServer()
	proto.RegisterGoodsServer(server, &handler.GoodsServer{}) // 注册 GoodsServer 服务

	// 监听指定的 IP 地址和端口
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error()) // 监听失败则崩溃
	}

	// 注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 配置 Consul 客户端并注册服务
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err) // Consul 客户端创建失败则崩溃
	}

	// 创建服务注册对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = global.ServerConfig.Tags
	registration.Address = global.ServerConfig.Host
	registration.Check = check

	// 注册服务到 Consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err) // 服务注册失败则崩溃
	}

	// 启动 gRPC 服务器
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error()) // gRPC 服务器启动失败则崩溃
		}
	}()

	// 接收终止信号并注销服务
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败") // 注销服务失败则记录日志
	}
	zap.S().Info("注销成功") // 记录注销成功日志
}
