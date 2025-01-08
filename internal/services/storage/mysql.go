package storage

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math"
	"oneinstack/app"
	"oneinstack/internal/models"
)

type MysqlOP struct {
	ID       int64
	Addr     string
	Port     string
	Root     string
	Password string
	Type     string
	Lib      string
	DB       *gorm.DB
}
type DbInfo struct {
	DbName string
	Usage  float64
}
type UserPrivilege struct {
	Db   string
	User string
	Host string
}

func NewMysqlOP(p *models.Storage, lib string) *MysqlOP {
	return &MysqlOP{
		ID:       p.ID,
		Addr:     p.Addr,
		Port:     p.Port,
		Root:     p.Root,
		Password: p.Password,
		Type:     p.Type,
		DB:       nil,
		Lib:      lib,
	}
}

func (s *MysqlOP) Connet() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local", s.Root, s.Password, s.Addr, s.Port, s.Lib)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	s.DB = db
	return nil
}

func (s *MysqlOP) Sync() error {
	// 获取所有数据库及其大小
	var dbInfos []DbInfo
	err := s.DB.Raw(`
		SELECT 
			TABLE_SCHEMA as DbName, 
			SUM(DATA_LENGTH + INDEX_LENGTH) as ` + "`Usage`" + `
		FROM 
			TABLES
		GROUP BY 
			TABLE_SCHEMA
	`).Scan(&dbInfos).Error
	if err != nil {
		return err
	}
	ls := []models.Library{}
	// 获取每个数据库的用户权限信息
	for _, dbInfo := range dbInfos {
		if dbInfo.DbName == "information_schema" || dbInfo.DbName == "mysql" || dbInfo.DbName == "performance_schema" || dbInfo.DbName == "sys" {
			continue
		}
		var userPrivileges []UserPrivilege
		err = s.DB.Raw(`
			SELECT DISTINCT 
				DB, 
				User, 
				Host
			FROM 
				mysql.db
			WHERE 
				DB = ?
		`, dbInfo.DbName).Scan(&userPrivileges).Error
		if err != nil {
			return err
		}

		l := models.Library{
			PID:      s.ID,
			Name:     dbInfo.DbName,
			User:     "",
			Password: "",
			Capacity: ConvertBytes(dbInfo.Usage),
			PAddr:    fmt.Sprintf("%s:%v", s.Addr, s.Port),
			Type:     s.Type,
		}
		// 输出数据库信息和其访问用户
		fmt.Printf("Database: %s, Usage: %.2f bytes\n", dbInfo.DbName, dbInfo.Usage)
		if len(userPrivileges) > 0 {
			l.User = userPrivileges[0].User
		}
		ls = append(ls, l)
	}
	tx := app.DB().Where("p_id = ? ", s.ID).Delete(&models.Library{})
	if tx.Error != nil {
		return tx.Error
	}
	tx = app.DB().Create(ls)
	return tx.Error
}

func ConvertBytes(bytes float64) string {
	if bytes == 0 {
		return "0 B"
	}
	units := []string{"B", "KB", "MB", "GB", "TB"}
	log := math.Log(bytes) / math.Log(1024)
	unitIndex := int(log)
	if unitIndex >= len(units) {
		unitIndex = len(units) - 1
	}
	return fmt.Sprintf("%.2f %s", bytes/math.Pow(1024, float64(unitIndex)), units[unitIndex])
}
