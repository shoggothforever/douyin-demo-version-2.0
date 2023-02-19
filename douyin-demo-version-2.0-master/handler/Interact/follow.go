package Interact

import (
	"douyin.core/dal"
	"douyin.core/handler/USer"
	"fmt"
	"strconv"
)

const (
	FOLLOWED   int = 1
	UNFOLLOWED int = 2
)

var conn = dal.Redisclient
var ctx = dal.Ctx
var db = dal.DB

type Relation struct {
	FollowCount   int64 `gorm:"follow_count" json:"follow_count"`     // 关注总数
	FollowerCount int64 `gorm:"follower_count" json:"follower_count"` // 粉丝总数
	ID            int64 `gorm:"id" json:"id"`                         // 用户id
}

func NewRelation(user USer.USer) *Relation {
	return &Relation{user.FollowCount, user.FollowerCount, user.ID}
}
func FollowCountkey(id int64) string {
	name := fmt.Sprintf("CommonUser%d'sFollowCount", id)
	return name
}
func FollowerCountkey(id int64) string {
	name := fmt.Sprintf("CommonUser%d'sFollowerCount", id)
	return name
}
func FollowSetkey(id string) string {
	name := fmt.Sprintf("CommonUser%d'sFollowCount", id)
	return name
}
func FollowerSetkey(id string) string {
	name := fmt.Sprintf("CommonUser%d'sFollowerCount", id)
	return name
}

//前者关注后者，前者关注数量+1，后者粉丝数量+1,
//还要更新关注列表，粉丝列表
func Follow(userid, id int64) error {
	var key = string(rune(userid))
	var val = string(rune(id))
	pipe := conn.Pipeline()
	if exi, err := conn.SIsMember(ctx, FollowSetkey(key), val).Result(); exi != true {
		if err != nil {
			return err
		}
		//关注操作，前者的关注集合添加后者，后者的粉丝集合添加前者
		_, adderr := pipe.SAdd(ctx, FollowSetkey(key), val).Result()
		_, adderr = pipe.SAdd(ctx, FollowerSetkey(val), key).Result()
		if adderr != nil {
			return adderr
		}
	} else {
		return nil
	}
	//前者的关注数以及后者的粉丝数加1
	if follow := pipe.Incr(ctx, FollowCountkey(userid)); follow.Err() != nil {
		return follow.Err()
	} else {
		if followed := pipe.Incr(ctx, FollowerCountkey(id)); followed.Err() != nil {
			return follow.Err()
		}
	}
	_, err := pipe.Exec(ctx)
	if err != nil { // 报错后进行一次额外尝试
		_, err = pipe.Exec(ctx)
		if err != nil {
			return nil
		}
	}
	return nil
}

//前者取关后者，前者关注数量-1，后者粉丝数量-1
func UnFollow(userid, id int64) error {
	var key = string(rune(userid))
	var val = string(rune(id))
	pipe := conn.Pipeline()
	if exi, err := conn.SIsMember(ctx, FollowSetkey(key), val).Result(); exi == true {
		if err != nil {
			return err
		}
		_, adderr := pipe.SRem(ctx, FollowSetkey(key), val).Result() //val==1，关注操作，前者的关注集合添加后者，后者的粉丝集合添加前者
		_, adderr = pipe.SRem(ctx, FollowerSetkey(val), key).Result()
		if adderr != nil {
			return adderr
		}
	} else {
		return nil
	}
	//前者的关注数以及后者的粉丝数减1
	if follow := pipe.Decr(ctx, FollowCountkey(userid)); follow.Err() != nil {
		return follow.Err()
	} else {
		if followed := pipe.Decr(ctx, FollowerCountkey(id)); followed.Err() != nil {
			return follow.Err()
		}
	}
	_, err := pipe.Exec(ctx)
	if err != nil { // 报错后进行一次额外尝试
		_, err = pipe.Exec(ctx)
		if err != nil {
			return nil
		}
	}
	return nil
}
func GetFollowList(id int64) ([]USer.USer, error) {
	var user []USer.USer
	vals, getallerr := conn.SMembers(ctx, FollowSetkey(strconv.FormatInt(id, 10))).Result()
	if getallerr != nil {
		return user, getallerr
	}
	err := db.Where("id in (?)", vals).Find(&user).Error
	return user, err
}
func GetFollowerList(id int64) ([]USer.USer, error) {
	var user []USer.USer
	vals, getallerr := conn.SMembers(ctx, FollowerSetkey(strconv.FormatInt(id, 10))).Result()
	if getallerr != nil {
		return user, getallerr
	}
	err := db.Where("id in (?)", vals).Find(&user).Error
	return user, err
}

//查询双向关注的好友
func GetFriendsList(id int64) ([]USer.USer, error) {
	var user []USer.USer
	var friends []string
	vals, getallerr := conn.SMembers(ctx, FollowerSetkey(strconv.FormatInt(id, 10))).Result()
	if getallerr != nil {
		return user, getallerr
	}
	for _, val := range vals {
		if exi, err := conn.SIsMember(ctx, FollowerSetkey(val), id).Result(); err != nil {
			return user, err
		} else {
			if exi {
				friends = append(friends, val)
			}
		}
	}

	err := db.Where("id in (?)", friends).Find(&user).Error
	return user, err
}
