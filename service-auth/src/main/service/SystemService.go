package service

import (
	"service-common/src/main/models"
	"shared/constants"
	"shared/utils"
	"shared/vo"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"github.com/PurpleScorpion/go-sweet-orm/v3/mapper"
)

type SystemService struct {
}

func (that *SystemService) UserPageData(userVO vo.UserPageVO) utils.R {
	qw := mapper.BuilderQueryWrapper()
	qw.Like(keqing.IsNotEmpty(userVO.Username), "username", userVO.Username)
	qw.Ne(true, "role", constants.ROOT_ROLE_ID)
	qw.Eq(true, "deleted", constants.NO_DELETE_CODE)
	page := mapper.BuilderPageUtils(userVO.Current, userVO.PageSize, qw)
	pageData := mapper.Page[models.User](page)

	arr := pageData.List
	var vos []vo.UserVO
	for i := 0; i < len(arr); i++ {
		var tmp = vo.UserVO{}
		tmp.Id = arr[i].Id
		tmp.Username = arr[i].Username
		tmp.Role = arr[i].Role
		tmp.LastModifiedDate = arr[i].LastModifiedDate
		tmp.Status = arr[i].Status
		tmp.RoleName = ""
		role := mapper.SelectById[models.SysRole](arr[i].Role)
		if keqing.IsNotEmpty(role) {
			if role[0].Deleted == constants.NO_DELETE_CODE {
				tmp.RoleName = role[0].RoleName
			}
		}
		vos = append(vos, tmp)
	}
	data := mapper.PageData[vo.UserVO]{}
	data.PageSize = pageData.PageSize
	data.Current = pageData.Current
	data.TotalCount = pageData.TotalCount
	data.TotalPage = pageData.TotalPage
	data.List = vos
	return utils.Success(data)
}

func (that *SystemService) RolePageData(roleVO vo.RolePageVO) utils.R {
	qw := mapper.BuilderQueryWrapper()
	qw.Like(keqing.IsNotEmpty(roleVO.RoleName), "role_name", roleVO.RoleName)
	qw.Eq(true, "deleted", constants.NO_DELETE_CODE)
	qw.OrderByTimeAsc(true, "created_by")
	page := mapper.BuilderPageUtils(roleVO.Current, roleVO.PageSize, qw)
	pageData := mapper.Page[models.SysRole](page)
	return utils.Success(pageData)
}

func (that *SystemService) MenuPageData() utils.R {
	pvos := getParentMenu()
	if len(pvos) == 0 {
		return utils.Success(pvos)
	}
	getChildMenu(pvos)
	return utils.Success(pvos)
}

func (that *SystemService) DeleteUserById(id int) utils.R {

	count := mapper.SelectCount[models.User](mapper.BuilderQueryWrapper().
		Eq(true, "id", id).
		Eq(true, "deleted", constants.NO_DELETE_CODE),
	)

	if count == 0 {
		return utils.Fail(constants.DATA_NOT_EXIST, "The user does not exist")
	}

	num := mapper.Update[models.User](mapper.BuilderUpdateWrapper(false).
		Eq(true, "id", id).
		Set(true, "last_modified_date", keqing.NowDateStr()).
		Set(true, "deleted", constants.DELETE_CODE),
	)
	if num == 0 {
		return utils.Fail(constants.UPDATE_ERROR, "Update failed")
	}
	return utils.Success("")
}

func (that *SystemService) ChangeUserStatus(userVO vo.UserVO) utils.R {
	count := mapper.SelectCount[models.User](mapper.BuilderQueryWrapper().
		Eq(true, "id", userVO.Id).
		Eq(true, "deleted", constants.NO_DELETE_CODE),
	)
	if count == 0 {
		return utils.Fail(constants.DATA_NOT_EXIST, "The user does not exist")
	}

	num := mapper.Update[models.User](mapper.BuilderUpdateWrapper(false).
		Eq(true, "id", userVO.Id).
		Set(true, "last_modified_date", keqing.NowDateStr()).
		Set(true, "last_modified_by", userVO.LastModifiedBy).
		Set(true, "status", userVO.Status),
	)

	if num == 0 {
		return utils.Fail(constants.UPDATE_ERROR, "Update failed")
	}
	if userVO.Status == constants.FAIL_STATUS {
		utils.DeleteCache(constants.GetUserExpireTimeKey(userVO.Id))
	}

	return utils.Success(constants.SUCCESS)
}

