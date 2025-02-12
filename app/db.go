package app

import (
	"errors"
	"fmt"
	"log"
	"oneinstack/internal/models"
	"oneinstack/utils"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dbPath string) error {
	fmt.Println("创建db..." + dbPath)
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
	err = db.AutoMigrate(&models.Remark{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.Dictionary{})
	if err != nil {
		return err
	}
	err = initSoftware()
	if err != nil {
		return err
	}
	err = initDic()
	if err != nil {
		return err
	}
	err = initRemark()
	if err != nil {
		return err
	}
	return nil
}

func initSoftware() error {
	softToSeed := []*models.Software{
		{
			Name:      "Mysql-5.5",
			Key:       "db",
			Icon:      "data:image/png;base64,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",
			Type:      "",
			Status:    0,
			Resource:  "local",
			Installed: false,
			Version:   "5.5",
			Describe:  "Mysql数据库",
			Tags:      "数据库",
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
			Name:      "Mysql-5.7",
			Key:       "db",
			Icon:      "data:image/png;base64,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",
			Type:      "",
			Status:    0,
			Resource:  "local",
			Installed: false,
			Describe:  "Mysql数据库",

			Version: "5.7",
			Tags:    "数据库",
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
			Name:      "Mysql-8.0",
			Key:       "db",
			Icon:      "data:image/png;base64,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",
			Type:      "",
			Status:    0,
			Resource:  "local",
			Installed: false,
			Version:   "8.0",
			Describe:  "Mysql数据库",

			Tags: "数据库",
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
			Name:      "Redis-6",
			Key:       "redis",
			Icon:      "",
			Type:      "",
			Status:    0,
			Resource:  "local",
			Installed: false,
			Version:   "6.2.0",
			Describe:  "Redis缓存",

			Tags:   "缓存",
			Params: "",
		},
		{
			Name:      "Redis-7",
			Key:       "redis",
			Icon:      "",
			Type:      "",
			Status:    0,
			Resource:  "local",
			Installed: false,
			Version:   "7.0.5",
			Describe:  "Redis缓存",

			Tags:   "缓存",
			Params: "",
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
			Describe:  "Nginx代理服务",
			Tags:      "web服务",
			Params:    "",
		},
		{
			Name:      "PHP-5",
			Key:       "php",
			Icon:      "",
			Type:      "",
			Status:    0,
			Resource:  "local",
			Installed: false,
			Version:   "5.6",
			Tags:      "php",
			Params:    "",
		},
		{
			Name:      "PHP-7",
			Key:       "php",
			Icon:      "",
			Type:      "",
			Status:    0,
			Resource:  "local",
			Installed: false,
			Version:   "7.4",
			Tags:      "php",
			Params:    "",
		},
		{
			Name:      "PHP-8",
			Key:       "php",
			Icon:      "",
			Type:      "",
			Status:    0,
			Resource:  "local",
			Installed: false,
			Version:   "8.1",
			Tags:      "php",
			Params:    "",
		},
		{
			Name:      "phpmyadmin",
			Key:       "phpmyadmin",
			Icon:      "",
			Type:      "",
			Status:    0,
			Resource:  "local",
			Installed: false,
			Version:   "5.2.1",
			Tags:      "php",
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

func initDic() error {
	r := []*models.Dictionary{{
		Key:   "数据库",
		Value: "数据库",
		Q:     "soft_tags",
	},
		{
			Key:   "缓存",
			Value: "缓存",
			Q:     "soft_tags",
		},
		{
			Key:   "web服务",
			Value: "web服务",
			Q:     "soft_tags",
		},
		{
			Key:   "php",
			Value: "php",
			Q:     "soft_tags",
		}}
	var dic models.Dictionary
	result := db.First(&dic)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	if dic.ID > 0 {
		return nil
	}
	tx := db.CreateInBatches(r, len(r))
	return tx.Error
}

func initRemark() error {
	r := &models.Remark{
		Content: "",
	}
	result := db.First(r)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	if r.ID > 0 {
		return nil
	}
	tx := db.Create(r)
	return tx.Error
}

func InitUser() error {
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
	fmt.Printf("用户创建成功.\n用户名: %s\n密码: %s\n", username, password)
	return nil
}
