package constants

var (
	SUCCESS                  int32 = 200
	SYSTEM_ERROR             int32 = 500 // 系统错误
	PARAMETER_ERROR          int32 = 501 // 参数错误
	PARAMETER_MISSING_ERROR  int32 = 502 // 参数缺失
	TOKEN_ERROR              int32 = 503 // token错误
	DATA_NOT_EXIST           int32 = 504 // 数据不存在
	BASE_CODE_ERROR          int32 = 505 // 据点不存在
	INSUFFICIENT_PRIVILEGES  int32 = 506 // 权限不足
	DATA_ALREADY_EXIST       int32 = 507 // 数据已存在
	USER_EMPTY_CODE          int32 = 510 // 用户不存在
	USER_PASSWORD_ERROR_CODE int32 = 511 // 密码错误
	ROLE_NOT_CONFIGURED      int32 = 512 // 角色未配置
	BASE_ROLE_NOT_CONFIGURED int32 = 513 // 据点角色未配置
	USER_OLD_PASSWORD_ERROR  int32 = 514 // 旧密码错误
	USER_FORBIDDEN_CODE      int32 = 515 // 用户被禁用
	FUNC_LOCK                int32 = 516 // 接口已被锁定
	BASE_CODE_MISSING_ERROR  int32 = 517 // 据点未选择
	UPDATE_ERROR             int32 = 610 // 更新失败
	DELETE_ERROR             int32 = 620 // 删除失败
	INSERT_ERROR             int32 = 630 // 插入失败

	// IotHub/EventHub Error
	IOT_HUB_ERROR               int32 = 700 // Iothub 错误
	IOT_HUB_RESTORE_ERROR       int32 = 701 // Iothub 恢复失败
	IOT_HUB_CHANGE_STATUS_ERROR int32 = 702 // Iothub 状态更改失败
	EVENTHUB_PARAMETER_MISSING  int32 = 703 // EventHub 参数缺失
	EVENTHUB_ALREADY_EXIST      int32 = 704 // EventHub 已存在
	IOT_HUB_CHECK_TIME_ERROR    int32 = 705 // Iothub 时间校验错误

	// System Error
	USER_NOT_EXIST                int32 = 800 // 用户不存在
	USER_NAME_REPEAT_ERROR        int32 = 801 // 用户名重复
	BASE_NAME_REPEAT_ERROR        int32 = 802 // 据点名称重复
	BASE_NOT_EXIST                int32 = 803 // 据点不存在
	MENU_NOT_EXIST                int32 = 804 // 菜单不存在
	MENU_CHILD_EXIST_DELETE_ERROR int32 = 805 // 子菜单存在，无法删除
	BASE_CHILD_EXIST_DELETE_ERROR int32 = 806 // 子据点存在，无法删除
	ROLE_NOT_EXIST                int32 = 807 // 角色不存在
	FILE_NAME_REPEAT_ERROR        int32 = 808 // 文件名称重复
	FILE_DELETE_ERROR             int32 = 809 // 文件删除失败
	READ_FILE_ERROR               int32 = 810 // 文件读取失败
	ROLE_NAME_MISSING_ERROR       int32 = 811 // 角色名称参数缺失
	ROLE_NAME_REPEAT_ERROR        int32 = 812 // 角色名称已存在
	USER_NAME_MISSING_ERROR       int32 = 813 // 用户名称参数缺失
	ROLE_NOT_SELECT               int32 = 814 // 角色未选择
	PASSWORD_MISSING_ERROR        int32 = 815 // 密码参数缺失
	BASE_NOT_SELECT               int32 = 816 // 据点未选择

	// File Error
	FILE_NOT_EXIST          int32 = 1000 // 文件不存在
	FILE_OPEN_ERROR         int32 = 1001 // 文件打开失败
	FILE_CREATE_GZ_ERROR    int32 = 1002 // 压缩文件创建失败
	FILE_ANALYSIS_IMG_ERROR int32 = 1003 // 图片解析失败
	FILE_IMG_FORMAT_ERROR   int32 = 1004 // 图片格式错误
	FILE_WRITE_ERROR        int32 = 1005 //文件写入失败

	// Cache Server Error
	CACHE_SERVER_ERROR int32 = 1100 // 缓存服务器异常

	// CAD Error
	CAD_CONF_MISSING_ERROR int32 = 1200 // CAD 错误

	// Calendar Error
	HOLIDAY_NOT_EXIST          int32 = 1300 // 假日不存在
	TEMPLATE_NOT_EXIST         int32 = 1301 // 模板不存在
	TEMPLATE_NAME_REPEAT_ERROR int32 = 1302 // 模板名称重复
	TEMPLATE_USED_DELETE_ERROR int32 = 1303 // 删除失败-模板被使用

	// Device Error
	DEVICE_NOT_EXIST int32 = 1400 // 设备不存在

	// Device Setting Error
	DEVICE_SETTING_NOT_EXIST int32 = 1500 // 配置不存在
)
