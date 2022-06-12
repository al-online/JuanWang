package controller

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/RaymondCode/simple-demo/mysql"
	"github.com/RaymondCode/simple-demo/pkg/jwt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	//获取token
	token := c.PostForm("token")
	//检验token
	myClaims, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 401,
			StatusMsg:  err.Error(),
		})
		return
	}
	//获取视频标题
	title := c.PostForm("title")
	if title == "" {
		c.JSON(200, Response{
			StatusCode: 401,
			StatusMsg:  "title is not empty!",
		})
		return
	}
	//获取视频数据
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 401,
			StatusMsg:  err.Error(),
		})
		return
	}
	userId := myClaims.UserID
	//保存视频到本地
	err, finalName, coverName := SaveVideo(data, c, userId)
	if err != nil {
		c.JSON(200, Response{
			StatusCode: 501,
			StatusMsg:  err.Error(),
		})
		return
	}
	//保存到数据库
	err = mysql.PublishVideo(finalName, userId, coverName, title)
	if err != nil {
		c.JSON(200, Response{
			StatusCode: 501,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
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
	//获取数据库中的视频信息
	err, videos := mysql.PublishList(UserId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 401,
			StatusMsg:  err.Error(),
		})
		return
	}
	//获取用户信息
	err, User := GetUserInfoById(UserId, UserId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 401,
			StatusMsg:  err.Error(),
		})
		return
	}
	VideosList := []Video{}
	for _, item := range videos {
		err, isFavorite := mysql.IsFavorite(UserId, item.Id)
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
			IsFavorite:    isFavorite,
		}
		VideosList = append(VideosList, video)
	}

	c.JSON(200, VideoListResponse{
		VideoList: VideosList,
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "Successfully obtained the publishList",
		},
	})
}

func SaveVideo(data *multipart.FileHeader, c *gin.Context, UserId int64) (error, string, string) {
	filename := filepath.Base(data.Filename)
	fileExt := strings.ToLower(path.Ext(filename))
	//检查文件后缀
	if fileExt != ".mp4" && fileExt != ".avi" && fileExt != ".mpg" {
		return errors.New("this file type is not supported"), "", ""
	}
	finalName := fmt.Sprintf("%d_%s", UserId, filename)
	saveFile := filepath.Join("./public/", finalName)
	//保存文件
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		return err, "", ""
	}
	//获取视频封面
	err, coverName := GetCover(finalName)
	if err != nil {
		return err, "", ""
	}
	return nil, finalName, coverName
}

func GetCover(filepath string) (error, string) {
	//获取视频中的第一帧
	err, coverFileName := ReadFrameAsJpeg(filepath, 1)
	return err, coverFileName
}

func ReadFrameAsJpeg(inFileName string, frameNum int) (error, string) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(filepath.Join("./public/", inFileName)).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		return err, ""
	}
	//解码
	img, err := imaging.Decode(buf)
	if err != nil {
		return err, ""
	}
	//保存截图coverName
	coverName := strings.Replace(inFileName, ".mp4", ".jpeg", -1)
	err = imaging.Save(img, filepath.Join("./public/", coverName))
	if err != nil {
		return err, ""
	}
	return nil, coverName
}
