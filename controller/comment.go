package controller

import (
	"github.com/RaymondCode/simple-demo/mysql"
	"github.com/RaymondCode/simple-demo/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
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
	videoId := c.Query("video_id")
	if videoId == "" {
		c.JSON(200, Response{
			StatusCode: 401,
			StatusMsg:  "videoId does not empty",
		})
		return
	}
	VideoId, err := strconv.ParseInt(videoId, 10, 64)
	if err != nil {
		c.JSON(200, Response{
			StatusCode: 501,
			StatusMsg:  err.Error(),
		})
		return
	}
	actionType := c.Query("action_type")
	//如果是评论操作
	if actionType == "1" {
		text := c.Query("comment_text")
		if text == "" {
			c.JSON(200, Response{
				StatusCode: 401,
				StatusMsg:  "text does not empty",
			})
			return
		}
		//插入评论记录到数据库
		err, CommentId, Content, Time := mysql.InsertComment(UserId, VideoId, text)
		if err != nil {
			c.JSON(200, Response{
				StatusCode: 501,
				StatusMsg:  err.Error(),
			})
			return
		}
		//格式化时间字符串
		t := Time.Format("01-02")
		err, User := GetUserInfoById(UserId, UserId)
		if err != nil {
			c.JSON(200, Response{
				StatusCode: 501,
				StatusMsg:  err.Error(),
			})
			return
		}
		c.JSON(200, CommentActionResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "comment successfully!",
			},
			Comment: Comment{
				Id:         CommentId,
				User:       User,
				Content:    Content,
				CreateDate: t,
			},
		})
	} else if actionType == "2" {
		//删除评论操作
		comment_id := c.Query("comment_id")
		ComentId, err := strconv.ParseInt(comment_id, 10, 64)
		if err != nil {
			c.JSON(200, Response{
				StatusCode: 501,
				StatusMsg:  err.Error(),
			})
			return
		}
		err = mysql.DeleteComment(UserId, VideoId, ComentId)
		if err != nil {
			c.JSON(200, Response{
				StatusCode: 501,
				StatusMsg:  err.Error(),
			})
			return
		}
		c.JSON(200, CommentActionResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "delete comment successfully!",
			},
			Comment: Comment{},
		})
	} else {
		//其他非法操作
		c.JSON(404, Response{
			StatusCode: 422,
			StatusMsg:  "illegal action type",
		})
		return
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
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
	videoId := c.Query("video_id")
	if videoId == "" {
		c.JSON(200, Response{
			StatusCode: 401,
			StatusMsg:  "videoId does not empty",
		})
		return
	}
	VideoId, err := strconv.ParseInt(videoId, 10, 64)
	//从数据库中获取此视频的评论
	comments, err := mysql.GetCommentList(VideoId)
	if err != nil {
		c.JSON(200, Response{
			StatusCode: 501,
			StatusMsg:  err.Error(),
		})
		return
	}
	var CommentList []Comment
	for _, item := range comments {
		Time := item.CreatedAt
		//格式化时间字符串
		t := Time.Format("01-02")
		err, User := GetUserInfoById(UserId, item.UserId)
		if err != nil {
			c.JSON(200, Response{
				StatusCode: 501,
				StatusMsg:  err.Error(),
			})
			return
		}
		comment := Comment{
			Id:         item.Id,
			User:       User,
			Content:    item.Content,
			CreateDate: t,
		}
		CommentList = append(CommentList, comment)
	}

	c.JSON(200, CommentListResponse{
		CommentList: CommentList,
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "video comment list obtained successfully!",
		},
	})
}
