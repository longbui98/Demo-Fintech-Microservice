package main

import (
	"flag"
	"fmt"
	"micro-project/order_service/internal"
	"micro-project/order_service/internal/config"
	"micro-project/order_service/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

func main() {
	var configFile = flag.String("f", "etc/user.yaml", "the config file")
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(rest.RestConf{
		Host: c.Host,
		Port: c.Port,
	})
	defer server.Stop()

	internal.RegisterRoutes(server, ctx)

	fmt.Printf("User Service running on %s:%d\n", c.Host, c.Port)
	server.Start()
}
