package logic

import (
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	dao "yiran/dao/mysql"
	"yiran/model"
	"yiran/pkg/snowflake"
	"yiran/pkg/util"
)

type userLogic struct{}

var DefaultUserLogic = userLogic{}

// 处理微信小程序注册登录逻辑
func (u userLogic) MiniLogin(addAccount *model.MiniAccount) (account *model.Account, err error) {
	account = &model.Account{}
	// 1. 根据 code 获取 openid
	weChatId, err := util.GetWeChatId(addAccount.Code)
	if err != nil {
		zap.L().Error("util.GetWeChatId failed", zap.Error(err))
		return
	}
	zap.L().Debug("1.获取openid", zap.Any("weChatId", weChatId))

	// 2. 根据 openid 判断是否已经注册过
	bindInfo, err := dao.DefaultUser.GetOauthByOpenid(weChatId.Openid, 1)
	if err != nil {
		zap.L().Error("dao.DefaultUser.GetOauthByOpenid", zap.Error(err))
		return
	}
	zap.L().Debug("2.获取bindInfo", zap.Any("bindInfo", bindInfo))

	// 3. 未注册就进行注册
	if bindInfo == nil {
		// 创建
		member := &model.OauthMemberBind{
			Type:   1,
			Openid: weChatId.Openid,
			Extra:  "微信小程序注册",
		}
		// 1. 生成用户编号
		userID, err1 := snowflake.GenID()
		if err1 != nil {
			zap.L().Error("snowflake.GenID failed", zap.Error(err))
			return nil, err1
		}
		member.UserSn = userID
		// 2. 保存到用户表、用户统计表、用户和第三方登录绑定表
		err2 := dao.DefaultUser.MemberRegister(member)
		if err2 != nil {
			zap.L().Error("dao.DefaultUser.MemberRegister", zap.Error(err))
			return nil, err2
		}
		account.UserSn = member.UserSn // 保存注册成功的用户编号
	} else {
		account.UserSn = bindInfo.UserSn // 保存查询到的用户编号
	}
	// 4. 已注册就返回用户信息
	return
}

// 通过邮箱登录的逻辑处理
func (u userLogic) EmailLogin(account *model.AccountByEmail) (user *model.Account, err error) {
	user, err = dao.DefaultUser.GetUserByEmail(account.Email)
	if err != nil {
		return
	}
	// 4. 校验密码
	return user, bcrypt.CompareHashAndPassword(user.StorePasswd, []byte(account.Password))
}

// 通过邮箱注册的逻辑处理
func (u userLogic) EmailRegister(account model.AccountByEmail) (user *model.Account, err error) {
	// 1. 构造用户结构体
	user = &model.Account{
		Email:    account.Email,
		Password: account.Password,
	}
	// 2. 使用 bcrypt 进行密码加密
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("bcrypt.GenerateFromPassword failed", zap.Error(err))
		return
	}
	// 3. 生成用户编号
	userID, err := snowflake.GenID()
	if err != nil {
		zap.L().Error("snowflake.GenID failed", zap.Error(err))
		return
	}
	user.StorePasswd = password
	user.UserSn = userID
	// 4. 保存到数据库
	if err = dao.DefaultUser.EmailRegister(user); err != nil {
		zap.L().Error("userLogic email register failed", zap.Error(err))
		return
	}
	return
}

// 根据微信用户授权更新用户信息
func (u userLogic) WxUpdateUser(account model.MiniAccount, userSn uint64) (user *model.Account, err error) {
	// 1. 构造用户结构体
	user = &model.Account{
		UserSn:   userSn,
		Nickname: account.UserInfo.NickName,
		Avatar:   account.UserInfo.AvatarURL,
		Gender:   uint8(account.UserInfo.Gender),
	}
	// 2. 根据 userSn 判断是否已经注册过
	_, err = dao.DefaultUser.GetUserBySn(strconv.Itoa(int(userSn)))
	if err != nil {
		zap.L().Error("dao.DefaultUser.GetOauthByOpenid", zap.Error(err))
		return
	}
	// 3. 更新用户信息
	err = dao.DefaultUser.WxUpdateUser(user)
	if err != nil {
		zap.L().Error("dao.DefaultUser.UpdateUser failed", zap.Error(err))
		return
	}
	return
}
