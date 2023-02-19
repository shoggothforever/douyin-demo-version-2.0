package USer

import (
	"douyin.core/dal"
	"douyin.core/middleware"
	"errors"
)

// UserLoginTable 用户登录表
type UserLoginTable struct {
	Id       int64  `gorm:"primary_key"`
	UserId   int64  `gorm:"notnull"`
	Username string `gorm:"notnull"`
	Password string `gorm:"notnull"`
}

// UserRigestDao 用户注册表数据操作结构体
type UserRigestDao struct {
}

// NewUserLoginTable 新建用户
func NewUserLoginTable(username, password string) *UserLoginTable {
	return &UserLoginTable{Username: username, Password: password}
}

// NewUserRigisterDao 用户注册表数据操作结构体构造函数
func NewUserRigisterDao() *UserRigestDao {
	return &UserRigestDao{}
}

// RegistUsertoDb 将用户信息持久化到数据库
func (u *UserRigestDao) RegistUsertoDb(userid int64, username, password string) error {
	user := UserLoginTable{
		UserId:   userid,
		Username: username,
		Password: password,
	}
	return dal.DB.Create(&user).Error
}

// QueryUserLogin 用户登录时检查用户的参数是否正确
func (u UserRigestDao) QueryUserLogin(username, password string, login *UserLoginTable) error {
	err := dal.DB.Where("username=?", username).First(&login).Error
	if err != nil {
		return err
	}
	//对查询得到的用户密码进行解密
	bytes, err := middleware.DePwdCode(login.Password)
	if err != nil {
		return err
	}
	//判断密码是否正确
	if string(bytes) != password {
		err = errors.New("密码错误")
		return err
	}
	return nil
}
