package app

import (
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"oneinstack/internal/models"
	"oneinstack/utils"
)

func InitDB(dbPath string) error {
	d, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}
	db = d
	// 检查是否存在用户，如果不存在提示创建管理员
	err = createTables()
	if err != nil {
		log.Fatal("failed to migrate the database:", err)
	}
	err = initUser()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func createTables() error {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.Storage{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.Library{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.Software{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.Website{})
	if err != nil {
		return err
	}
	err = initSoftware()
	if err != nil {
		return err
	}
	return nil
}

func initSoftware() error {
	softToSeed := []*models.Software{
		{
			Name:      "Mysql",
			Key:       "db",
			Icon:      "",
			Type:      "",
			Status:    0,
			Resource:  "local",
			Installed: false,
			Version:   "5.5,5.7,8.0",
			Tags:      "",
			Params: `
	[{
		"key": "port",
		"name": "port",
		"rule": "port",
		"required": "true",
		"type": "number"
	},
	{
		"key": "pwd",
		"name": "pwd",
		"rule": "pwd",
		"required": "true",
		"type": "input"
	},
	{
		"key": "username",
		"name": "username",
		"rule": "username",
		"required": "true",
		"type": "username"
	}]`,
		},
		{
			Name:      "Redis",
			Key:       "redis",
			Icon:      "",
			Type:      "",
			Status:    0,
			Resource:  "local",
			Installed: false,
			Version:   "6.2.0,7.0.5",
			Tags:      "",
			Params:    "",
		},
		{
			Name:      "Nginx",
			Key:       "webserver",
			Icon:      "",
			Type:      "",
			Status:    0,
			Resource:  "local",
			Installed: false,
			Version:   "1.24.0",
			Tags:      "",
			Params:    "",
		},
		{
			Name:      "PHP",
			Key:       "php",
			Icon:      "",
			Type:      "",
			Status:    0,
			Resource:  "local",
			Installed: false,
			Version:   "5.6,7.4,8.1",
			Tags:      "",
			Params:    "",
		},
	}
	var soft models.Software
	result := db.Where("resource = ?", "local").First(&soft)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	if soft.Id > 0 {
		return nil
	}
	tx := db.CreateInBatches(softToSeed, len(softToSeed))
	return tx.Error
}

func initUser() error {
	var count int64 = 0
	tx := DB().Model(models.User{}).Count(&count)
	if tx.Error != nil {
		return tx.Error
	}
	if count > 0 {
		return nil
	}
	err := setupAdminUser()
	if err != nil {
		return err
	}
	return nil
}

func setupAdminUser() error {
	username := utils.GenerateRandomString(8, 12)
	password := utils.GenerateRandomString(8, 12) // 生成 8-12 位随机密码
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	user := &models.User{
		Username: username,
		Password: hashed,
		IsAdmin:  true,
	}
	tx := DB().Create(user)
	if tx.Error != nil {
		return tx.Error
	}
	fmt.Printf("Admin user created successfully.\nUsername: %s\nPassword: %s\n", "admin", password)
	return nil
}
