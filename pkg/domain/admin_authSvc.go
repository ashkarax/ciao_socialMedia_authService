package domain_authSvc

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Name     string
	Email    string `gorm:"primary key"`
	Password string
}
