package controller

import (
	"errors"
	"fmt"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/mysql"
	"github.com/RaymondCode/simple-demo/pkg/jwt"
	"github.com/RaymondCode/simple-demo/pkg/snowflake"
	"github.com/gin-gonic/gin"
	"net/http"
)

var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	//获取和检验参数
	username := c.Query("username")
	password := c.Query("password")
	fmt.Println(username, password)
	if username == "" || password == "" {
		c.JSON(404, Response{StatusCode: 422, StatusMsg: "username or password cannot be empty!"})
		return
	}
	//是否存在用户表如果不存在则创建
	err := mysql.Db.AutoMigrate(&config.UserForm{})
	if err != nil {
		c.JSON(200, Response{StatusCode: 200,
			StatusMsg: "An error occurred on the server. Contact your administrator"})
		return
	}

	//检验用户是否存在
	err = mysql.CheckUserExist(username)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 200,
			StatusMsg: err.Error()})
		return
	}

	//生成UID
	userId, err := snowflake.GetID()
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 200,
			StatusMsg: errors.New("failed to create a user").Error()})
		return
	}
	// 构造一个UserForm实例
	user := config.UserForm{
		UserId:   int64(userId),
		Username: username,
		Password: password,
	}

	//保存进数据库
	Id, err := mysql.InsertUser(user)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 200,
			StatusMsg: errors.New("failed to create a user").Error()})
		return
	}
	// 生成JWT
	token, err := jwt.GenToken(Id, username)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 200,
			StatusMsg: errors.New("failed to create a token").Error()})
		return
	}
	//注册成功返回token和用户id
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0, StatusMsg: "Registered successfully"},
		UserId:   Id,
		Token:    token,
	})
}

func Login(c *gin.Context) {
	//获取和检验参数
	username := c.Query("username")
	password := c.Query("password")
	fmt.Println(username, password)
	if username == "" || password == "" {
		c.JSON(404, Response{StatusCode: 422, StatusMsg: "username or password cannot be empty!"})
		return
	}
	//数据库登录检验
	err, Id := mysql.Login(username, password)
	if err != nil {
		c.JSON(404, Response{StatusCode: 422, StatusMsg: err.Error()})
		return
	}
	// 生成JWT
	token, err := jwt.GenToken(Id, username)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 200,
			StatusMsg: errors.New("failed to create a token").Error()})
		return
	}
	//登录成功返回token和用户id
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0, StatusMsg: "Login successfully"},
		UserId:   Id,
		Token:    token,
	})
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	id := c.Query("user_id")
	if token == "" || id == "" {
		c.JSON(404, Response{StatusCode: 404, StatusMsg: "token or id cannot be empty!"})
		return
	}
	//解析token
	claims, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 200,
			StatusMsg: err.Error()})
		return
	}
	userID := claims.UserID
	username := claims.Username
	followCount, err := mysql.GetFollowCount(userID)
	if err != nil {
		c.JSON(501, Response{
			StatusCode: 501,
			StatusMsg:  err.Error(),
		})
		return
	}
	FollowerCount, err := mysql.GetFollowerCount(userID)
	if err != nil {
		c.JSON(501, Response{
			StatusCode: 501,
			StatusMsg:  err.Error(),
		})
		return
	}
	//创建实例
	user := User{
		Id:            userID,
		Name:          username,
		FollowCount:   followCount,
		FollowerCount: FollowerCount,
		IsFollow:      true,
	}

	//返回响应
	c.JSON(http.StatusOK, UserResponse{
		Response{StatusCode: 0, StatusMsg: ""},
		user,
	})
}

func GetUserInfoById(userId int64, ToUserId int64) (error, User) {
	err, user := mysql.GetUserInfoById(ToUserId)
	if err != nil {
		return err, User{}
	}
	followCount, err := mysql.GetFollowCount(ToUserId)
	if err != nil {
		return err, User{}
	}
	FollowerCount, err := mysql.GetFollowerCount(ToUserId)
	if err != nil {
		return err, User{}
	}
	err, IsFollow := mysql.IsFollow(userId, ToUserId)
	if err != nil {
		return err, User{}
	}
	User := User{
		Id:            user.UserId,
		Name:          user.Username,
		FollowCount:   followCount,
		FollowerCount: FollowerCount,
		IsFollow:      IsFollow,
	}
	return nil, User
}
