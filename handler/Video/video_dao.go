package Video

import (
	"douyin.core/Model"
	"douyin.core/config"
	user "douyin.core/handler/User"
	"errors"
	"gorm.io/gorm"
	"time"
)

// Video
type Video struct {
	Author        *user.User `json:"author"`         // 视频作者信息
	UserID        int64      `json:"user_id"`        //用户id
	CommentCount  int64      `json:"comment_count"`  // 视频的评论总数
	CoverURL      string     `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64      `json:"favorite_count"` // 视频的点赞总数
	ID            int64      `json:"id"`             // 视频唯一标识
	IsFavorite    bool       `json:"is_favorite"`    // true-已点赞，false-未点赞
	PlayURL       string     `json:"play_url"`       // 视频播放地址
	Title         string     `json:"title"`          // 视频标题
	UserVideocode int64      `json:"videocode"`      //用户视频编号
}

// 视频表数据操作结构体
type VideoDao struct {
}

// 视频表数据操作结构体构造函数
func NewVideoDao() *VideoDao {
	return &VideoDao{}
}

// todo   查询用户视频关系表排除用户看过的视频，嵌套结构体查询问题
func (v *VideoDao) QueryVideoListLogin(userid int64, last_time time.Time) (*VideoList, error) {
	var videolist VideoList
	videolist.Videos = make([]*Video, 0, config.MaxVideoList)
	err := Model.DB.Where("publish_time<?", last_time).Order("publish_time desc").Limit(config.MaxVideoList).Find(&videolist.Videos).Error
	if err != nil {
		return nil, err
	}
	return &videolist, nil
}

// todo 嵌套结构体查询问题
func (v *VideoDao) QueryVideoListUnLogin(last_time time.Time) (*VideoList, error) {
	var videolist VideoList
	videolist.Videos = make([]*Video, 0, config.MaxVideoList)
	err := Model.DB.Where("publish_time<?", last_time).Order("publish_time desc").Limit(config.MaxVideoList).Find(&videolist.Videos).Error
	if err != nil {
		return nil, err
	}
	return &videolist, nil
}

// todo 未完成，等待粉丝和点赞系统完成
func (v *VideoDao) PersistNewVideo(title string, userid int64, user *user.UserInfoDao) error {
	userinfo, err := user.GetUserByuserID(userid)
	if err != nil {
		return err
	}
	video := &Video{
		Author:        userinfo,
		CommentCount:  0,
		CoverURL:      "",
		FavoriteCount: 0,
		ID:            0,
		IsFavorite:    false,
		PlayURL:       "",
		Title:         title,
		UserVideocode: 0,
	}
	return Model.DB.Create(video).Error
}

// 获取用户视频序号吗，场景：用于用户将视频上传的时候生成视频文件名
func (v *VideoDao) GetUserVideoCode(userid int64) (int64, error) {
	var videocode int64
	err := Model.DB.Select("videocode").Where("userid=?", userid).First(&videocode).Error
	is := errors.Is(err, gorm.ErrRecordNotFound)
	if is {
		return 0, err
	}
	if err != nil {
		return -1, err
	}
	return videocode + 1, nil
}

// todo 嵌套结构体查询问题
func (v *VideoDao) QueryUserPublishList(userid int64) (*VideoList, error) {
	var videos []*Video
	err := Model.DB.Where("userid=?", userid).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return &VideoList{Videos: videos}, nil
}

func QueryVideoById(vid int64) (Video, error) {
	var video Video
	err := Model.DB.Where("videoid = ?", vid).First(&video).Error
	return video, err
}
