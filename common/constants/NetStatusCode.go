package constants

var (
	SUCCESS                  int32 = 200
	SYSTEM_ERROR             int32 = 500 // 系统错误
	PARAMETER_ERROR          int32 = 501 // 参数错误
	PARAMETER_MISSING_ERROR  int32 = 502 // 参数缺失
	TOKEN_ERROR              int32 = 503 // token错误
	DATA_NOT_EXIST           int32 = 504 // 数据不存在
	USER_EMPTY_CODE          int32 = 510 // 用户不存在
	USER_PASSWORD_ERROR_CODE int32 = 511 // 密码错误
	ROLE_NOT_CONFIGURED      int32 = 512 // 角色未配置
	USER_FORBIDDEN_CODE      int32 = 511 // 用户被禁用
	USER_OLD_PASSWORD_ERROR  int32 = 514 // 旧密码错误
	UPDATE_ERROR             int32 = 610 // 更新失败
	DELETE_ERROR             int32 = 620 // 删除失败
	INSERT_ERROR             int32 = 630 // 插入失败
)
