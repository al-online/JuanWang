package controller

import (
	"github.com/RaymondCode/simple-demo/mysql"
	"github.com/RaymondCode/simple-demo/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {

	//获取请求参数并校验参数
	token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")
	if token == "" || to_user_id == "" || action_type == "" {
		c.JSON(404, Response{StatusCode: 422, StatusMsg: "param cannot be empty!"})
		return
	}
	//处理ToUserID
	ToUserID, err := strconv.ParseInt(to_user_id, 10, 64)
	if err != nil {
		c.JSON(501, Response{
			StatusCode: 501,
			StatusMsg:  err.Error(),
		})
		return
	}
	//检验token合法性
	myClaims, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(404, Response{
			StatusCode: 422,
			StatusMsg:  err.Error(),
		})
		return
	}
	userId := myClaims.UserID
	//检验该用户是否存在
	err = mysql.CheckUserExistById(userId)
	if err != nil {
		c.JSON(404, Response{
			StatusCode: 422,
			StatusMsg:  err.Error(),
		})
		return
	}
	//关注操作
	if action_type == "1" {
		//向数据库插入数据
		err := mysql.InsertRelation(userId, ToUserID)
		if err != nil {
			c.JSON(502, Response{
				StatusCode: 502,
				StatusMsg:  err.Error(),
			})
			return
		}
		c.JSON(200, Response{
			StatusCode: 0,
			StatusMsg:  "follow successfully!",
		})
		return
	} else if action_type == "2" { //取消关注
		err := mysql.DeleteRelation(userId, ToUserID)
		if err != nil {
			c.JSON(502, Response{
				StatusCode: 502,
				StatusMsg:  err.Error(),
			})
			return
		}
		c.JSON(200, Response{
			StatusCode: 0,
			StatusMsg:  "cancel follow successfully!",
		})
		return
	} else { //非法操作
		c.JSON(404, Response{
			StatusCode: 422,
			StatusMsg:  "illegal action type",
		})
		return
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	//检验参数合法性
	token := c.Query("token")
	userId := c.Query("user_id")
	if token == "" || userId == "" {
		c.JSON(404, Response{StatusCode: 422, StatusMsg: "param cannot be empty!"})
		return
	}
	//检验token合法性
	myClaims, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(404, Response{
			StatusCode: 422,
			StatusMsg:  err.Error(),
		})
		return
	}
	user_id := myClaims.UserID
	userlist, err := mysql.GetFollowList(user_id)
	if err != nil {
		c.JSON(501, Response{
			StatusCode: 501,
			StatusMsg:  err.Error(),
		})
		return
	}
	var Users []User
	//将用户列表信息完善
	for _, item := range userlist {
		//该用户的关注数
		followCount, err := mysql.GetFollowCount(item.UserId)
		if err != nil {
			c.JSON(501, Response{
				StatusCode: 501,
				StatusMsg:  err.Error(),
			})
			return
		}
		//该用户的粉丝数
		followerCount, err := mysql.GetFollowerCount(item.UserId)
		if err != nil {
			if err != nil {
				c.JSON(501, Response{
					StatusCode: 501,
					StatusMsg:  err.Error(),
				})
				return
			}
		}
		//该用户是否已关注
		err, isFollow := mysql.IsFollow(user_id, item.UserId)
		if err != nil {
			if err != nil {
				c.JSON(501, Response{
					StatusCode: 501,
					StatusMsg:  err.Error(),
				})
				return
			}
		}
		user := User{
			Id:            item.UserId,
			Name:          item.Username,
			FollowCount:   followCount,
			FollowerCount: followerCount,
			IsFollow:      isFollow,
		}
		Users = append(Users, user)
	}
	c.JSON(200, UserListResponse{
		UserList: Users,
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "get followList successfully!",
		},
	})

}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	//检验参数合法性
	token := c.Query("token")
	userId := c.Query("user_id")
	if token == "" || userId == "" {
		c.JSON(404, Response{StatusCode: 422, StatusMsg: "param cannot be empty!"})
		return
	}
	//检验token合法性
	myClaims, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(404, Response{
			StatusCode: 422,
			StatusMsg:  err.Error(),
		})
		return
	}
	user_id := myClaims.UserID
	userlist, err := mysql.GetFollowerList(user_id)
	if err != nil {
		c.JSON(501, Response{
			StatusCode: 501,
			StatusMsg:  err.Error(),
		})
		return
	}
	var Users []User
	//将用户列表信息完善
	for _, item := range userlist {
		//该用户的关注数
		followCount, err := mysql.GetFollowCount(item.UserId)
		if err != nil {
			c.JSON(501, Response{
				StatusCode: 501,
				StatusMsg:  err.Error(),
			})
			return
		}
		//该用户的粉丝数
		followerCount, err := mysql.GetFollowerCount(item.UserId)
		if err != nil {
			if err != nil {
				c.JSON(501, Response{
					StatusCode: 501,
					StatusMsg:  err.Error(),
				})
				return
			}
		}
		//该用户是否已关注
		err, isFollow := mysql.IsFollow(user_id, item.UserId)
		if err != nil {
			if err != nil {
				c.JSON(501, Response{
					StatusCode: 501,
					StatusMsg:  err.Error(),
				})
				return
			}
		}
		user := User{
			Id:            item.UserId,
			Name:          item.Username,
			FollowCount:   followCount,
			FollowerCount: followerCount,
			IsFollow:      isFollow,
		}
		Users = append(Users, user)
	}
	c.JSON(200, UserListResponse{
		UserList: Users,
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "get followerList successfully!",
		},
	})
}
