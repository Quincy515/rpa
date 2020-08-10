package dao

import (
	"database/sql"
	"go.uber.org/zap"
	"yiran/model"
)

type userDAO struct{}

var DefaultUser = userDAO{}

// 通过邮箱注册保存用户表、用户统计表
// [golang 中的transaction（事务）的用法](https://studygolang.com/topics/1472)
func (u userDAO) EmailRegister(user *model.Account) (err error) {
	// 涉及处理多表，需要使用事务
	tx, err := db.Begin()
	if err != nil {
		zap.L().Error("mysql transaction failed", zap.Error(err))
		return
	}
	defer ClearTransaction(tx) // 如果出现异常情况,导致没有 commit和rollback,可以用来收尾

	// 1. 先判断当前用户是否已经被注册过
	if u.checkByEmail(user.Email) {
		return ErrorUserExit
	}

	// 2. 入库
	sqlStr := `INSERT INTO bt_user (user_sn, email, store_passwd) values (?, ?, ?)`
	if _, err = tx.Exec(sqlStr, user.UserSn, user.Email, user.StorePasswd); err != nil {
		zap.L().Error("insert bt_user failed", zap.Error(err))
		return
	}

	// 3. 新增用户统计表
	userCount := &model.UserCount{
		UserSn: user.UserSn,
	}
	sqlStr = `INSERT INTO bt_user_count (user_sn) VALUES (?)`
	if _, err = tx.Exec(sqlStr, userCount.UserSn); err != nil {
		zap.L().Error("insert bt_user_count failed", zap.Error(err))
		return
	}

	// 5. 事务提交
	if errTx := tx.Commit(); errTx != nil {
		zap.L().Error("commit view failed", zap.Any("user", user), zap.Error(err))
		return errTx
	}
	return
}

// 保存到用户文件表
func (u userDAO) SaveFile(fileMeta *model.FileMeta) (err error) {
	sqlStr := `INSERT IGNORE INTO bt_user_file (user_sn, file_sha1, file_size, file_name, upload_at, status) VALUES (?,?,?,?,?,1)`
	_, err = db.Exec(sqlStr, fileMeta.UserSn, fileMeta.FileSha1, fileMeta.FileSize, fileMeta.FileName, fileMeta.UploadAt)
	if err != nil {
		zap.L().Error("insert user_file failed", zap.Error(err))
		return
	}
	return
}

// 根据邮箱判断用户是否已经被注册过
func (u userDAO) checkByEmail(email string) bool {
	sqlStr := "SELECT COUNT(id) FROM bt_user WHERE email = ?"
	var count int64
	err := db.Get(&count, sqlStr, email)
	if err != nil && err != sql.ErrNoRows {
		zap.L().Error("by email query user exist failed", zap.Error(err))
		return true
	}
	return count > 0
}

// 根据邮箱获取用户信息
func (u userDAO) GetUserByEmail(email string) (user *model.Account, err error) {
	user = &model.Account{}
	sqlStr := `SELECT user_sn, email, store_passwd, nickname, avatar, gender, introduce, state, is_root, created_at  FROM bt_user WHERE email = ? LIMIT 1`
	err = db.Get(user, sqlStr, email)
	if err != nil && err != sql.ErrNoRows {
		// 查询数据库出错
		zap.L().Error("by email query user info failed", zap.Error(err))
		return
	}
	if err == sql.ErrNoRows {
		// 用户不存在
		zap.L().Error("by email query user info not exit", zap.Error(err))
		return nil, ErrorUserNotExit
	}
	return
}

// 根据用户编号获取用户信息
func (u userDAO) GetUserBySn(sn string) (user *model.Account, err error) {
	user = &model.Account{}
	sqlStr := `SELECT user_sn, email, store_passwd, nickname, avatar, gender, introduce, state, is_root, created_at  FROM bt_user WHERE user_sn = ? LIMIT 1`
	err = db.Get(user, sqlStr, sn)
	if err != nil && err != sql.ErrNoRows {
		// 查询数据库出错
		zap.L().Error("by sn query user info failed", zap.Error(err))
		return
	}
	if err == sql.ErrNoRows {
		// 用户不存在
		zap.L().Error("by sn query user info not exit", zap.Error(err))
		return nil, ErrorUserNotExit
	}
	return
}

// 根据 openid 查询第三方登录绑定关系
func (u userDAO) GetOauthByOpenid(openid string, minType int) (bind *model.OauthMemberBind, err error) {
	bind = &model.OauthMemberBind{}
	sqlStr := `SELECT user_sn, type, openid, unionid, extra FROM oauth_member_bind WHERE openid = ? AND type = ? LIMIT 1`
	err = db.Get(bind, sqlStr, openid, minType)
	//if err != nil && err != sql.ErrNoRows { // no rows in result set 没有找到在逻辑层处理
	if err != nil { // 查询不到对应记录，返回参数及错误均为nil
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			zap.L().Error("by openid query member bind failed", zap.Error(err))
			return
		}
	}
	return
}

// 通过 openid 进行注册保存用户表、用户统计表、用户和第三方登录关系绑定表
func (u userDAO) MemberRegister(member *model.OauthMemberBind) (err error) {
	// 涉及处理多表，需要使用事务
	tx, err := db.Begin()
	if err != nil {
		zap.L().Error("mysql transaction failed", zap.Error(err))
		return
	}
	defer ClearTransaction(tx) // 如果出现异常情况,导致没有 commit和rollback,可以用来收尾

	// 1. 保存到用户表 状态0-第三方登录绑定用户;1-正常注册;2-冻结账户
	sqlStr := `INSERT INTO bt_user (user_sn, state) values (?, 0)`
	if _, err = tx.Exec(sqlStr, member.UserSn); err != nil {
		zap.L().Error("insert bt_user failed", zap.Error(err))
		return
	}

	// 2. 新增用户统计表
	userCount := &model.UserCount{
		UserSn: member.UserSn,
	}
	sqlStr = `INSERT INTO bt_user_count (user_sn) VALUES (?)`
	if _, err = tx.Exec(sqlStr, userCount.UserSn); err != nil {
		zap.L().Error("insert bt_user_count failed", zap.Error(err))
		return
	}

	// 3. 将创建成功的成员表和第三方登录绑定关系表进行绑定
	sqlStr = `INSERT INTO oauth_member_bind (user_sn, type, openid, unionid, extra) VALUES (?, ?, ?, ?, ?)`
	if _, err = tx.Exec(sqlStr, userCount.UserSn, member.Type, member.Openid, member.Unionid, member.Extra); err != nil {
		zap.L().Error("insert bt_user_count failed", zap.Error(err))
		return
	}

	// 事务提交
	if errTx := tx.Commit(); errTx != nil {
		zap.L().Error("commit transaction failed", zap.Any("member", member), zap.Error(err))
		return errTx
	}
	return
}

// 根据微信小程序用户授权更新用户信息
func (u userDAO) WxUpdateUser(user *model.Account) (err error) {
	sqlStr := `UPDATE bt_user SET nickname = ?, avatar = ?, gender = ? WHERE user_sn = ? AND state = 0`
	_, err = db.Exec(sqlStr, user.Nickname, user.Avatar, user.Gender, user.UserSn)
	if err != nil {
		zap.L().Error("WxUpdateUser failed", zap.Error(err))
		return
	}
	return
}
