package domain_authSvc

import "gorm.io/gorm"

type userStatus string

const (
	Blocked  userStatus = "blocked"
	Deleted  userStatus = "deleted"
	Pending  userStatus = "pending"
	Active   userStatus = "active"
	verified userStatus = "verified"
	Rejected userStatus = "rejected"
)

type Users struct {
	gorm.Model
	Name          string
	UserName      string
	Email         string
	Password      string
	Bio           string
	ProfileImgUrl string
	Links         string
	Status        userStatus `gorm:"default:pending"`
}
