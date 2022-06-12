package mysql

import "github.com/RaymondCode/simple-demo/config"

func InsertFavorite(userId int64, VideoId int64) error {
	var err error
	//是否存在视频点赞表如果不存在则创建
	err = Db.AutoMigrate(&config.Favorite{})
	if err != nil {
		return err
	}
	favorite := config.Favorite{
		UserId:  userId,
		VideoId: VideoId,
	}
	//添加记录
	res := Db.Create(&favorite)
	if res.Error != nil {
		return res.Error
	}
	return err
}

func DeleteFavorite(userId int64, VideoId int64) error {
	var err error
	favorite := config.Favorite{
		UserId:  userId,
		VideoId: VideoId,
	}
	//删除记录
	res := Db.Where("user_id=? and video_id=?", userId, VideoId).Delete(&favorite)
	if res.Error != nil {
		return res.Error
	}
	return err
}

// GetFavoriteList 获取点赞列表
func GetFavoriteList(user_id int64) ([]config.Video, error) {
	var videos []config.Video
	var favorites []config.Favorite
	//从点赞表中获取用户的点赞列表
	res := Db.Where("user_id = ?", user_id).Find(&favorites)
	if res.Error != nil {
		return []config.Video{}, res.Error
	}
	//从数据库中的视频表查找视频信息
	for _, favorite := range favorites {
		video := config.Video{}
		res := Db.Where("id=?", favorite.VideoId).Find(&video)
		if res.Error != nil {
			return []config.Video{}, res.Error
		}
		videos = append(videos, video)
	}
	if res.Error != nil {
		return videos, res.Error
	}
	return videos, nil
}

func GetFavoriteCount(video_id int64) (int64, error) {
	var favorites []config.Favorite
	res := Db.Where("video_id=?", video_id).Find(&favorites)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func IsFavorite(user_id int64, video_id int64) (error, bool) {

	res := Db.Where("user_id=? and video_id=?", user_id, video_id).Find(&config.Favorite{})
	if res.Error != nil {
		return res.Error, false
	}
	if res.RowsAffected == 0 {
		return nil, false
	}
	return nil, true
}
