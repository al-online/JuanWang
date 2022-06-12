package mysql

import (
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/pkg/snowflake"
	"time"
)

func InsertComment(UserId int64, VideoId int64, Content string) (error, int64, string, time.Time) {
	var err error
	//是否存在评论表,如果不存在则创建
	err = Db.AutoMigrate(&config.Comment{})
	if err != nil {
		return err, 0, "", time.Time{}
	}
	//雪花算法生成评论id
	id, err := snowflake.GetID()
	if err != nil {
		return err, 0, "", time.Time{}
	}
	comment := config.Comment{
		Id:      int64(id),
		UserId:  UserId,
		VideoId: VideoId,
		Content: Content,
	}
	//添加记录
	res := Db.Create(&comment)
	if res.Error != nil {
		return res.Error, 0, "", time.Time{}
	}
	return err, comment.Id, comment.Content, comment.CreatedAt
}

func DeleteComment(userId int64, VideoId int64, CommentId int64) error {
	var err error
	//删除记录
	res := Db.Where("user_id=? and video_id=? and id=?", userId, VideoId, CommentId).Delete(&config.Comment{})
	if res.Error != nil {
		return res.Error
	}
	return err
}

func GetCommentList(videoId int64) ([]config.Comment, error) {
	commentList := []config.Comment{}
	res := Db.Where("video_id=?", videoId).Find(&commentList)
	if res.Error != nil {
		return []config.Comment{}, res.Error
	}
	return commentList, nil
}

func GetCommentCount(videoId int64) (int64, error) {
	var comments []config.Comment
	res := Db.Where("video_id=?", videoId).Find(&comments)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}
