package test

import (
	"encoding/json"
	"gomall/global"
	"gomall/models/users"
	"testing"
)

func TestGorm(t *testing.T) {
	var user users.User
	if err := global.Db.Model(&users.User{}).Where("id = ?", 1).Find(&user).Error; err != nil {
		t.Errorf("sql err：%v", err)
	}
	bytes, _ := json.Marshal(user)
	t.Logf("user info:%s", string(bytes))
}
