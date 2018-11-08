package service

import (
	"github.com/Waitfantasy/tmq/config"
	"github.com/Waitfantasy/unicorn/rpc"
	"github.com/Waitfantasy/unicorn/rpc/client"
)

func GetId(cfg *config.Config) {
	client.New(&rpc.Config{
		Addr: cfg.GRpc.Addr,
	})
}
