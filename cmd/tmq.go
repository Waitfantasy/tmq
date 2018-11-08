package main

import (
	"flag"
	"fmt"
	"github.com/Waitfantasy/tmq/config"
	tmqHttp "github.com/Waitfantasy/tmq/http"
	"github.com/Waitfantasy/tmq/logger"
	"github.com/Waitfantasy/tmq/message/manager"
	"github.com/Waitfantasy/tmq/message/persistent"
	"github.com/Waitfantasy/tmq/mq"
	tmpRpc "github.com/Waitfantasy/tmq/rpc"
	"github.com/Waitfantasy/unicorn/rpc"
	"github.com/Waitfantasy/unicorn/rpc/client"
	"runtime"
	"strings"
	"time"
)

const banner string = `
 ________  _______ 
/_  __/  |/  / __ \
 / / / /|_/ / /_/ /
/_/ /_/  /_/\___\_\
`

var configFile *string = flag.String("config", "/etc/tmq.yaml", "tmq config file")

func main() {
	fmt.Print(banner)
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	if len(*configFile) == 0 {
		fmt.Println("must use a config file")
		return
	}

	var (
		err          error
		cfg          *config.Config
		mqer         mq.Mqer
		rpcSrv       *tmpRpc.MqServer
		httpSrv      *tmqHttp.MqServer
		mg           *manager.Manager
		persistenter persistent.Persistenter
	)

	if cfg, err = config.ParseConfigFile(*configFile); err != nil {
		fmt.Printf("parse config file error:%v\n", err.Error())
		return
	}

	// new persistent service
	switch strings.ToLower(cfg.Persistent) {
	case persistent.MysqlPersistentType:
		if persistenter, err = persistent.NewMysqlPersistent(cfg.MysqlDSN); err != nil {
			fmt.Printf("persistent.NewMysqlPersistent(%s) create mysql persistent type error: %v\n", cfg.MysqlDSN, err.Error())
		}
	case persistent.RedisPersistentType:
		persistenter = persistent.NewRedisPersistent(cfg)
	case persistent.MongoPersistentType:
		panic("implement me")
	default:
		fmt.Printf("unsupport persistent type\n")
		return
	}

	// new mq service
	switch strings.ToLower(cfg.Mq) {
	case mq.NsqType:
		mqer = mq.NewNsqMq(cfg.Nsq.Addr, time.Duration(cfg.Nsq.Delay))
	case mq.RedisType:
		mqer = mq.NewRedisMq(cfg)
	default:
		fmt.Printf("unsupport mq type\n")
		return
	}

	// new manager
	rpcCli, err := client.New(&rpc.Config{
		Addr:       cfg.IdService.Addr,
		EnableTLS:  cfg.IdService.EnableTLS,
		CertFile:   cfg.IdService.CertFile,
		KeyFile:    cfg.IdService.KeyFile,
		ServerName: cfg.IdService.ServerName,
	})

	if err != nil {
		fmt.Printf("new id rpc client error: %v\n", err.Error())
		return
	}

	// new logger
	log := logger.New(&logger.Config{
		Level:      cfg.Logger.Level,
		Output:     cfg.Logger.Output,
		Split:      cfg.Logger.Split,
		FilePath:   cfg.Logger.FilePath,
		FilePrefix: cfg.Logger.FilePrefix,
		FileSuffix: cfg.Logger.FileSuffix,
	})

	if err = log.Run(); err != nil {
		fmt.Printf("log run error: %v\n", err.Error())
		return
	}

	// new manager
	mg = manager.New(mqer, persistenter, rpcCli)

	rpcSrv = tmpRpc.New(&tmpRpc.Config{
		Addr:       cfg.GRpc.Addr,
		EnableTLS:  cfg.GRpc.EnableTLS,
		CertFile:   cfg.GRpc.CertFile,
		KeyFile:    cfg.GRpc.KeyFile,
		ServerName: cfg.GRpc.ServerName,
	}, mg)

	// start rpc server
	go func() {
		log.Debug("tmq grpc server listen in: %s", cfg.GRpc.Addr)
		rpcSrv.Run()
	}()

	// start http server
	httpSrv = tmqHttp.New(&tmqHttp.Config{
		Addr:       cfg.Http.Addr,
		EnableTLS:  cfg.Http.EnableTLS,
		CaFile:     cfg.Http.CaFile,
		CertFile:   cfg.Http.CertFile,
		KeyFile:    cfg.Http.KeyFile,
		ClientAuth: cfg.Http.ClientAuth,
	}, mg)

	log.Debug("tmq http server listen in: %s", cfg.Http.Addr)
	httpSrv.Run()
}
