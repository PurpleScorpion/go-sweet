package service

import (
	"github.com/PurpleScorpion/go-sweet-json/jsonutil"
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"github.com/PurpleScorpion/go-sweet-orm/mapper"
	"sweet-common/constants"
	"sweet-common/utils"
	"sweet-common/vo"
	"sweet-src/main/golang/models"
	"time"
)

type UserService struct {
}

var USER_SALT = "babalachongya"

func (that *UserService) HealthCheck(id int32) utils.R {

	expire := utils.GetCache(constants.GetUserExpireTimeKey(id))
	if keqing.IsEmpty(expire) {
		return utils.Fail(constants.TOKEN_ERROR, "user expired")
	}

	utcDate := keqing.ParseUTC(expire)
	now := time.Now().UTC()
	if now.After(utcDate) {
		return utils.Fail(constants.TOKEN_ERROR, "user expired")
	}

	utils.SetCache(constants.GetHealthCheckKey(id), keqing.NowUTCDateStr())

	subMinutes := utcDate.Sub(now).Minutes()

	if subMinutes <= 30 && subMinutes > 0 {
		js := jsonutil.NewJSONObject()
		js.FluentPut("id", id)
		// 计算过期时间
		expireTime := keqing.NowDate().Add(4 * time.Hour)
		newUtcDate := keqing.FormatDate(expireTime.UTC(), keqing.DEFAULT_UTC_FORMAT)
		js.FluentPut("expirationTime", newUtcDate)
		token := keqing.RsaEncrypt(js.ToJsonString())
		utils.SetCache(constants.GetUserExpireTimeKey(id), newUtcDate)
		return utils.Success(token)
	}

	return utils.Success("")
}

func (that *UserService) RePassword(userVO vo.UserVO) utils.R {
	var list []models.User
	mapper.SelectById(&list, userVO.Id)
	if len(list) == 0 {
		return utils.Fail(constants.USER_EMPTY_CODE, "user does not exist")
	}
	user := list[0]

	password := userVO.Password
	oldPassword := userVO.OldPassword

	oldPasswordMd5 := keqing.MD5Salt(oldPassword, USER_SALT)
	if oldPasswordMd5 != user.Password {
		return utils.Fail(constants.USER_OLD_PASSWORD_ERROR, "Old password error")
	}

	passwordMd5 := keqing.MD5Salt(password, USER_SALT)

	var u models.User
	uqw := mapper.BuilderQueryWrapper(&u)
	uqw.Eq(true, "id", user.Id)
	uqw.Set(true, "password", passwordMd5)
	count := mapper.Update(uqw)
	if count == 0 {
		return utils.Fail(constants.UPDATE_ERROR, "Password update failed, please try again later")
	}
	return utils.Success("")
}

func (that *UserService) Login(user models.User) utils.R {
	password := user.Password

	var list []models.User
	qw := mapper.BuilderQueryWrapper(&list)
	qw.Eq(keqing.IsNotEmpty(user.Username), "username", user.Username)
	//qw.Eq(true, "status", constants.NORMAL_STATUS)
	qw.Eq(true, "deleted", constants.NO_DELETE_CODE)

	mapper.SelectList(qw)

	if len(list) == 0 {
		return utils.Fail(constants.USER_EMPTY_CODE, "Incorrect username or password")
	}
	u := list[0]

	if u.Status == constants.FAIL_STATUS {
		return utils.Fail(constants.USER_FORBIDDEN_CODE, "The account has been frozen")
	}

	// 加盐加密
	hashedPassword := keqing.MD5Salt(password, USER_SALT)

	if hashedPassword != u.Password {
		return utils.Fail(constants.USER_EMPTY_CODE, "Incorrect username or password")
	}
	// 获取权限角色
	roleList := getRoleList(u)
	if len(roleList) == 0 {
		return utils.Fail(constants.ROLE_NOT_CONFIGURED, "Role not configured")
	}

	js := jsonutil.NewJSONObject()
	js.FluentPut("id", u.Id)

	// 计算过期时间
	expireTime := keqing.NowDate().Add(4 * time.Hour)
	utcDate := keqing.FormatDate(expireTime.UTC(), keqing.DEFAULT_UTC_FORMAT)
	js.FluentPut("expirationTime", utcDate)
	token := keqing.RsaEncrypt(js.ToJsonString())

	utils.SetCache(constants.GetHealthCheckKey(u.Id), keqing.NowUTCDateStr())
	utils.SetCache(constants.GetUserExpireTimeKey(u.Id), utcDate)

	var userVO vo.UserVO
	userVO.Id = u.Id
	userVO.Token = token
	userVO.Role = u.Role
	userVO.Username = u.Username
	userVO.Routers = roleList
	return utils.Success(userVO)
}

func getRoleList(u models.User) []string {
	var routers []string
	var menuList []models.SysMenu
	qw := mapper.BuilderQueryWrapper(&menuList)
	// 如果是超管 , 则直接返回所有的路由
	if u.Role == constants.ROOT_ROLE_ID {
		mapper.SelectList(qw)
		for i := 0; i < len(menuList); i++ {
			routers = append(routers, menuList[i].RouterName)
		}
		return routers
	}
	// 检查角色状态
	flag, _ := checkRole(u)
	if !flag {
		return routers
	}

	var roleMenuList []models.SysRoleMenu
	rmQw := mapper.BuilderQueryWrapper(&roleMenuList)
	rmQw.Eq(true, "role_id", u.Role)
	mapper.SelectList(rmQw)
	// 角色没有配置菜单
	if len(roleMenuList) == 0 {
		return routers
	}
	var ids []int32
	for i := 0; i < len(roleMenuList); i++ {
		ids = append(ids, roleMenuList[i].MenuId)
	}
	qw.InInt32(true, "id", ids)
	mapper.SelectList(qw)
	if len(menuList) == 0 {
		return routers
	}

	for i := 0; i < len(menuList); i++ {
		routers = append(routers, menuList[i].RouterName)
	}

	return routers
}

/*
检查角色状态
true : 角色正常
false : 角色已删除或不存在
*/
func checkRole(u models.User) (bool, []models.SysRole) {
	var roleList []models.SysRole
	// 没配置角色
	if u.Role == 0 {
		return false, roleList
	}
	// 检查角色状态

	mapper.SelectById(&roleList, u.Role)
	if len(roleList) == 0 {
		return false, roleList
	}
	// 已删除
	if roleList[0].Deleted == constants.DELETE_CODE {
		return false, roleList
	}

	return true, roleList
}