func (that *SystemService) UserUpdate(userVO vo.UserVO) utils.R {
	msg := validateUser(userVO)
	if msg != "ok" {
		return utils.Fail(constants.PARAMETER_ERROR, msg)
	}

	count := mapper.SelectCount[models.User](mapper.BuilderQueryWrapper().
		Eq(true, "id", userVO.Id).
		Eq(true, "deleted", constants.NO_DELETE_CODE),
	)
	if count == 0 {
		return utils.Fail(constants.DATA_NOT_EXIST, "The user does not exist")
	}

	num := mapper.Update[models.User](mapper.BuilderUpdateWrapper(false).
		Eq(true, "id", userVO.Id).
		Set(keqing.IsNotEmpty(userVO.Password), "password", keqing.MD5Salt(userVO.Password, USER_SALT)).
		Set(true, "role", userVO.Role).
		Set(true, "last_modified_date", keqing.NowDateStr()).
		Set(true, "last_modified_by", userVO.LastModifiedBy),
	)
	if num == 0 {
		return utils.Fail(constants.UPDATE_ERROR, "Update failed")
	}

	return utils.Success("")
}

func (that *SystemService) UserInsert(userVO vo.UserVO) utils.R {
	msg := validateUser(userVO)
	if msg != "ok" {
		return utils.Fail(constants.PARAMETER_ERROR, msg)
	}

	// 判断用户名是否重复
	users := mapper.SelectList[models.User](mapper.BuilderQueryWrapper().
		Eq(true, "username", userVO.Username),
	)
	if keqing.IsNotEmpty(users) {
		// 判断账号状态
		if users[0].Deleted == constants.NO_DELETE_CODE {
			return utils.Fail(constants.INSERT_ERROR, "The username already exists")
		}
		userVO.Id = users[0].Id
		// 恢复账号
		return recoveryUser(userVO)
	}
	var user models.User
	user.Username = userVO.Username
	user.Password = keqing.MD5Salt(userVO.Password, USER_SALT)
	user.Role = userVO.Role
	user.CreatedDate = keqing.NowDateStr()
	user.CreatedBy = userVO.CreatedBy
	user.LastModifiedBy = userVO.LastModifiedBy
	user.LastModifiedDate = keqing.NowDateStr()
	user.Deleted = constants.NO_DELETE_CODE
	user.Status = constants.NORMAL_STATUS
	count := mapper.Insert[models.User](&user, nil)
	if count == 0 {
		return utils.Fail(constants.INSERT_ERROR, "Insert failed")
	}

	return utils.Success("")
}

/*
恢复账号删除状态
*/
func recoveryUser(userVO vo.UserVO) utils.R {
	count := mapper.Update[models.User](mapper.BuilderUpdateWrapper(false).
		Eq(true, "id", userVO.Id).
		Set(keqing.IsNotEmpty(userVO.Password), "password", keqing.MD5Salt(userVO.Password, USER_SALT)).
		Set(true, "role", userVO.Role).
		Set(true, "last_modified_date", keqing.NowDateStr()).
		Set(true, "last_modified_by", userVO.LastModifiedBy).
		Set(true, "deleted", constants.NO_DELETE_CODE),
	)
	if count == 0 {
		return utils.Fail(constants.UPDATE_ERROR, "Update failed")
	}
	return utils.Success("")
}

func (that *SystemService) GetUserById(id int) utils.R {
	users := mapper.SelectById[models.User](id)

	if keqing.IsEmpty(users) {
		return utils.Fail(constants.DATA_NOT_EXIST, "The user does not exist")
	}

	user := users[0]
	var userVO vo.UserVO

	userVO.Id = user.Id
	userVO.Username = user.Username
	userVO.Role = user.Role

	return utils.Success(userVO)
}

func (that *SystemService) MenuInsert(menu models.SysMenu) utils.R {
	menu.IsSys = 0
	count := mapper.Insert[models.SysMenu](&menu, nil)
	if count == 0 {
		return utils.Fail(constants.INSERT_ERROR, "Insert failed")
	}
	return utils.Success("")
}

