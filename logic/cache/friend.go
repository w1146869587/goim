package cache

import (
	"goim/logic/db"
	"goim/logic/model"
	"goim/public/imctx"
	"goim/public/logger"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	FriendKey    = "friend:"
	FriendExpire = 2 * time.Hour
)

type friendCache struct{}

var FriendCache = new(friendCache)

func (c *friendCache) Key(appId, userId, friendId int64) string {
	return FriendKey + strconv.FormatInt(appId, 10) + ":" + strconv.FormatInt(userId, 10) + ":" + strconv.FormatInt(friendId, 10)
}

// Get 获取朋友信息
func (c *friendCache) Get(ctx *imctx.Context, appId, userId, friendId int64) (friend *model.Friend, err error) {
	err = get(c.Key(appId, userId, friendId), friend)
	if err != nil && err != redis.Nil {
		logger.Sugar.Error(err)
		return
	}
	if err == redis.Nil {
		return nil, nil
	}
	return
}

// Set 设置朋友信息
func (c *friendCache) Set(ctx *imctx.Context, friend *model.Friend) error {
	err := set(c.Key(friend.AppId, friend.UserId, friend.FriendId), friend, FriendExpire)
	if err != nil {
		logger.Sugar.Error(err)
	}
	return err
}

// Delele 删除朋友信息
func (c *friendCache) Delele(ctx *imctx.Context, appId, userId, friendId int64) error {
	err := db.RedisCli.Del(c.Key(appId, userId, friendId)).Err()
	if err != nil {
		logger.Sugar.Error(err)
	}
	return err
}