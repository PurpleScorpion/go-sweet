package service

import (
	"github.com/PurpleScorpion/go-sweet-orm/mapper"
	"sweet-common/constants"
	"sweet-common/utils"
	"sweet-common/vo"
	"sweet-src/main/golang/models"
)

type SystemService struct {
}

func (that *SystemService) UserPageData(userVO vo.UserPageVO) utils.R {
	var list []models.User
	qw := mapper.BuilderQueryWrapper(&list)
	qw.Like(utils.IsNotEmpty(userVO.Username), "username", userVO.Username)
	qw.Ne(true, "role", constants.ROOT_ROLE_ID)
	qw.Eq(true, "deleted", constants.NO_DELETE_CODE)
	page := mapper.BuilderPageUtils(userVO.Current, userVO.PageSize, qw)
	pageData := mapper.Page(page)

	arr := *pageData.List.(*[]models.User)
	var vos []vo.UserVO
	for i := 0; i < len(arr); i++ {
		var tmp = vo.UserVO{}
		tmp.Id = arr[i].Id
		tmp.Username = arr[i].Username
		tmp.Role = arr[i].Role
		tmp.LastModifiedDate = arr[i].LastModifiedDate
		tmp.Status = arr[i].Status
		tmp.RoleName = ""
		var role []models.SysRole
		mapper.SelectById(&role, arr[i].Role)
		if len(role) > 0 {
			if role[0].Deleted == constants.NO_DELETE_CODE {
				tmp.RoleName = role[0].RoleName
			}
		}
		vos = append(vos, tmp)
	}
	pageData.List = vos
	return utils.Success(pageData)
}

