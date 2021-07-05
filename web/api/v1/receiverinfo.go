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

// GetReceiverInfoListPage get receiverinfo list by page
func GetReceiverInfoListPage(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.Query("page")
		pagea, _ := strconv.Atoi(page)
		filter := c.Query("name")
		ril := db.GetReceiverInfoList(filter)
		len := len(ril)
		st := (pagea - 1) * 20
		ed := pagea * 20
		if ed > len {
			ed = len
		}
		ril = ril[st:ed]
		data, err := json.Marshal(ril)
		if err != nil {
			logger.Warn("RECEIVERINFO", fmt.Sprintf("get rota list page find err: %v", err))
		}

		c.JSON(http.StatusOK, gin.H{
			"code":  e.SUCCESS,
			"msg":   e.GetMsg(e.SUCCESS),
			"total": len,
			"data":  string(data),
		})
	}
}

// DeleteReceiverInfo delete one receiveringo
func DeleteReceiverInfo(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var l config.ListInfo
		if err := c.BindJSON(&l); err != nil {
			logger.Warn("ReceiverInfo", fmt.Sprintf("deletereceiverinfo bind json find err: %v", err))
			c.JSON(http.StatusBadRequest, gin.H{
				"code": e.INVALID_PARAMS,
				"msg":  e.GetMsg(e.INVALID_PARAMS),
			})
			return
		}
		ids := []int{l.ID}
		if err := db.DeleteReceiverInfos(ids); err != nil {
			logger.Warn("ReceiverInfo", fmt.Sprintf("deletereceiver db operation find err: %v", err))
			c.JSON(http.StatusOK, gin.H{
				"code": e.ERROR_DELETE_RECEIVERINFO,
				"msg":  e.GetMsg(e.ERROR_DELETE_RECEIVERINFO),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": e.SUCCESS,
			"msg":  e.GetMsg(e.SUCCESS),
		})
	}
}

// DeleteReceiverInfoBatch delete some ReceiverInfo
func DeleteReceiverInfoBatch(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var l config.ListInfo
		if err := c.BindJSON(&l); err != nil {
			logger.Warn("ReceiverInfo", fmt.Sprintf("deletereceiverinfobatch bind json find err: %v", err))
			c.JSON(http.StatusBadRequest, gin.H{
				"code": e.INVALID_PARAMS,
				"msg":  e.GetMsg(e.INVALID_PARAMS),
			})
			return
		}
		if err := db.DeleteReceiverInfos(l.Ids); err != nil {
			logger.Warn("ReceiverInfo", fmt.Sprintf("deletereceiverinfobatch db operation find err: %v", err))
			c.JSON(http.StatusOK, gin.H{
				"code": e.ERROR_DELETE_RECEIVERINFO,
				"msg":  e.GetMsg(e.ERROR_DELETE_RECEIVERINFO),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": e.SUCCESS,
			"msg":  e.GetMsg(e.SUCCESS),
		})
	}
}

// EditReceiverInfo edit receiverinfo
func EditReceiverInfo(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r db.ReceiverInfo
		if err := c.BindJSON(&r); err != nil {
			logger.Warn("ReceiverInfo", fmt.Sprintf("editreceiver bind json find err: %v", err))
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   e.INVALID_PARAMS,
				"msg":    e.GetMsg(e.INVALID_PARAMS),
				"detail": "",
			})
			return
		}
		if err := db.EditReceiverInfo(&r); err != nil {
			logger.Warn("ReceiverInfo", fmt.Sprintf("editreceiver db operation find err: %v", err))
			c.JSON(http.StatusOK, gin.H{
				"code":   e.ERROR_EDIT_RECEIVERINFO,
				"msg":    e.GetMsg(e.ERROR_EDIT_RECEIVERINFO),
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

// AddReceiverInfo add receiver info
func AddReceiverInfo(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r db.ReceiverInfo
		if err := c.BindJSON(&r); err != nil {
			logger.Warn("ReceiverInfo", fmt.Sprintf("addreceiverinfo bind json find err: %v", err))
			c.JSON(http.StatusBadRequest, gin.H{
				"code":   e.INVALID_PARAMS,
				"msg":    e.GetMsg(e.INVALID_PARAMS),
				"detail": err,
			})
			return
		}
		if err := db.AddReceiverInfo(&r); err != nil {
			logger.Warn("ReceiverInfo", fmt.Sprintf("addreceviverinfo db operation find err: %v", err))
			c.JSON(http.StatusOK, gin.H{
				"code":   e.ERROR_ADD_RECEIVERINFO,
				"msg":    e.GetMsg(e.ERROR_ADD_RECEIVERINFO),
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