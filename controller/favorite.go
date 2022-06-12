package controller

import (
	"github.com/RaymondCode/simple-demo/mysql"
	"github.com/RaymondCode/simple-demo/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	//获取请求参数并校验参数
	token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")
	if token == "" || video_id == "" || action_type == "" {
		c.JSON(404, Response{StatusCode: 422, StatusMsg: "param cannot be empty!"})
		return
	}
	//处理video_id
	VideoId, err := strconv.ParseInt(video_id, 10, 64)
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
	//点赞操作
	if action_type == "1" {
		//向数据库插入数据
		err := mysql.InsertFavorite(userId, VideoId)
		if err != nil {
			c.JSON(502, Response{
				StatusCode: 502,
				StatusMsg:  err.Error(),
			})
			return
		}
		c.JSON(200, Response{
			StatusCode: 0,
			StatusMsg:  "Favorite successfully!",
		})
		return
	} else if action_type == "2" { //取消点赞
		err := mysql.DeleteFavorite(userId, VideoId)
		if err != nil {
			c.JSON(502, Response{
				StatusCode: 502,
				StatusMsg:  err.Error(),
			})
			return
		}
		c.JSON(200, Response{
			StatusCode: 0,
			StatusMsg:  "cancel favorite successfully!",
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

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	//获取token
	token := c.Query("token")
	//检验token
	myClaims, err := jwt.ParseToken(token)
	//如果鉴权失败
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 401,
			StatusMsg:  err.Error(),
		})
		return
	}
	//获取用户id
	UserId := myClaims.UserID
	//获取数据库中的点赞视频列表
	videos, err := mysql.GetFavoriteList(UserId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 401,
			StatusMsg:  err.Error(),
		})
		return
	}
	if err != nil {
		c.JSON(200, Response{
			StatusCode: 501,
			StatusMsg:  err.Error(),
		})
	}
	VideosList := []Video{}
	for _, item := range videos {
		//获取该视频用户信息
		err, User := GetUserInfoById(item.UserId, item.UserId)
		if err != nil {
			c.JSON(200, Response{
				StatusCode: 501,
				StatusMsg:  err.Error(),
			})
		}
		FavoriteCount, err := mysql.GetFavoriteCount(item.Id)
		if err != nil {
			c.JSON(200, Response{
				StatusCode: 501,
				StatusMsg:  err.Error(),
			})
		}
		commentCount, err := mysql.GetCommentCount(item.Id)
		if err != nil {
			c.JSON(200, Response{
				StatusCode: 501,
				StatusMsg:  err.Error(),
			})
		}
		video := Video{
			Id:            item.Id,
			Author:        User,
			PlayUrl:       item.VideoUrl,
			CoverUrl:      item.VideoCover,
			FavoriteCount: FavoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    true,
		}
		VideosList = append(VideosList, video)
	}

	c.JSON(200, VideoListResponse{
		VideoList: VideosList,
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "Successfully obtained the videosList",
		},
	})
}
