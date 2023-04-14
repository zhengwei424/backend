package tools

import "golang.org/x/crypto/bcrypt"

// CryptPassword 加密
func CryptPassword(password string) (hashPassword string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// DeCryptPassword 判断密码正确性
func DeCryptPassword(hashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}