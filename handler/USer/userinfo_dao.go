package USer

import (
	"douyin.core/dal"
)

// USer 用户信息表
type USer struct {
	FollowCount   int64  `gorm:"follow_count" json:"follow_count,omitempty"`     // 关注总数
	FollowerCount int64  `gorm:"follower_count" json:"follower_count,omitempty"` // 粉丝总数
	ID            int64  `gorm:"id" json:"id,omitempty"`                         // 用户id
	IsFollow      bool   `gorm:"is_follow" json:"is_follow,omitempty"`           // true-已关注，false-未关注
	Name          string `gorm:"name" json:"name,omitempty"`                     // 用户名称
}

// UserInfoDao 用户信息数据操作结构体
type UserInfoDao struct {
}

// NewUserInfoDao 用户信息数据操作结构体构造函数
func NewUserInfoDao() *UserInfoDao {
	return &UserInfoDao{}
}

// GetUserByUserName 通过用户名查找用户
func (u *UserInfoDao) GetUserByUserName(username string) (*USer, error) {
	var User USer
	err := dal.DB.Where("username=?", username).First(&User).Error
	if err != nil {
		return nil, err
	}
	return &User, nil
}

// InsertToUserInfoTable 将用户信息持久化到数据库
func (u *UserInfoDao) InsertToUserInfoTable(userid int64, username string) error {
	user := USer{
		FollowCount:   0,
		FollowerCount: 0,
		ID:            userid,
		IsFollow:      false,
		Name:          username,
	}
	return dal.DB.Create(&user).Error
}

// GetUserByuserID 通过用户ID查找用户
func (u *UserInfoDao) GetUserByuserID(userid interface{}) (*USer, error) {
	var user USer
	err := dal.DB.Where("id=?", userid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserNameByUserID 通过用户id查询用户名，场景：用户上传视频的时候用于生成视频的文件名
func (u *UserInfoDao) GetUserNameByUserID(userid int64) (string, error) {
	var username string
	err := dal.DB.Select("username").Where("user_id=?", userid).First(&username).Error
	if err != nil {
		return "", err
	}
	return username, nil
}