func (that *SystemService) MenuUpdate(menu models.SysMenu) utils.R {
	count := mapper.Update[models.SysMenu](mapper.BuilderUpdateWrapper(false).
		Eq(true, "id", menu.Id).
		Set(keqing.IsNotEmpty(menu.MenuName), "menu_name", menu.MenuName).
		Set(keqing.IsNotEmpty(menu.RouterName), "router_name", menu.RouterName).
		Set(menu.MenuType > 0, "menu_type", menu.MenuType).
		Set(true, "order_num", menu.OrderNum).
		Set(true, "parent_id", menu.ParentId),
	)
	if count == 0 {
		return utils.Fail(constants.UPDATE_ERROR, "Update failed")
	}
	return utils.Success("")
}

func (that *SystemService) AllParentMenu() utils.R {
	pvos := getParentMenu()
	return utils.Success(pvos)
}

func (that *SystemService) AllRole() utils.R {
	roles := mapper.SelectList[models.SysRole](mapper.BuilderQueryWrapper().
		Eq(true, "deleted", constants.NO_DELETE_CODE),
	)
	return utils.Success(roles)
}

func (that *SystemService) GetMenuById(id int) utils.R {
	menus := mapper.SelectById[models.SysMenu](id)
	if keqing.IsEmpty(menus) {
		return utils.Fail(constants.DATA_NOT_EXIST, "Menu does not exist")
	}
	return utils.Success(menus[0])
}

func (that *SystemService) DeleteMenuById(id int) utils.R {
	menus := mapper.SelectById[models.SysMenu](id)
	if keqing.IsEmpty(menus) {
		return utils.Fail(constants.DATA_NOT_EXIST, "Menu does not exist")
	}
	// 判断是否是父节点 , 且父节点下有未删除的子节点

	count := mapper.SelectCount[models.SysMenu](mapper.BuilderQueryWrapper().
		Eq(true, "parent_id", id),
	)
	if count > 0 {
		return utils.Fail(constants.DELETE_ERROR, "There are submenus in the current directory that have not been deleted")
	}
	num := mapper.DeleteById[models.SysMenu](id, nil)
	if num == 0 {
		return utils.Fail(constants.DELETE_ERROR, "Delete failed, please try again later")
	}
	return utils.Success("")
}

func (that *SystemService) DeleteRoleById(id int) utils.R {

	count := mapper.SelectCount[models.SysRole](mapper.BuilderQueryWrapper().
		Eq(true, "id", id),
	)

	if count == 0 {
		return utils.Fail(constants.DATA_NOT_EXIST, "Menu does not exist")
	}
	num := mapper.DeleteById[models.SysRole](id, nil)
	if num == 0 {
		return utils.Fail(constants.DELETE_ERROR, "Delete failed. Please try again later")
	}
	// 删除权限
	mapper.Delete[models.SysRoleMenu](mapper.BuilderUpdateWrapper(false).
		Eq(true, "role_id", id),
	)
	return utils.Success("")
}

func (that *SystemService) RoleUpdate(roleVO vo.RolePageVO) utils.R {
	msg := validateRole(roleVO)
	if msg != "ok" {
		return utils.Fail(constants.PARAMETER_ERROR, msg)
	}

	// 先删除旧权限
	mapper.Delete[models.SysRoleMenu](mapper.BuilderUpdateWrapper(false).
		Eq(true, "role_id", roleVO.Id),
	)

	// 在添加新权限
	for _, menuIds := range roleVO.MenuIds {
		var roleMenu models.SysRoleMenu
		roleMenu.RoleId = roleVO.Id
		roleMenu.MenuId = menuIds
		mapper.Insert[models.SysRoleMenu](&roleMenu, nil)
	}
	// 更新角色信息
	mapper.Update[models.SysRole](mapper.BuilderUpdateWrapper(false).
		Eq(true, "id", roleVO.Id).
		Set(true, "role_name", roleVO.RoleName).
		Set(true, "last_modified_date", keqing.NowDateStr()),
	)
	return utils.Success("")
}

