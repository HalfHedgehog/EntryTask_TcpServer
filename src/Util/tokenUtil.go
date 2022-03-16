package Util

import (
	"TcpServer/src/global"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type Token struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

// CreateToken 创建token
func CreateToken(userId int64) string {
	//加密
	c := Token{
		UserId: strconv.FormatInt(userId, 10),
		StandardClaims: jwt.StandardClaims{
			//什么时候生效，现在生效
			NotBefore: time.Now().Unix(),
			//什么时候失效，半个小时后
			ExpiresAt: time.Now().Unix() + global.Config.JWT.ExpiresTime,
			//签发人
			Issuer: global.Config.JWT.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	fmt.Println(token)
	newToken, e := token.SignedString([]byte(global.Config.JWT.Key))
	if e != nil {
		fmt.Println(e)
		return ""
	}
	return newToken
}