func (that *SystemService) RolePageData(roleVO vo.RolePageVO) utils.R {
	var list []models.SysRole
	qw := mapper.BuilderQueryWrapper(&list)
	qw.Like(utils.IsNotEmpty(roleVO.RoleName), "role_name", roleVO.RoleName)
	qw.Eq(true, "deleted", constants.NO_DELETE_CODE)
	qw.OrderByTimeAsc(true, "created_by")
	page := mapper.BuilderPageUtils(roleVO.Current, roleVO.PageSize, qw)
	pageData := mapper.Page(page)
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

func (that *SystemService) DeleteUserById(id int32, userId int32) utils.R {
	sqw := mapper.BuilderQueryWrapper(&models.User{})
	sqw.Eq(true, "id", id)
	sqw.Eq(true, "deleted", constants.NO_DELETE_CODE)
	count := mapper.SelectCount(sqw)
	if count == 0 {
		return utils.Fail(constants.DATA_NOT_EXIST, "The user does not exist")
	}

	qw := mapper.BuilderQueryWrapper(&models.User{})
	qw.Eq(true, "id", id)
	qw.Set(true, "last_modified_date", utils.GetNowDate())
	qw.Set(true, "last_modified_by", userId)
	qw.Set(true, "deleted", constants.DELETE_CODE)
	count = mapper.Update(qw)
	if count == 0 {
		return utils.Fail(constants.UPDATE_ERROR, "Update failed")
	}
	return utils.Success("")
}

func (that *SystemService) ChangeUserStatus(userVO vo.UserVO) utils.R {
	sqw := mapper.BuilderQueryWrapper(&models.User{})
	sqw.Eq(true, "id", userVO.Id)
	sqw.Eq(true, "deleted", constants.NO_DELETE_CODE)
	count := mapper.SelectCount(sqw)
	if count == 0 {
		return utils.Fail(constants.DATA_NOT_EXIST, "The user does not exist")
	}
	qw := mapper.BuilderQueryWrapper(&models.User{})
	qw.Eq(true, "id", userVO.Id)
	qw.Set(true, "last_modified_date", utils.GetNowDate())
	qw.Set(true, "last_modified_by", userVO.LastModifiedBy)
	qw.Set(true, "status", userVO.Status)
	count = mapper.Update(qw)
	if count == 0 {
		return utils.Fail(constants.UPDATE_ERROR, "Update failed")
	}
	if userVO.Status == constants.FAIL_STATUS {
		utils.DeleteCache(constants.GetUserExpireTimeKey(userVO.Id))
	}

	return utils.Success(constants.SUCCESS)
}

func (that *SystemService) UserUpdate(userVO vo.UserVO) utils.R {
	r := validateUser(userVO)
	if r.Code != 200 {
		return r
	}
	sqw := mapper.BuilderQueryWrapper(&models.User{})
	sqw.Eq(true, "id", userVO.Id)
	sqw.Eq(true, "deleted", constants.NO_DELETE_CODE)
	count := mapper.SelectCount(sqw)
	if count == 0 {
		return utils.Fail(constants.DATA_NOT_EXIST, "The user does not exist")
	}

	qw := mapper.BuilderQueryWrapper(&models.User{})
	qw.Eq(true, "id", userVO.Id)
	qw.Set(utils.IsNotEmpty(userVO.Password), "password", utils.GetMD5(userVO.Password, USER_SALT))
	qw.Set(true, "role", userVO.Role)
	qw.Set(true, "last_modified_date", utils.GetNowDate())
	qw.Set(true, "last_modified_by", userVO.LastModifiedBy)
	count = mapper.Update(qw)
	if count == 0 {
		return utils.Fail(constants.UPDATE_ERROR, "Update failed")
	}

	return utils.Success("")
}

func (that *SystemService) UserInsert(userVO vo.UserVO) utils.R {
	r := validateUser(userVO)
	if r.Code != 200 {
		return r
	}

	var users []models.User
	// 判断用户名是否重复
	qw := mapper.BuilderQueryWrapper(&users)
	qw.Eq(true, "username", userVO.Username)
	mapper.SelectList(qw)
	if len(users) > 0 {
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
	user.Password = utils.GetMD5(userVO.Password, USER_SALT)
	user.Role = userVO.Role
	user.CreatedDate = utils.GetNowDate()
	user.CreatedBy = userVO.CreatedBy
	user.LastModifiedBy = userVO.LastModifiedBy
	user.LastModifiedDate = utils.GetNowDate()
	user.Deleted = constants.NO_DELETE_CODE
	user.Status = constants.NORMAL_STATUS
	count := mapper.InsertCustom(&user, true, false)
	if count == 0 {
		return utils.Fail(constants.INSERT_ERROR, "Insert failed")
	}

	return utils.Success("")
}

/*
恢复账号删除状态
*/
func recoveryUser(userVO vo.UserVO) utils.R {
	qw := mapper.BuilderQueryWrapper(&models.User{})
	qw.Eq(true, "id", userVO.Id)
	qw.Set(utils.IsNotEmpty(userVO.Password), "password", utils.GetMD5(userVO.Password, USER_SALT))
	qw.Set(true, "role", userVO.Role)
	qw.Set(true, "last_modified_date", utils.GetNowDate())
	qw.Set(true, "last_modified_by", userVO.LastModifiedBy)
	qw.Set(true, "deleted", constants.NO_DELETE_CODE)
	count := mapper.Update(qw)
	if count == 0 {
		return utils.Fail(constants.UPDATE_ERROR, "Update failed")
	}

	return utils.Success("")
}

func (that *SystemService) GetUserById(id int32) utils.R {
	var users []models.User
	mapper.SelectById(&users, id)

	if len(users) == 0 {
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
	count := mapper.InsertCustom(&menu, true, false)
	if count == 0 {
		return utils.Fail(constants.INSERT_ERROR, "Insert failed")
	}
	return utils.Success("")
}

func (that *SystemService) MenuUpdate(menu models.SysMenu) utils.R {
	qw := mapper.BuilderQueryWrapper(&models.SysMenu{})
	qw.Eq(true, "id", menu.Id)
	qw.Set(utils.IsNotEmpty(menu.MenuName), "menu_name", menu.MenuName)
	qw.Set(utils.IsNotEmpty(menu.RouterName), "router_name", menu.RouterName)
	qw.Set(menu.MenuType > 0, "menu_type", menu.MenuType)
	qw.Set(true, "order_num", menu.OrderNum)
	qw.Set(true, "parent_id", menu.ParentId)
	count := mapper.Update(qw)
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
	var roles []models.SysRole
	qw := mapper.BuilderQueryWrapper(&roles)
	qw.Eq(true, "deleted", constants.NO_DELETE_CODE)
	mapper.SelectList(qw)
	return utils.Success(roles)
}

func (that *SystemService) GetMenuById(id int32) utils.R {
	var menus []models.SysMenu
	mapper.SelectById(&menus, id)
	if len(menus) == 0 {
		return utils.Fail(constants.DATA_NOT_EXIST, "Menu does not exist")
	}
	return utils.Success(menus[0])
}

func (that *SystemService) DeleteMenuById(id int32) utils.R {
	var menus []models.SysMenu
	mapper.SelectById(&menus, id)
	if len(menus) == 0 {
		return utils.Fail(constants.DATA_NOT_EXIST, "Menu does not exist")
	}
	// 判断是否是父节点 , 且父节点下有未删除的子节点

	qw := mapper.BuilderQueryWrapper(&models.SysMenu{})
	qw.Eq(true, "parent_id", id)
	count := mapper.SelectCount(qw)
	if count > 0 {
		return utils.Fail(constants.DELETE_ERROR, "There are submenus in the current directory that have not been deleted")
	}
	count = mapper.DeleteById(&models.SysMenu{}, id)
	if count == 0 {
		return utils.Fail(constants.DELETE_ERROR, "Delete failed, please try again later")
	}
	return utils.Success("")
}

func (that *SystemService) DeleteRoleById(id int32) utils.R {
	qw := mapper.BuilderQueryWrapper(&models.SysRole{})
	qw.Eq(true, "id", id)
	count := mapper.SelectCount(qw)
	if count == 0 {
		return utils.Fail(constants.DATA_NOT_EXIST, "Menu does not exist")
	}
	count = mapper.DeleteById(&models.SysRole{}, id)
	if count == 0 {
		return utils.Fail(constants.DELETE_ERROR, "Delete failed. Please try again later")
	}
	// 删除权限
	delQw := mapper.BuilderQueryWrapper(&models.SysRoleMenu{})
	delQw.Eq(true, "role_id", id)
	mapper.Delete(delQw)
	return utils.Success("")
}

func (that *SystemService) RoleUpdate(roleVO vo.RolePageVO) utils.R {
	r := validateRole(roleVO)
	if r.Code != 200 {
		return r
	}

	// 先删除旧权限
	delQw := mapper.BuilderQueryWrapper(&models.SysRoleMenu{})
	delQw.Eq(true, "role_id", roleVO.Id)
	mapper.Delete(delQw)
	// 在添加新权限
	for _, menuIds := range roleVO.MenuIds {
		var roleMenu models.SysRoleMenu
		roleMenu.RoleId = roleVO.Id
		roleMenu.MenuId = menuIds
		mapper.Insert(&roleMenu)
	}
	// 更新角色信息
	qw := mapper.BuilderQueryWrapper(&models.SysRole{})
	qw.Eq(true, "id", roleVO.Id)
	qw.Set(true, "role_name", roleVO.RoleName)
	qw.Set(true, "last_modified_date", utils.GetNowDate())
	qw.Set(true, "last_modified_by", roleVO.UserId)
	mapper.Update(qw)
	return utils.Success("")
}

func (that *SystemService) GetRoleById(id int32) utils.R {
	var roles []models.SysRole
	mapper.SelectById(&roles, id)
	if len(roles) == 0 {
		return utils.Fail(constants.DATA_NOT_EXIST, "The role does not exist. Please refresh the page and try again")
	}
	var roleMenuList []models.SysRoleMenu
	qw := mapper.BuilderQueryWrapper(&roleMenuList)
	qw.Eq(true, "role_id", id)
	mapper.SelectList(qw)
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
	r := validateRole(roleVO)
	if r.Code != 200 {
		return r
	}
	var role models.SysRole
	role.RoleName = roleVO.RoleName
	role.Deleted = constants.NO_DELETE_CODE
	role.CreatedDate = utils.GetNowDate()
	role.CreatedBy = roleVO.UserId
	role.LastModifiedDate = utils.GetNowDate()
	role.LastModifiedBy = roleVO.UserId

	count := mapper.InsertCustom(&role, true, false)
	if count == 0 {
		return utils.Fail(constants.INSERT_ERROR, "Insert failed")
	}
	for _, menuIds := range roleVO.MenuIds {
		var roleMenu models.SysRoleMenu
		roleMenu.RoleId = role.Id
		roleMenu.MenuId = menuIds
		mapper.Insert(&roleMenu)
	}
	return utils.Success("")
}

func validateRole(role vo.RolePageVO) utils.R {
	if utils.IsEmpty(role.RoleName) {
		return utils.Fail(constants.PARAMETER_ERROR, "Role name cannot be empty")
	}
	qw := mapper.BuilderQueryWrapper(&models.SysRole{})
	qw.Eq(true, "role_name", role.RoleName)
	qw.Ne(role.Id > 0, "id", role.Id)
	qw.Eq(true, "deleted", constants.NO_DELETE_CODE)
	count := mapper.SelectCount(qw)

	if count > 0 {
		return utils.Fail(constants.PARAMETER_ERROR, "Role name already exists")
	}
	return utils.Success("")
}

func getChildMenu(pvos []vo.MenuVO) {
	if len(pvos) == 0 {
		return
	}
	for i := 0; i < len(pvos); i++ {
		pid := pvos[i].Id
		var list []models.SysMenu
		var cvos []vo.MenuChildVO
		qw := mapper.BuilderQueryWrapper(&list)
		qw.Eq(true, "parent_id", pid)
		qw.OrderByAsc(true, "order_num")
		mapper.SelectList(qw)
		if len(list) == 0 {
			continue
		}
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
	var list []models.SysMenu
	var vos []vo.MenuVO
	qw := mapper.BuilderQueryWrapper(&list)
	qw.Eq(true, "parent_id", 0)
	qw.OrderByAsc(true, "order_num")
	mapper.SelectList(qw)
	if len(list) == 0 {
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

func validateUser(userVO vo.UserVO) utils.R {
	if utils.IsEmpty(userVO.Username) {
		return utils.Fail(constants.PARAMETER_ERROR, "Username cannot be empty")
	}
	if userVO.Role <= 0 {
		return utils.Fail(constants.PARAMETER_ERROR, "Please select Role")
	}
	qw := mapper.BuilderQueryWrapper(&models.SysRole{})
	qw.Eq(true, "id", userVO.Role)
	qw.Eq(true, "deleted", constants.NO_DELETE_CODE)
	count := mapper.SelectCount(qw)
	if count == 0 {
		return utils.Fail(constants.PARAMETER_ERROR, "Role does not exist")
	}

	if userVO.Id == 0 {
		if utils.IsEmpty(userVO.Password) {
			return utils.Fail(constants.PARAMETER_ERROR, "Password cannot be empty")
		}
	}
	return utils.Success("")
}
