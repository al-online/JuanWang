package config

import "time"

// UserForm 用户表单信息
type UserForm struct {
	UserId   int64  `json:"user_id" gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// TableName 自定义表名
func (UserForm) TableName() string {
	return "t_users"
}

type Relation struct {
	UserId   int64    `json:"user_id" gorm:"primary_key"`
	ToUserId int64    `json:"to_user_id" gorm:"primary_key"`
	User     UserForm `gorm:"foreignkey:UserId;foreignkey:ToUserId"` //用户所属用户表外键
}

func (Relation) TableName() string {
	return "t_relation"
}

type Video struct {
	Id         int64     `json:"id" gorm:"primary_key"`
	UserId     int64     `json:"user_id"`
	VideoUrl   string    `json:"video_url"`
	Title      string    `json:"title"`
	VideoCover string    `json:"video_cover"`
	CreatedAt  time.Time `json:"created_at"`
	User       UserForm  `gorm:"foreignkey:UserId"`
}

func (Video) TableName() string {
	return "t_video"
}
