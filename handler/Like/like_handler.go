package Like

import (
	"douyin.core/middleware"
	"github.com/gin-gonic/gin"
	"strconv"
)

func LikeHandler(c *gin.Context) {
	token, ok := c.GetQuery("token")
	if !ok {
		LikeResponse(c, 1, "未能成功获取token，请重试")
		return
	}
	userclaim, err := middleware.JwtParseUser(token)
	if err != nil {
		LikeResponse(c, 0, "token已过期，请重新登录")
		return
	}
	userid := userclaim.Userid
	//获取视频id
	videoid, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if !ok {
		LikeResponse(c, 1, "未能成功获取视频id，请重试")
		return
	}
	actionType, ok := c.GetQuery("action_type")
	if !ok {
		LikeResponse(c, 1, "未能成功获取操作类型，请重试")
		return
	}
	dao := LikeDAO{}
	switch actionType {
	case "1":
		//点赞
		err := dao.AddLike(userid, videoid)
		if err != nil {
			LikeResponse(c, 1, "点赞失败")
			return
		}
		LikeResponse(c, 0, "点赞成功")
	case "2":
		//取消点赞
		err := dao.CancelLike(userid, videoid)
		if err != nil {
			LikeResponse(c, 1, "取消点赞失败")
			return
		}
		LikeResponse(c, 0, "取消点赞成功")
	default:
		LikeResponse(c, 1, "未知操作类型")
		return

	}
}

func LikeResponse(c *gin.Context, statuscode int64, statusmsg string) {
	c.JSON(200, gin.H{
		"status_code": statuscode,
		"status_msg":  statusmsg,
	})
}
