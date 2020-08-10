package myJwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

// MyClaims 自定义声明结构体并内嵌 jwt.StandardClaims
// jwt 包自带的 jwtStandardClaims 只包含了官方字段
// 这里需要额外记录一个 UserSn 字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserSn uint64 `json:"userSn"`
	jwt.StandardClaims
}

// GenToken 生成 JWT
func GenToken(userSn uint64, args ...int) (string, error) {
	// 创建一个自定义的声明
	c := MyClaims{
		userSn, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Second * time.Duration(viper.GetInt("jwt.expire_seconds"))).Unix(), // 过期时间
			Issuer: viper.GetString("app.name"), // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的 secret 签名并获得完成的编码后的字符串 token
	return token.SignedString([]byte(viper.GetString("jwt.secret")))
}

// ParseToken 解析 JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{},
		func(token *jwt.Token) (i interface{}, err error) {
			return []byte(viper.GetString("jwt.secret")), nil
		})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
