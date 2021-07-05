package db

import (
	"alertmanager_notifier/pkg/utils"

	// import postgres driver dependent
	_ "gorm.io/driver/postgres"
	_ "gorm.io/driver/sqlite"
)

// TableUserInfo table userinfo name
const TableUserInfo = "userinfo"

// TableUserRota table userrota name
const TableUserRota = "userrota"

// TableReceiverInfo table receiverinfo name
const TableReceiverInfo = "receiverinfo"

// TableAdmin table admin name
const TableAdmin = "admin"

// Model db tabel common attribute
type Model struct {
	ID int `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
}

// UserInfo table userinfo
type UserInfo struct {
	Model
	UID        string `json:"uid" gorm:"unique"`
	Uname      string `json:"uname"`
	Department string `json:"department"`
	Group      string `json:"group"`
	Enabled    bool   `json:"enabled"`
	Comment    string `json:"comment"`
}

// TableName specify table name for gorm
func (UserInfo) TableName() string {
	return TableUserInfo
}

// UserRota table Rota
type UserRota struct {
	Model
	UID  string         `json:"uid"`
	Date utils.JSONDate `form:"date" json:"date" gorm:"column:date"`
}

// TableName specify table name for gorm
func (UserRota) TableName() string {
	return TableUserRota
}

// ReceiverInfo table userinfo
type ReceiverInfo struct {
	Model
	UID        string `json:"uid" gorm:"uniqueIndex:idx_name"`
	ReceiverType      string `json:"receiverType" gorm:"uniqueIndex:idx_name"`
	ReceiverName string `json:"receiverName"`
}

// TableName specify table name for gorm
func (ReceiverInfo) TableName() string {
	return TableReceiverInfo
}

// Admin table Admin
type Admin struct {
	Model
	Username string `gorm:"unique"`
	Password string
}

// TableName specify table name for gorm
func (Admin) TableName() string {
	return TableAdmin
}
