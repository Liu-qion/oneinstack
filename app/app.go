package app

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"oneinstack/config"
)

var (
	db         *gorm.DB
	ONE_CONFIG config.Server
	ONE_VIP    *viper.Viper
)

func DB() *gorm.DB {
	return db
}
