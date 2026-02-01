package models

type SysRoleMenu struct {
	Id     int   `json:"id" tableId:"id"` //
	RoleId int32 `json:"roleId"`          // 角色ID
	MenuId int32 `json:"menuId"`          // 菜单ID
}

func (SysRoleMenu) TableName() string {
	return "sys_role_menu"
}
