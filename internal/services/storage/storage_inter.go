package storage

import (
	"fmt"
	"oneinstack/internal/models"
)

type StorageOP struct {
	DB   StorageOPI
	Type string
}

type StorageOPI interface {
	Connet() error
	Sync() error
}

func NewStorageOP(p *models.Storage) (StorageOPI, error) {
	switch p.Type {
	case "mysql":
		return NewMysqlOP(p), nil
	case "pg":
	case "sqlserver":
	case "redis":
		return NewRedisOP(p), nil
	case "mongo":
	}
	return nil, fmt.Errorf("未知的存储服务")
}

func (s StorageOP) Connet() error {
	return s.DB.Connet()
}

func (s StorageOP) Sync() error {
	return s.DB.Sync()
}
