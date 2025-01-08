package models

import (
	"gorm.io/gorm"
	"time"
)

const (
	Soft_Status_Default = 0
	Soft_Status_Ing     = 1
	Soft_Status_Suc     = 2
	Soft_Status_Err     = 3
)

type Software struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Key        string    `json:"key"`
	Icon       string    `json:"icon"`
	Type       string    `json:"type"`
	Status     int       `json:"status"` //0待安装,1安装中,2安装成功,3安装失败
	Resource   string    `json:"resource"`
	Installed  bool      `json:"installed"`
	Tags       string    `json:"tags"`
	Version    string    `json:"version"`
	Params     string    `json:"params"`
	Log        string    `json:"log"`
	CreateTime time.Time `json:"create_time"`
}

func (m *Software) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	return
}
