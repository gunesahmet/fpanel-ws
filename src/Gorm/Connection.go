package Gorm

import (
	"Helpers"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() error {
	dsn := fmt.Sprintf("fpaneladmin:%s@tcp(127.0.0.1:3306)/fpanel?charset=utf8mb4&parseTime=True&loc=Local", Env.FpanelAdminPasswd)
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		Helpers.HandleError(err, 1007, "Database connection failed")
		return err
	}
	Conn = Connection{DB: db}
	return nil
}
