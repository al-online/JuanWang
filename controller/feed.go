package controller

import (
	"github.com/RaymondCode/simple-demo/mysql"
	"github.com/RaymondCode/simple-demo/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	var err error
	token := c.Query("token")
	//处理时间
	lastTime := c.Query("last_time")
	var last_time int64
	if lastTime == "" {
		last_time = time.Now().Unix()
	} else {
		last_time, err = strconv.ParseInt(lastTime, 10, 64)
		if err != nil {
			c.JSON(200, Response{
				StatusCode: 501,
				StatusMsg:  err.Error(),
			})
		}
	}
	//用户未登录
	if token == "" {
		videos, err := mysql.VideosList(last_time)
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
			video := Video{
				Id:            item.Id,
				Author:        User,
				PlayUrl:       item.VideoUrl,
				CoverUrl:      item.VideoCover,
				FavoriteCount: 90909,
				CommentCount:  20919,
				IsFavorite:    false,
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

	} else { //用户已经登录
		//鉴权,获取用户id
		myClaims, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(200, Response{
				StatusCode: 501,
				StatusMsg:  err.Error(),
			})
			return
		}
		UserId := myClaims.UserID
		//数据库中获取视频列表
		videos, err := mysql.VideosList(last_time)
		if err != nil {
			c.JSON(200, Response{
				StatusCode: 501,
				StatusMsg:  err.Error(),
			})
		}
		//获取next_time
		var NextTime int64
		if len(videos) != 0 {
			NextTime = videos[len(videos)-1].CreatedAt.Unix()
		} else {
			NextTime = 0
		}
		VideosList := []Video{}
		for _, item := range videos {
			//获取该视频用户信息
			err, User := GetUserInfoById(UserId, item.UserId)
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
				FavoriteCount: 90909,
				CommentCount:  20919,
				IsFavorite:    false,
			}
			VideosList = append(VideosList, video)
		}

		c.JSON(200, FeedResponse{
			VideoList: VideosList,
			NextTime:  NextTime,
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "Successfully obtained the videosList",
			},
		})
	}
}
