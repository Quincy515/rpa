package model

import (
	"errors"
	"time"
)

// 用户表结构
type Account struct {
	ID          int       `json:"id" db:"id"`
	UserSn      uint64    `json:"userSn" db:"user_sn"`
	Nickname    string    `json:"nickname" db:"nickname"`
	Email       string    `json:"email" db:"email"`
	Password    string    `json:"password"`
	Avatar      string    `json:"avatar" db:"avatar"`
	Gender      uint8     `json:"gender" db:"gender"`
	Introduce   string    `json:"introduce" db:"introduce"`
	State       uint8     `json:"state" db:"state"`
	IsRoot      bool      `json:"isRoot" db:"is_root"`
	StorePasswd []byte    `json:"storePasswd" db:"store_passwd"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// 用户统计表结构
type UserCount struct {
	UserSn    uint64 `json:"userSn" db:"user_sn"`
	FansNum   int    `json:"fansNum" db:"fans_num"`
	FollowNum int    `json:"followNum" db:"follow_num"`
	PlanNum   int    `json:"planNum" db:"plan_num"`
	PlanDone  int    `json:"planDone" db:"plan_done"`
	ZanNum    int    `json:"zanNum" db:"zan_num"`
}

// 邮箱注册登录的请求结构体
type AccountByEmail struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// 邮箱注册的数据校验
func (a AccountByEmail) Validation() error {
	switch {
	case len(a.Email) == 0:
		return errors.New("缺少必填字段邮箱")
	case len(a.Password) == 0:
		return errors.New("缺少必填字段密码")
	case a.Password != a.ConfirmPassword:
		return errors.New("两次密码不一致")
	default:
		return nil
	}
}

// 微信小程序注册登录请求的结构体
type MiniAccount struct {
	Code     string `json:"js_code"`
	UserInfo struct {
		NickName  string `json:"nickName"`
		Gender    int    `json:"gender"`
		Language  string `json:"language"`
		City      string `json:"city"`
		Province  string `json:"province"`
		Country   string `json:"country"`
		AvatarURL string `json:"avatarURL"`
	}
	Ip     string `json:"ip"`
	Mobile string `json:"mobile"`
}

type Code2Session struct {
	Openid     string `json: "openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}
