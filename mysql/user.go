package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/RaymondCode/simple-demo/config"
)

const secret = "chery"

func encryptPassword(data []byte) (result string) {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum(data))
}

// CheckUserExist 检验用户是否存在
func CheckUserExist(username string) (error error) {
	user := &config.UserForm{}
	res := Db.Where("username=?", username).Find(&user)
	//数据库查询失败
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected > 0 {
		return errors.New("the user already exists")
	}
	return nil
}

// CheckUserExistById 检验用户Id是否存在
func CheckUserExistById(userId int64) (error error) {
	user := &config.UserForm{}
	res := Db.Where("user_id =?", userId).Find(&user)
	//数据库查询失败
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected > 0 {
		return nil
	}
	return errors.New("the user does not exists")
}

// InsertUser 向数据库中插入数据
func InsertUser(user config.UserForm) (int64, error) {
	// 对密码进行加密
	user.Password = encryptPassword([]byte(user.Password))
	//插入数据库
	result := Db.Create(&user) // 通过数据的指针来创建
	return user.UserId, result.Error
}

//Login 登录处理
func Login(username string, password string) (error, int64) {
	user := &config.UserForm{}
	originPassword := password // 记录一下原始密码(用户登录的密码)
	res := Db.Where("username=?", username).Find(&user)
	//数据库查询失败
	if res.Error != nil {
		return res.Error, 0
	}
	//该用户不存在
	if res.RowsAffected == 0 {
		return errors.New("user does not exist"), 0
	}
	// 生成加密密码与查询到的密码比较
	password = encryptPassword([]byte(originPassword))
	if user.Password != password {
		return errors.New("the user name or password is incorrect"), 0
	}
	return nil, user.UserId
}
