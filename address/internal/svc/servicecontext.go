package svc

import (
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"go-zero-demo/address/internal/config"
	"go-zero-demo/address/internal/model"
)

type ServiceContext struct {
	Config config.Config

	AddressModel model.AddressModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		AddressModel: model.NewAddressModel(sqlx.NewMysql(c.Mysql.DataSource), c.Cache),
	}
}
