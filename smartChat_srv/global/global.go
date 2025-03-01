package global

import (
	"gorm.io/gorm"
	"myshop_srvs/smartChat_srv/config"
	"myshop_srvs/smartChat_srv/proto"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig

	OrderSrvClient proto.OrderClient
)
