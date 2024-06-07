package models

type User struct {
	Id               int32  `json:"id" tableId:"id"`  //
	Username         string `json:"username"`         // 用户名
	Password         string `json:"password"`         // 密码
	Status           int32  `json:"status"`           // 状态（1：启用；0：禁用）
	Role             int32  `json:"role"`             // 角色ID
	Deleted          int32  `json:"deleted"`          // 状态 0:不可用 1:可用
	CreatedBy        int32  `json:"createdBy"`        //
	CreatedDate      string `json:"createdDate"`      //
	LastModifiedBy   int32  `json:"lastModifiedBy"`   //
	LastModifiedDate string `json:"lastModifiedDate"` //
}

func (User) TableName() string {
	return "user"
}
