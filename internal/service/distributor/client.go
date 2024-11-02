package distributor

import (
	"challenge2016/internal/model"
)

type client struct{}

type Distributor interface {
	Add(*model.Distributor) (*model.DistributorResponse, error)
	Get(string) (*model.DistributorResponse, error)
	CheckAccess(string, string) (string, error)
}

func New() Distributor {
	return client{}
}
