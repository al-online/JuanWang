package mysql

import (
	"github.com/RaymondCode/simple-demo/config"
)

func InsertRelation(userId int64, toUserId int64) error {
	var err error
	//是否存在用户关注表如果不存在则创建
	err = Db.AutoMigrate(&config.Relation{})
	if err != nil {
		return err
	}
	relation := config.Relation{
		UserId:   userId,
		ToUserId: toUserId,
	}
	//添加记录
	res := Db.Create(&relation)
	if res.Error != nil {
		return res.Error
	}
	return err
}

func DeleteRelation(userId int64, toUserId int64) error {
	var err error
	relation := config.Relation{
		UserId:   userId,
		ToUserId: toUserId,
	}
	//删除记录
	res := Db.Where("user_id=? and to_user_id=?", userId, toUserId).Delete(&relation)
	if res.Error != nil {
		return res.Error
	}
	return err
}

func GetFollowList(user_id int64) ([]config.UserForm, error) {
	var users []config.UserForm
	var relations []config.Relation
	//从关注表中获取用户的关注列表
	res := Db.Where("user_id = ?", user_id).Find(&relations)
	if res.Error != nil {
		return []config.UserForm{}, res.Error
	}
	//从数据库中查找用户信息
	for _, relation := range relations {
		user := config.UserForm{}
		res := Db.Where("user_id=?", relation.ToUserId).Find(&user)
		if res.Error != nil {
			return []config.UserForm{}, res.Error
		}
		users = append(users, user)
	}
	if res.Error != nil {
		return users, res.Error
	}
	return users, nil
}

func IsFollow(user_id int64, to_user_id int64) (error, bool) {

	res := Db.Where("user_id=? and to_user_id=?", user_id, to_user_id).Find(&config.Relation{})
	if res.Error != nil {
		return res.Error, false
	}
	if res.RowsAffected == 0 {
		return nil, false
	}
	return nil, true
}
func GetFollowCount(user_id int64) (int64, error) {
	var relations []config.Relation
	res := Db.Where("user_id=?", user_id).Find(&relations)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func GetFollowerCount(user_id int64) (int64, error) {
	var relations []config.Relation
	res := Db.Where("to_user_id=?", user_id).Find(&relations)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func GetFollowerList(user_id int64) ([]config.UserForm, error) {
	var users []config.UserForm
	var relations []config.Relation
	//从关注表中获取用户的粉丝列表
	res := Db.Where("to_user_id = ?", user_id).Find(&relations)
	if res.Error != nil {
		return []config.UserForm{}, res.Error
	}
	//从数据库中查找用户信息
	for _, relation := range relations {
		user := config.UserForm{}
		res := Db.Where("user_id=?", relation.UserId).Find(&user)
		if res.Error != nil {
			return []config.UserForm{}, res.Error
		}
		users = append(users, user)
	}
	if res.Error != nil {
		return users, res.Error
	}
	return users, nil
}
