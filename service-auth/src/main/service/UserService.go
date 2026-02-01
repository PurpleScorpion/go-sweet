package service

import (
	"shared/constants"
	"shared/utils"

	"github.com/PurpleScorpion/go-sweet-json/jsonutil"
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"github.com/PurpleScorpion/go-sweet-orm/v3/mapper"

	"service-common/src/main/models"
	"shared/vo"
	"time"
)

type UserService struct {
}

var USER_SALT = "babalachongya"

func (that *UserService) RePassword(userVO vo.UserVO) utils.R {
	list := mapper.SelectById[models.User](userVO.Id)

	if keqing.IsEmpty(list) {
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

	count := mapper.Update[models.User](mapper.BuilderUpdateWrapper(false).
		Eq(true, "id", user.Id).
		Set(true, "password", passwordMd5),
	)
	if count == 0 {
		return utils.Fail(constants.UPDATE_ERROR, "Password update failed, please try again later")
	}
	return utils.Success("")
}

func (that *UserService) Login(user models.User) utils.R {

	if keqing.IsEmpty(user.Username) {
		return utils.Fail(constants.USER_EMPTY_CODE, "Please enter the correct username")
	}

	password := user.Password

	list := mapper.SelectList[models.User](mapper.BuilderQueryWrapper().
		Eq(true, "username", user.Username).
		Eq(true, "deleted", constants.NO_DELETE_CODE),
	)

	if keqing.IsEmpty(list) {
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
	routers := keqing.NewSet[string](nil)
	qw := mapper.BuilderQueryWrapper()

	// 如果是超管 , 则直接返回所有的路由
	if u.Role == constants.ROOT_ROLE_ID {
		menuList := mapper.SelectList[models.SysMenu](qw)

		for i := 0; i < len(menuList); i++ {
			routers.Add(menuList[i].RouterName)
		}
		return routers.GetData()
	}
	// 检查角色状态
	flag, _ := checkRole(u)
	if !flag {
		return routers.GetData()
	}

	roleMenuList := mapper.SelectList[models.SysRoleMenu](mapper.BuilderQueryWrapper().
		Eq(true, "role_id", u.Role),
	)
	// 角色没有配置菜单
	if keqing.IsEmpty(roleMenuList) {
		return routers.GetData()
	}
	var ids []int32
	for i := 0; i < len(roleMenuList); i++ {
		ids = append(ids, roleMenuList[i].MenuId)
	}
	qw.InInt32(true, "id", ids)
	menuList := mapper.SelectList[models.SysMenu](qw)
	if keqing.IsEmpty(menuList) {
		return routers.GetData()
	}

	for i := 0; i < len(menuList); i++ {
		routers.Add(menuList[i].RouterName)
	}

	return routers.GetData()
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

	mapper.SelectById[models.SysRole](u.Role)
	if len(roleList) == 0 {
		return false, roleList
	}
	// 已删除
	if roleList[0].Deleted == constants.DELETE_CODE {
		return false, roleList
	}

	return true, roleList
}
