package USer

import (
	"douyin.core/middleware"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

const (
	MaxUsernameLength = 20
	MaxPasswordLength = 18
	MinPasswordLength = 5
)

// RegisterReponse 用户注册回复结构体
type RegisterReponse struct {
	CommonResponse
	Token  string `json:"token"`   // 用户鉴权token
	UserID int64  `json:"user_id"` // 用户id
}

// PostUserLogin 用户登录校验
type PostUserLogin struct {
	Username string
	Password string
	Userid   int64
	Token    string
}

// UserRegistHandler 用户注册处理函数
func UserRegistHandler(c *gin.Context) {
	//获取用户名和密码
	username := c.Query("username")
	get, exists := c.Get("password")
	if !exists {
		RegisterErr(c, "密码不能为空")
		return
	}
	password := get.(string)
	//创建新对象
	newPostUserLogin := NewPostUserLogin(username, password)
	//注册新用户
	err := newPostUserLogin.Register()
	if err != nil {
		//返回错误信息
		RegisterErr(c, err.Error())
		return
	}
	//返回正确信息
	RegisterOK(c, newPostUserLogin)
}

// Register 注册新用户
func (u *PostUserLogin) Register() error {
	//校验参数
	err := u.CheckPost()
	if err != nil {
		return err
	}

	//生成userid
	u.UserIdGenarate()

	//持久化到数据库
	err = u.PersistData()
	if err != nil {
		return err
	}

	//获取token
	err = u.SetToken()
	if err != nil {
		return err
	}

	return nil
}

// NewPostUserLogin 创建对象，用于注册
func NewPostUserLogin(username, password string) *PostUserLogin {
	return &PostUserLogin{Username: username, Password: password}
}

// PersistData 将用户数据持久化到数据库并检查是否出现用户名重复的现象
func (u *PostUserLogin) PersistData() error {
	//创建用户表数据操作对象
	userDAO := NewUserInfoDao()
	//创建用户注册表数据操作对象
	userRigestDao := NewUserRigisterDao()
	//检查用户名是否已经存在
	_, err := userDAO.GetUserByUserName(u.Username)
	is := errors.Is(err, gorm.ErrRecordNotFound)
	if is {
		//将数据持久化到用户表
		err = userDAO.InsertToUserInfoTable(u.Userid, u.Username)
		if err != nil {
			return err
		}
		//对用户密码进行AES-256加密，保障用户安全
		passwordstr := []byte(u.Password)
		aes, err := middleware.EnPwdCode(passwordstr)
		if err != nil {
			return err
		}
		u.Password = aes
		//将数据持久化到用户注册表
		err = userRigestDao.RegistUsertoDb(u.Userid, u.Username, u.Password)
		if err != nil {
			return err
		}
		return nil
	} else {
		err := errors.New("用户名已存在，请重试")
		return err
	}
}

// SetToken 获取token
func (u *PostUserLogin) SetToken() error {
	token, err := middleware.JwtGenerateToken(u.Userid, time.Hour)
	if err != nil {
		return err
	}
	u.Token = token
	return nil
}

// UserIdGenarate 用户id生成
func (u *PostUserLogin) UserIdGenarate() {
	worker, _ := middleware.NewWorker(1)
	id := worker.GetId()
	u.Userid = id
}

// CheckPost 检查用户登录时的用户名和密码是否合法
func (u *PostUserLogin) CheckPost() error {
	if u.Username == "" {
		return errors.New("用户名不能为空")
	}
	if len(u.Username) > MaxUsernameLength {
		return errors.New("用户名过长")
	}

	if len(u.Password) > MaxPasswordLength {
		return errors.New("密码过长")
	}
	if len(u.Password) < MinPasswordLength {
		return errors.New("密码不能少于5位")
	}
	return nil
}

// RegisterOK 返回正确信息
func RegisterOK(c *gin.Context, login *PostUserLogin) {
	c.JSON(http.StatusOK, RegisterReponse{
		CommonResponse: CommonResponse{
			StatusCode: 0,
		},
		UserID: login.Userid,
		Token:  login.Token,
	})
}

// RegisterErr 返回错误信息
func RegisterErr(c *gin.Context, errmessage string) {
	c.JSON(http.StatusOK, RegisterReponse{
		CommonResponse: CommonResponse{
			StatusCode: 1,
			StatusMsg:  errmessage,
		},
	})
}
