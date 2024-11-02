package distributor

import (
	"challenge2016/internal/model"
	"gofr.dev/pkg/gofr"
)

func (h *handler) Add(ctx *gofr.Context) (interface{}, error) {
	var distributor = &model.Distributor{}

	err := ctx.Bind(&distributor)
	if err != nil {
		ctx.Logger.Errorf("Bind error: %v", err.Error())
		return nil, err
	}

	return h.Service.Add(distributor)
}

func (h *handler) Get(ctx *gofr.Context) (interface{}, error) {
	var name = model.Sanitize(ctx.PathParam("name"))

	return h.Service.Get(name)
}

func (h *handler) Check(ctx *gofr.Context) (interface{}, error) {
	var (
		name   = model.Sanitize(ctx.PathParam("name"))
		region = model.Sanitize(ctx.Param("region"))
	)

	return h.Service.CheckAccess(name, region)
}
