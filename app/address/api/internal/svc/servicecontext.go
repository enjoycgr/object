package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"object/app/address/api/internal/config"
	"object/app/address/api/internal/model"
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
