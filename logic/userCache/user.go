package userCache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"gomall/consts"
	"gomall/global"
	"gomall/models/users"
	"time"
)

// LoadUserByUsername 查找用户信息
func LoadUserByUsername(username string) (user *users.User, err error) {
	//记得先初始化，不然json.Unmarshal会报错：将string解析到一个空指针上
	user = &users.User{}
	//先去缓存里找，找不到的话去mysql里找
	key := fmt.Sprintf("%s:%s:%s", consts.RedisDatabase, consts.REDIS_KEY_MEMBER, username)
	result, err := global.RedisDb.Get(key).Result()
	if err != nil && err != redis.Nil {
		return nil, errors.New("query user from redis failed:%v" + err.Error())
	}
	//缓存里有就返回
	if len(result) != 0 {
		//将result转换为user类型然后返回
		if err := json.Unmarshal([]byte(result), user); err != nil {
			global.Logger.Errorf("redis中的用户信息转换为response failed:%v", err)
		}
		return user, nil
	}
	//缓存里没有就去查数据库并且放入缓存中去
	if err := global.Db.Model(&users.User{}).Where("username = ?", username).Find(&user).Error; err != nil {
		return nil, errors.New("query user from mysql failed:%v" + err.Error())
	}
	//判断一下用户是否存在
	if user.ID == 0 {
		return nil, errors.New("query user from mysql failed:user not exist")
	}
	//放入缓存中去
	bytes, _ := json.Marshal(user)
	if err := global.RedisDb.Set(key, string(bytes), 5*time.Minute).Err(); err != nil {
		global.Logger.Errorf("add userinfo to redis failed:%v", err)
	}
	return user, nil
}
