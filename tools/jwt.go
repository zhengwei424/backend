package tools

import (
	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	UserID   uint64
	Username string
	jwt.StandardClaims
}

func (claims MyClaims) Valid() (err error) {
	return err
}

// 用户生成token的key
var key = []byte("kubernetes")

// GenerateToken 生成token
func (claims MyClaims) GenerateToken() (token string, err error) {
	result := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 可以通过自己的自定义key来生成token
	token, err = result.SignedString(key)
	return
}

// ParseToken 解析token
func ParseToken(generateToken string) (token *jwt.Token, claims MyClaims, err error) {
	token, err = jwt.ParseWithClaims(generateToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	//if claims, ok := token.Claims.(*Claims); ok && token.Valid {
	//	fmt.Printf("%v %v", claims.Username, claims.StandardClaims.ExpiresAt)
	//} else {
	//	fmt.Println(err)
	//}
	return token, claims, err
}
