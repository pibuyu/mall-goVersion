package userCache

import (
	"errors"
	"gomall/models/users"
)

// LoadUserByUsername 查找用户信息
func LoadUserByUsername(username string) (user *users.User, err error) {
	//记得先初始化，不然json.Unmarshal会报错：将string解析到一个空指针上
	user = &users.User{}
	if err := user.GetByName(username); err != nil {
		return nil, errors.New("查找用户信息failed:" + err.Error())
	}
	return user, nil
}
