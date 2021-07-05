package v1

import (
	"alertmanager_notifier/config"
	"alertmanager_notifier/log"
	"alertmanager_notifier/modules/db"
	e "alertmanager_notifier/pkg/err"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const moduleName = "USERINFO"

// GetUserListPage get user list by page
func GetUserListPage(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.Query("page")
		pagea, _ := strconv.Atoi(page)
		filter := c.Query("name")
		ul := db.GetUserList(filter)
		len := len(ul)
		st := (pagea - 1) * 20
		ed := pagea * 20
		if ed > len {
			ed = len
		}
		ul = ul[st:ed]
		data, err := json.Marshal(ul)
		if err != nil {
			logger.Warn("USERINFO", fmt.Sprintf("get rota list page find err: %v", err))
		}

		c.JSON(http.StatusOK, gin.H{
			"code":  e.SUCCESS,
			"msg":   e.GetMsg(e.SUCCESS),
			"total": len,
			"data":  string(data),
		})
	}
}

// DeleteUser delete one user
func DeleteUser(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var l config.ListInfo
		if err := c.BindJSON(&l); err != nil {
			logger.Warn("USERINFO", fmt.Sprintf("deleteuser bind json find err: %v", err))
			c.JSON(http.StatusBadRequest, gin.H{
				"code": e.INVALID_PARAMS,
				"msg":  e.GetMsg(e.INVALID_PARAMS),
			})
			return
		}
		ids := []int{l.ID}
		if err := db.DeleteUsers(ids); err != nil {
			logger.Warn("USERINFO", fmt.Sprintf("deleteuser db operation find err: %v", err))
			c.JSON(http.StatusOK, gin.H{
				"code": e.ERROR_DELETE_USER,
				"msg":  e.GetMsg(e.ERROR_DELETE_USER),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": e.SUCCESS,
			"msg":  e.GetMsg(e.SUCCESS),
		})
	}
}

// DeleteUserBatch delete some users
func DeleteUserBatch(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var l config.ListInfo
		if err := c.BindJSON(&l); err != nil {
			logger.Warn("USERINFO", fmt.Sprintf("deleteuserbatch bind json find err: %v", err))
			c.JSON(http.StatusBadRequest, gin.H{
				"code": e.INVALID_PARAMS,
				"msg":  e.GetMsg(e.INVALID_PARAMS),
			})
			return
		}
		if err := db.DeleteUsers(l.Ids); err != nil {
			logger.Warn("USERINFO", fmt.Sprintf("deleteuserbatch db operation find err: %v", err))
			c.JSON(http.StatusOK, gin.H{
				"code": e.ERROR_DELETE_USER,
				"msg":  e.GetMsg(e.ERROR_DELETE_USER),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": e.SUCCESS,
			"msg":  e.GetMsg(e.SUCCESS),
		})
	}
}

// EditUser edit user info
func EditUser(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var u db.UserInfo
		if err := c.BindJSON(&u); err != nil {
			logger.Warn("USERINFO", fmt.Sprintf("edituser bind json find err: %v", err))
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   e.INVALID_PARAMS,
				"msg":    e.GetMsg(e.INVALID_PARAMS),
				"detail": "",
			})
			return
		}
		if err := db.EditUser(&u); err != nil {
			logger.Warn("USERINFO", fmt.Sprintf("edituser db operation find err: %v", err))
			c.JSON(http.StatusOK, gin.H{
				"code":   e.ERROR_EDIT_USER,
				"msg":    e.GetMsg(e.ERROR_EDIT_USER),
				"detail": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   e.SUCCESS,
			"msg":    e.GetMsg(e.SUCCESS),
			"detail": "",
		})
	}
}

// AddUser add user info
func AddUser(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var u db.UserInfo
		if err := c.BindJSON(&u); err != nil {
			logger.Warn("USERINFO", fmt.Sprintf("adduser bind json find err: %v", err))
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   e.INVALID_PARAMS,
				"msg":    e.GetMsg(e.INVALID_PARAMS),
				"detail": err,
			})
			return
		}
		if err := db.AddUser(&u); err != nil {
			logger.Warn("USERINFO", fmt.Sprintf("adduser db operation find err: %v", err))
			c.JSON(http.StatusOK, gin.H{
				"code":   e.ERROR_ADD_USER,
				"msg":    e.GetMsg(e.ERROR_ADD_USER),
				"detail": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   e.SUCCESS,
			"msg":    e.GetMsg(e.SUCCESS),
			"detail": "",
		})
	}
}

// GetAllUID get all user uid
func GetAllUID(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// var uids []string
		uids, err := db.GetAllUser()
		if err != nil {
			logger.Warn("USERINFO", fmt.Sprintf("getalluid db operation find err: %v", err))
			c.JSON(http.StatusOK, gin.H{
				"code": e.ERROR_DB_OPERATION,
				"msg":  e.GetMsg(e.ERROR_DB_OPERATION),
				"data": make([]string, 0),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": e.SUCCESS,
			"msg":  e.GetMsg(e.SUCCESS),
			"data": uids,
		})
	}
}

// GetUserInfo Get user info by uid
func GetUserInfo(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.Query("uid")
		u, err := db.GetUserInfo(uid)
		if err != nil {
			logger.Warn("USERINFO", fmt.Sprintf("getuserinfo db operation find err: %v", err))
			c.JSON(http.StatusOK, gin.H{
				"code": e.ERROR_DB_OPERATION,
				"msg":  e.GetMsg(e.ERROR_DB_OPERATION),
				"data": make([]string, 0),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": e.SUCCESS,
			"msg":  e.GetMsg(e.SUCCESS),
			"data": u,
		})
	}
}
