package mysql

import (
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/pkg/snowflake"
)

func PublishVideo(fileName string, userId int64, coverName string, Title string) error {
	//是否存在视频表,如果不存在则创建
	err := Db.AutoMigrate(&config.Video{})
	if err != nil {
		return err
	}
	//生成视频id
	id, err := snowflake.GetID()
	if err != nil {
		return err
	}
	//拼接url
	videoUrl := "http://172.18.154.64:8080/static/" + fileName
	coverUrl := "http://172.18.154.64:8080/static/" + coverName
	video := config.Video{
		Id:         int64(id),
		UserId:     userId,
		Title:      Title,
		VideoUrl:   videoUrl,
		VideoCover: coverUrl,
	}
	//插入数据库
	res := Db.Create(&video)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func PublishList(userId int64) (error, []config.Video) {
	videos := []config.Video{}
	res := Db.Where("user_id=?", userId).Find(&videos)
	if res.Error != nil {
		return res.Error, videos
	}
	return nil, videos
}
