package USer

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserResponse 用户信息回复结构体
type UserResponse struct {
	CommonResponse
	User *USer `json:"user"` // 用户信息
}

// UserInfoHandler 用户信息处理函数，用于处理http请求
func UserInfoHandler(c *gin.Context) {
	//从请求中获取用户id
	id, exists := c.Get("user_id")
	if !exists {
		UserInfoErr(c, 3, "获取用户id失败")
		return
	}
	//查询用户信息
	userDao := NewUserInfoDao()
	userinfo, err := userDao.GetUserByuserID(id)
	if err != nil {
		UserInfoErr(c, 4, "获取用户信息失败")
		return
	}
	UserInfoOK(c, userinfo)
}

// UserInfoOK 返回正确信息
func UserInfoOK(c *gin.Context, login *USer) {
	c.JSON(http.StatusOK, UserResponse{
		CommonResponse: CommonResponse{
			StatusCode: 0,
		},
		User: &USer{
			ID:            login.ID,
			Name:          login.Name,
			FollowerCount: login.FollowerCount,
			FollowCount:   login.FollowCount,
			IsFollow:      login.IsFollow,
		},
	})
}

// UserInfoErr 返回错误信息
func UserInfoErr(c *gin.Context, code int64, errmessage string) {
	c.JSON(http.StatusOK, UserResponse{
		CommonResponse: CommonResponse{
			StatusCode: code,
			StatusMsg:  errmessage,
		},
	})
}
