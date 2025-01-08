package software

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"oneinstack/app"
	"oneinstack/internal/models"
	"oneinstack/web/input"
	"oneinstack/web/output"
	"strings"
)

func RunInstall(p *input.InstallParams) (string, error) {
	op, err := NewInstallOP(p)
	if err != nil {
		return "", err
	}
	return op.Install()
}

func List(param *input.SoftwareParam) ([]*output.Software, error) {
	tx := app.DB()
	if param.Id > 0 {
		tx = tx.Where("id = ?", param.Id)
	}

	if param.Name != "" {
		searchName := "%" + param.Name + "%"
		tx = tx.Where("name LIKE ?", searchName)
	}

	if param.Key != "" {
		searchKey := "%" + param.Key + "%"
		tx = tx.Where("key LIKE ?", searchKey)
	}

	if param.Type != "" {
		tx = tx.Where("type = ?", param.Type)
	}

	if param.Status != "" {
		tx = tx.Where("status = ?", param.Status)
	}

	if param.Resource != "" {
		tx = tx.Where("resource = ?", param.Resource)
	}

	if param.Installed {
		tx = tx.Where("installed = ?", param.Installed)
	}

	if param.Versions != "" {
		searchVersions := "%" + param.Versions + "%"
		tx = tx.Where("versions LIKE ?", searchVersions)
	}

	if param.Tags != "" {
		searchTags := "%" + param.Tags + "%"
		tx = tx.Where("tags LIKE ?", searchTags)
	}
	ls := []*models.Software{}
	find := tx.Find(&ls)
	if find.Error != nil && !errors.Is(find.Error, gorm.ErrRecordNotFound) {
		return nil, tx.Error
	}

	res := []*output.Software{}
	for _, v := range ls {
		toNew, err := convertOldToNew(v)
		if err != nil {
			return nil, err
		}
		res = append(res, toNew)
	}
	return res, nil
}

func convertOldToNew(old *models.Software) (*output.Software, error) {
	ps := []*output.SoftParam{}
	if old.Params != "" {
		err := json.Unmarshal([]byte(old.Params), &ps)
		if err != nil {
			return nil, err
		}
	}
	newSoftware := &output.Software{
		Id:        old.Id,
		Name:      old.Name,
		Key:       old.Key,
		Icon:      old.Icon,
		Type:      old.Type,
		Status:    old.Status,
		Resource:  old.Resource,
		Installed: old.Installed,
		Log:       old.Log,
		Params:    ps,
	}
	if old.Version != "" {
		newSoftware.Version = strings.Split(old.Version, ",")
	}
	return newSoftware, nil
}