func (that *SystemService) GetRoleById(id int) utils.R {
	roles := mapper.SelectById[models.SysRole](id)
	if keqing.IsEmpty(roles) {
		return utils.Fail(constants.DATA_NOT_EXIST, "The role does not exist. Please refresh the page and try again")
	}
	roleMenuList := mapper.SelectList[models.SysRoleMenu](mapper.BuilderQueryWrapper().
		Eq(true, "role_id", id),
	)

	var rolePage vo.RolePageVO

	rolePage.Id = roles[0].Id
	rolePage.RoleName = roles[0].RoleName
	var menuIds []int32

	for _, v := range roleMenuList {
		menuIds = append(menuIds, v.MenuId)
	}
	rolePage.MenuIds = menuIds
	return utils.Success(rolePage)
}

func (that *SystemService) RoleInsert(roleVO vo.RolePageVO) utils.R {
	msg := validateRole(roleVO)
	if msg != "ok" {
		return utils.Fail(constants.PARAMETER_ERROR, msg)
	}
	var role models.SysRole
	role.RoleName = roleVO.RoleName
	role.Deleted = constants.NO_DELETE_CODE
	role.CreatedDate = keqing.NowDateStr()
	role.LastModifiedDate = keqing.NowDateStr()

	count := mapper.Insert[models.SysRole](&role, nil)
	if count == 0 {
		return utils.Fail(constants.INSERT_ERROR, "Insert failed")
	}
	for _, menuIds := range roleVO.MenuIds {
		var roleMenu models.SysRoleMenu
		roleMenu.RoleId = role.Id
		roleMenu.MenuId = menuIds
		mapper.Insert[models.SysRoleMenu](&roleMenu, nil)
	}
	return utils.Success("")
}

func validateRole(role vo.RolePageVO) string {
	if keqing.IsEmpty(role.RoleName) {
		return "Role name cannot be empty"
	}
	count := mapper.SelectCount[models.SysRole](mapper.BuilderQueryWrapper().
		Eq(true, "role_name", role.RoleName).
		Ne(role.Id > 0, "id", role.Id).
		Eq(true, "deleted", constants.NO_DELETE_CODE),
	)

	if count > 0 {
		return "Role name already exists"
	}
	return "ok"
}

func getChildMenu(pvos []vo.MenuVO) {
	if len(pvos) == 0 {
		return
	}
	for i := 0; i < len(pvos); i++ {
		pid := pvos[i].Id

		list := mapper.SelectList[models.SysMenu](mapper.BuilderQueryWrapper().
			Eq(true, "parent_id", pid).
			OrderByAsc(true, "order_num"),
		)
		if keqing.IsEmpty(list) {
			continue
		}
		var cvos []vo.MenuChildVO
		for j := 0; j < len(list); j++ {
			var tmp vo.MenuChildVO
			tmp.Id = list[j].Id
			tmp.MenuName = list[j].MenuName
			tmp.RouterName = list[j].RouterName
			tmp.ParentId = list[j].ParentId
			tmp.MenuType = list[j].MenuType
			tmp.IsSys = list[j].IsSys
			tmp.OrderNum = list[j].OrderNum
			cvos = append(cvos, tmp)
		}
		pvos[i].Children = cvos
	}

}

func getParentMenu() []vo.MenuVO {
	var vos []vo.MenuVO
	list := mapper.SelectList[models.SysMenu](mapper.BuilderQueryWrapper().
		Eq(true, "parent_id", 0).
		OrderByAsc(true, "order_num"),
	)

	if keqing.IsEmpty(list) {
		return vos
	}
	for i := 0; i < len(list); i++ {
		var tmp vo.MenuVO
		tmp.Id = list[i].Id
		tmp.MenuName = list[i].MenuName
		tmp.RouterName = list[i].RouterName
		tmp.ParentId = list[i].ParentId
		tmp.MenuType = list[i].MenuType
		tmp.IsSys = list[i].IsSys
		tmp.OrderNum = list[i].OrderNum
		vos = append(vos, tmp)
	}
	return vos
}

func validateUser(userVO vo.UserVO) string {
	if keqing.IsEmpty(userVO.Username) {
		return "Username cannot be empty"
	}
	if userVO.Role <= 0 {
		return "Please select Role"
	}

	count := mapper.SelectCount[models.SysRole](mapper.BuilderQueryWrapper().
		Eq(true, "id", userVO.Role).
		Eq(true, "deleted", constants.NO_DELETE_CODE),
	)
	if count == 0 {
		return "Role does not exist"
	}

	if userVO.Id == 0 {
		if keqing.IsEmpty(userVO.Password) {
			return "Password cannot be empty"
		}
	}
	return "ok"
}
