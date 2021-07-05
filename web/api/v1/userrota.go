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

// GetRotaListPage get user list by page
func GetRotaListPage(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.Query("page")
		pagea, _ := strconv.Atoi(page)
		filter := c.Query("name")
		rl := db.GetRotaList(filter)
		len := len(rl)
		st := (pagea - 1) * 20
		ed := pagea * 20
		if ed > len {
			ed = len
		}
		rl = rl[st:ed]
		data, err := json.Marshal(rl)
		if err != nil {
			logger.Warn("USERROTA", fmt.Sprintf("get rota list page find err: %v", err))
		}

		c.JSON(http.StatusOK, gin.H{
			"code":  e.SUCCESS,
			"msg":   e.GetMsg(e.SUCCESS),
			"total": len,
			"data":  string(data),
		})
	}
}

// DeleteRota delete one rota
func DeleteRota(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var l config.ListInfo
		if err := c.BindJSON(&l); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": e.INVALID_PARAMS,
				"msg":  e.GetMsg(e.INVALID_PARAMS),
			})
			logger.Warn("USERROTA", fmt.Sprintf("deleterota bind json data find err: %v", err))
			return
		}
		ids := []int{l.ID}
		if err := db.DeleteRotas(ids); err != nil {
			logger.Warn("USERROTA", fmt.Sprintf("deleterota db operation find err: %v", err))
			c.JSON(http.StatusOK, gin.H{
				"code": e.ERROR_DELETE_ROTA,
				"msg":  e.GetMsg(e.ERROR_DELETE_ROTA),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": e.SUCCESS,
			"msg":  e.GetMsg(e.SUCCESS),
		})
	}
}

// DeleteRotaBatch delete some rotas
func DeleteRotaBatch(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var l config.ListInfo
		if err := c.BindJSON(&l); err != nil {
			logger.Warn("USERROTA", fmt.Sprintf("deleterotabatch bind json data find err: %v", err))
			c.JSON(http.StatusBadRequest, gin.H{
				"code": e.INVALID_PARAMS,
				"msg":  e.GetMsg(e.INVALID_PARAMS),
			})
			return
		}
		if err := db.DeleteRotas(l.Ids); err != nil {
			logger.Warn("USERROTA", fmt.Sprintf("deleterotabatch db operation find err: %v", err))
			c.JSON(http.StatusOK, gin.H{
				"code": e.ERROR_DELETE_ROTA,
				"msg":  e.GetMsg(e.ERROR_DELETE_ROTA),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": e.SUCCESS,
			"msg":  e.GetMsg(e.SUCCESS),
		})
	}
}

// EditRota edit rota info
func EditRota(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r db.UserRota
		if err := c.BindJSON(&r); err != nil {
			logger.Warn("USERROTA", fmt.Sprintf("editrota bind json data find err: %v", err))
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   e.INVALID_PARAMS,
				"msg":    e.GetMsg(e.INVALID_PARAMS),
				"detail": "",
			})
			return
		}
		if err := db.EditRota(&r); err != nil {
			logger.Warn("USERROTA", fmt.Sprintf("editrota db operation find err: %v", err))
			c.JSON(http.StatusOK, gin.H{
				"code":   e.ERROR_EDIT_ROTA,
				"msg":    e.GetMsg(e.ERROR_EDIT_ROTA),
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

// AddRota add rota info
func AddRota(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r db.UserRota
		if err := c.BindJSON(&r); err != nil {
			logger.Warn("USERROTA", fmt.Sprintf("addrota bind json data find err: %v", err))
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   e.INVALID_PARAMS,
				"msg":    e.GetMsg(e.INVALID_PARAMS),
				"detail": err,
			})
			return
		}
		if err := db.AddRota(&r); err != nil {
			logger.Warn("USERROTA", fmt.Sprintf("addrota db operation find err: %v", err))
			c.JSON(http.StatusOK, gin.H{
				"code":   e.ERROR_ADD_ROTA,
				"msg":    e.GetMsg(e.ERROR_ADD_ROTA),
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
