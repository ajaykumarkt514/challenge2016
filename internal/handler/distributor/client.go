package distributor

import (
	"challenge2016/internal/service/distributor"
	"gofr.dev/pkg/gofr"
)

type handler struct {
	Service distributor.Distributor
}

type Distributor interface {
	Add(ctx *gofr.Context) (interface{}, error)
	Check(ctx *gofr.Context) (interface{}, error)
	Get(ctx *gofr.Context) (interface{}, error)
}

func New(s distributor.Distributor) Distributor {
	return &handler{
		Service: s,
	}
}
