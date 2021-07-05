package db

import (
	"alertmanager_notifier/pkg/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// var db = DBInstance

// UserRotaDetail userrotedetail  table description
type UserRotaDetail struct {
	ID         int            `json:"id"`
	UID        string         `json:"uid"`
	Uname      string         `json:"uname"`
	Department string         `json:"department"`
	Group      string         `json:"group"`
	Date       utils.JSONDate `form:"date" json:"date" gorm:"column:date"`
}

// UserReceiverInfoDetail userreceiverinfodetail  table description
type UserReceiverInfoDetail struct {
	ID         int            `json:"id"`
	UID        string         `json:"uid"`
	Uname      string         `json:"uname"`
	ReceiverType string         `json:"receiverType"`
	ReceiverName      string         `json:"receiverName"`
}

// CheckAdmin check Admin user exists
// TODO: Fix this: 如果没有找到记录，err不为空
func CheckAdmin(username string, password string) bool {
	var u Admin
	encodePassword := utils.Encode(password)
	db := DBInstance.Where("username = ? AND password = ?", username, encodePassword).Find(&u)

	err := db.Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	return true

}

// GetUserList get userlist
func GetUserList(filterName string) (users []UserInfo) {
	if filterName == "" {
		DBInstance.Find(&users)
	} else {
		sql := fmt.Sprintf("SELECT * FROM %s WHERE (uid LIKE '%%%s%%' ) OR (uname LIKE '%%%s%%' )",
			TableUserInfo,
			filterName,
			filterName)
		DBInstance.Raw(sql).Scan(&users)
	}
	return
}

// DeleteUsers delete users
func DeleteUsers(ids []int) (err error) {
	err = DBInstance.Where("id IN (?)", ids).Delete(UserInfo{}).Error
	return
}

// EditUser edit user
func EditUser(u *UserInfo) (err error) {
	err = DBInstance.Save(&u).Error
	return
}

// AddUser add user
func AddUser(u *UserInfo) (err error) {
	err = DBInstance.Create(&u).Error
	return
}

// GetRotaList get userlist
func GetRotaList(filterName string) (rotas []UserRotaDetail) {
	if filterName == "" {
		// TODO: fix wthat?
		sql := "SELECT userrota.id, userrota.uid, uname, department, \"group\", date FROM userrota JOIN userinfo ON (userrota.uid = userinfo.uid)"
		DBInstance.Raw(sql).Scan(&rotas)
	} else {
		// TODO: fix what?
		sql := fmt.Sprintf("SELECT userrota.id,userrota.uid,uname,department,\"group\",date FROM userinfo RIGHT JOIN userrota ON (userinfo.uid = userrota.uid) WHERE (userrota.uid LIKE '%%%s%%' ) OR (uname LIKE '%%%s%%' );",
			filterName,
			filterName)
		DBInstance.Raw(sql).Scan(&rotas)
	}
	return
}

// DeleteRotas delete rotas
func DeleteRotas(ids []int) (err error) {
	err = DBInstance.Where("id IN (?)", ids).Delete(UserRota{}).Error
	return
}

// EditRota edit rota
func EditRota(r *UserRota) (err error) {
	err = DBInstance.Save(&r).Error
	return
}

// AddRota add rota
func AddRota(r *UserRota) (err error) {
	err = DBInstance.Create(&r).Error
	return
}

// GetReceiverInfoList get receiverinfo
func GetReceiverInfoList(filterName string) (receiverinfodetails []UserReceiverInfoDetail) {
	if filterName == "" {
		// TODO: fix wthat?
		sql := "SELECT receiverinfo.id, receiverinfo.uid, uname, receiver_type, receiver_name FROM receiverinfo JOIN userinfo ON (receiverinfo.uid = userinfo.uid)"
		DBInstance.Raw(sql).Scan(&receiverinfodetails)
	} else {
		// TODO: fix what?
		sql := fmt.Sprintf("SELECT receiverinfo.id,receiverinfo.uid, uname, receiver_type, receiver_name FROM userinfo RIGHT JOIN userrota ON (userinfo.uid = receiverinfo.uid) WHERE (receiverinfo.uid LIKE '%%%s%%' ) OR (uname LIKE '%%%s%%' );",
			filterName,
			filterName)
		DBInstance.Raw(sql).Scan(&receiverinfodetails)
	}
	return
}

// DeleteReceiverInfos delete receiverinfo
func DeleteReceiverInfos(ids []int) (err error) {
	err = DBInstance.Where("id IN (?)", ids).Delete(ReceiverInfo{}).Error
	return
}

// EditReceiverInfo edit receiverinfo
func EditReceiverInfo(r *ReceiverInfo) (err error) {
	err = DBInstance.Save(&r).Error
	return
}

// AddReceiverInfo add receiverinf
func AddReceiverInfo(r *ReceiverInfo) (err error) {
	err = DBInstance.Create(&r).Error
	return
}

// GetUserInfo get user info by user uid
func GetUserInfo(uid string) (userInfo UserInfo, err error) {
	err = DBInstance.Where("uid = ?", uid).Find(&userInfo).Error
	return
}

// GetStandbyUser get standby user by group\date
func GetStandbyUser(department string, group string, date string) (rotausers []string) {
	if utils.GetUseRote() {
		var rotas []UserRotaDetail
		sql := fmt.Sprintf("WITH r AS (SELECT * FROM userrota WHERE date='%s') SELECT r.id,r.uid,uname,department,\"group\",date FROM r JOIN userinfo ON (r.uid = userinfo.uid) WHERE enabled = true AND (department = '%s' ) AND (\"group\" = '%s' );",
			date,
			department,
			group)
		DBInstance.Raw(sql).Scan(&rotas)
		for _, u := range rotas {
			rotausers = append(rotausers, u.UID)
		}

	} else {
		// log.Info("DBHander", "don't use rota by setting")
		var users []UserInfo
		DBInstance.Where("enabled = true AND department = ? AND \"group\" = ?", department, group).Find(&users)
		for _, u := range users {
			rotausers = append(rotausers, u.UID)
		}
	}
	return
}

// GetComment get user work wechat name by username
func GetReceiverNameByType(users []string, receiverType string) (receiverNames []string) {
	var receiverInfos []ReceiverInfo
	DBInstance.Where("uid IN (?) AND receiver_type = ?", users, receiverType).Find(&receiverInfos)
	for _, r := range receiverInfos {
		receiverNames = append(receiverNames, r.ReceiverName)
	}
	return
}

// GetComment get user work wechat name by username
func GetComment(users []string) (comments []string) {
	var userinfos []UserInfo
	DBInstance.Where("uid IN (?)", users).Find(&userinfos)
	for _, u := range userinfos {
		comments = append(comments, u.Comment)
	}
	return
}

// GetAllUser get all user info
func GetAllUser() (uids []string, err error) {
	// var users []UserInfo
	// DBInstance.Find(&users).Pluck("uid", &uids)
	err = DBInstance.Model(&UserInfo{}).Pluck("uid", &uids).Error
	return
}
