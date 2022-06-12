package mysql

import (
	"github.com/RaymondCode/simple-demo/config"
	"time"
)

func VideosList(lastTime int64) ([]config.Video, error) {
	videos := []config.Video{}
	t := time.Unix(lastTime, 0)
	Time := t.Format("2006/01/02 15:04:05")
	//获取前30条数据，并倒序
	res := Db.Where("created_at < ?", Time).Order("created_at DESC").Limit(30).Find(&videos)
	if res.Error != nil {
		return videos, res.Error
	}
	//fmt.Println(videos)
	return videos, nil
}
