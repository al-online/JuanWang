package config

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
