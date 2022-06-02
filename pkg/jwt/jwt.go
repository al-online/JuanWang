package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个UserID字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

//定义Secret
var mySecret = []byte("代码敲栏,年薪过万")

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return mySecret, nil
}

// TokenExpireDuration 定义JWT的过期时间
const TokenExpireDuration = time.Hour * 2

// GenToken 生成access token 和 refresh token
func GenToken(userID int64, username string) (Token string, err error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		userID,   // 自定义字段
		username, // 自定义字段
		jwt.StandardClaims{ // JWT规定的7个官方字段
			ExpiresAt: time.Now().Add(
				time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(), // 过期时间
			Issuer: "bluebell", // 签发人
		},
	}
	// 加密并获得完整的编码后的字符串token
	Token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)
	return
}

func ParseToken(tokenString string) (claims *MyClaims, err error) {
	// 解析token
	var token *jwt.Token
	claims = new(MyClaims)
	token, err = jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return
	}
	if !token.Valid { // 校验token
		err = errors.New("invalid token")
	}
	return
}
