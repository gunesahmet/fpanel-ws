package Gorm

import (
	"Environments"
	"gorm.io/gorm"
)

var Env *Environments.Env

type Connection struct {
	DB *gorm.DB
}

var Conn Connection
