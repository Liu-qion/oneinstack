package storage

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"oneinstack/app"
	"oneinstack/internal/models"
	"oneinstack/internal/services"
	"oneinstack/web/input"
)

func Add(param *input.AddParam) error {
	s := &models.Storage{}
	tx := app.DB().Where("addr = ? and port =?", param.Addr, param.Port).First(s)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return tx.Error
	}
	if s.ID > 0 {
		return fmt.Errorf("%s:%v 已存在", param.Addr, param.Port)
	}
	m := &models.Storage{
		Addr:     param.Addr,
		Port:     param.Port,
		Root:     param.Root,
		Password: param.Password,
		Remark:   param.Remark,
		Type:     param.Type,
	}
	op, err := NewStorageOP(m)
	if err != nil {
		return err
	}
	err = op.Connet()
	if err != nil {
		return err
	}
	tx = app.DB().Create(m)
	return tx.Error
}

func List(ty string) ([]*models.Storage, error) {
	list := []*models.Storage{}
	tx := app.DB()
	if ty != "" && ty != "all" {
		tx = tx.Where("type = ?", ty)
	}
	tx = tx.Find(&list)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, tx.Error
	}
	return list, nil
}

func LibList(param *input.QueryParam) (*services.PaginatedResult[models.Library], error) {
	if param.Type == "redis" {
		s := &models.Storage{}
		tx := app.DB().Where("id = ?", param.ID).First(s)
		if tx.Error != nil {
			return nil, tx.Error
		}
		op := NewRedisOP(s)
		err := op.Connet()
		if err != nil {
			return nil, err
		}
		libs, err := op.GetLibs()
		if err != nil {
			return nil, err
		}
		return &services.PaginatedResult[models.Library]{
			Data:       libs,
			Total:      len(libs),
			Page:       1,
			PageSize:   100,
			TotalPages: 1,
		}, nil
	} else {
		return services.Paginate[models.Library](app.DB().Where("type = ?", param.Type), &models.Library{}, &input.Page{
			Page:     param.PageSize,
			PageSize: param.PageSize,
		})
	}
}

func Del(param *input.IDParam) error {
	tx := app.DB().Delete(&models.Storage{}, param.ID)
	return tx.Error
}

func Sync(param *input.IDParam) error {
	m := &models.Storage{}
	tx := app.DB().Where("id = ?", param.ID).First(m)
	if tx.Error != nil {
		return tx.Error
	}
	op, err := NewStorageOP(m)
	if err != nil {
		return err
	}
	err = op.Connet()
	if err != nil {
		return err
	}
	return op.Sync()
}

func Update(param *input.AddParam) error {
	s := &models.Storage{}
	tx := app.DB().Where("id = ?", param.ID).First(s)
	if tx.Error != nil {
		return tx.Error
	}
	s.Addr = param.Addr
	s.Port = param.Port
	s.Root = param.Root
	s.Password = param.Password
	s.Remark = param.Remark
	s.Remark = param.Remark

	op, err := NewStorageOP(s)
	if err != nil {
		return err
	}
	err = op.Connet()
	if err != nil {
		return err
	}
	tx = app.DB().Updates(s)
	return tx.Error
}

func RedisKeyList(param *input.QueryParam) (*PaginatedKeysInfo, error) {
	s := &models.Storage{}
	tx := app.DB().Where("id = ?", param.ID).First(s)
	if tx.Error != nil {
		return nil, tx.Error
	}
	op := NewRedisOP(s)
	err := op.Connet()
	if err != nil {
		return nil, err
	}
	return op.GetPaginatedKeyInfo(context.Background(), param.RDB, "", param.Page.Page, param.PageSize)

}
