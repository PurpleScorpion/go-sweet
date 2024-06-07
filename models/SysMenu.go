package models

type SysMenu struct {
	Id         int32  `json:"id" tableId:"id"` //
	MenuName   string `json:"menuName"`        // 菜单名称
	RouterName string `json:"routerName"`      // 路由权限名称
	MenuType   int32  `json:"menuType"`        // 菜单类型（1：目录 2：菜单）
	ParentId   int32  `json:"parentId"`        // 父菜单ID
	IsSys      int32  `json:"isSys"`           // 是否是系统菜单(1:系统菜单,不可删除 0:非系统菜单,可删除)
	OrderNum   int32  `json:"orderNum"`        // 展示顺序
}

func (SysMenu) TableName() string {
	return "sys_menu"
}
