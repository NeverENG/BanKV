package znet

import (
	"github.com/NeverENG/BanKV/internal/network/ziface"
)

type BaseRouter struct{}

var _ ziface.IRouter = &BaseRouter{}

func (B *BaseRouter) PreHandle(req ziface.IRequest) {}

func (B *BaseRouter) Handle(req ziface.IRequest) {}

func (B *BaseRouter) PostHandle(req ziface.IRequest) {}
