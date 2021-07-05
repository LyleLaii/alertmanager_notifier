package e

// MsgFlags  msg
var MsgFlags = map[int]string{
	SUCCESS:            "SUCCESS",
	ACCEPTED:           "SUCCESS",
	ERROR:              "内部错误",
	INVALID_PARAMS:     "请求参数错误",
	ERROR_AUTH:         "用户名或密码错误",
	ERROR_DB_OPERATION: "数据库操作错误",
	ERROR_ADD_USER:     "添加用户信息失败",
	ERROR_EDIT_USER:    "编辑用户信息失败",
	ERROR_DELETE_USER:  "删除用户信息失败",
	ERROR_ADD_ROTA:     "添加值班信息失败",
	ERROR_EDIT_ROTA:    "编辑值班信息失败",
	ERROR_DELETE_ROTA:  "删除值班表信息失败",
	ERROR_ADD_RECEIVERINFO:     "添加接收信息失败",
	ERROR_EDIT_RECEIVERINFO:    "编辑接收信息失败",
	ERROR_DELETE_RECEIVERINFO:  "删除接收信息失败",
}

// GetMsg return msg
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
