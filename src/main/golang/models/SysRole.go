package models

type SysRole struct {
	Id               int32  `json:"id" tableId:"id"`  //
	RoleName         string `json:"roleName"`         // 角色名称
	Deleted          int32  `json:"deleted"`          // 删除状态（1：已删除；０：未删除）
	CreatedBy        int32  `json:"createdBy"`        //
	CreatedDate      string `json:"createdDate"`      //
	LastModifiedBy   int32  `json:"lastModifiedBy"`   //
	LastModifiedDate string `json:"lastModifiedDate"` //
}

func (SysRole) TableName() string {
	return "sys_role"
}
