package vo

type UserVO struct {
	Id               int32    `json:"id"`               // id
	Username         string   `json:"username"`         // 用户名
	Password         string   `json:"password"`         //密码
	OldPassword      string   `json:"oldPassword"`      // 旧密码 - 修改密码用
	Role             int32    `json:"role"`             // 角色ID
	RoleName         string   `json:"roleName"`         // 角色Name
	Status           int32    `json:"status"`           // 状态 0:不可用 1:可用
	Token            string   `json:"token"`            // 登录token
	CreatedBy        int32    `json:"createdBy"`        // 创建者
	CreatedDate      string   `json:"createdDate"`      // 创建时间
	LastModifiedBy   int32    `json:"lastModifiedBy"`   // 更新者
	LastModifiedDate string   `json:"lastModifiedDate"` // 最后更新时间
	Routers          []string `json:"routers"`
}

type UserPageVO struct {
	DefaultPageVO
	Username string `json:"username"` // ユーザ名
}
