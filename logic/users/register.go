package users

import (
	"crypto/md5"
	"errors"
	"fmt"
	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gomall/consts"
	"gomall/global"
	receive "gomall/interaction/receive/users"
	"gomall/models/users"
	"gomall/utils/jwt"
	"math/rand"
	"time"
)

func Register(data *receive.UserRegisterStruct) (results interface{}, err error) {
	//先核验验证码是否正确。redis中验证码的key为database:authCode:telephone
	code, _ := global.RedisDb.Get(fmt.Sprintf("%s:%s:%s", consts.RedisDatabase, consts.REDIS_KEY_AUTH_CODE, data.Telephone)).Result()
	if data.AuthCode != code {
		return nil, errors.New("输入的验证码错误")
	}
	//然后查询库中是否已经存在该用户
	user := &users.User{}
	if user.IsExistByField("username", data.Username) {
		return nil, errors.New(fmt.Sprintf("%s已经被注册", data.Username))
	}
	//没有的话，插入新用户。获取用户的默认会员等级并设置
	user.Username = data.Username
	user.Phone = data.Telephone
	user.Password = encodePassword(data.Password)
	if err = global.Db.Create(user).Error; err != nil {
		return nil, errors.New("insert user to db failed:" + err.Error())
	}
	return "用户注册成功", nil
}

func GetAuthCode(data *receive.UserGetAuthCodeStruct) (results interface{}, err error) {
	authCode := generateAuthCode()
	//key:REDIS_DATABASE + ":" + REDIS_KEY_AUTH_CODE + ":" + telephone;
	var authCodeKey = fmt.Sprintf("%s:%s:%s", consts.RedisDatabase, consts.REDIS_KEY_AUTH_CODE, data.Telephone)
	if err = global.RedisDb.Set(authCodeKey, authCode, 5*time.Minute).Err(); err != nil {
		//如果没有写入redis的话会导致很严重的后果，后续注册用户验证时会取不到验证码，此时需要返回错误
		return nil, errors.New("generate auth code failed:" + err.Error())
	}
	return authCode, nil
}

// generateAuthCode 类内方法，生成6位随机验证码，以当前时间戳为seed
func generateAuthCode() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(1000000))
}

func UpdatePassword(data *receive.UserUpdatePasswordStruct) (results interface{}, err error) {
	//先判断当前telephone对应的用户是否存在
	user := &users.User{}
	if !user.IsExistByField("phone", data.Telephone) {
		return nil, errors.New("该手机号不存在")
	}
	//然后验证传过来的验证码是否等于redis中的验证码
	code, _ := global.RedisDb.Get(fmt.Sprintf("%s:%s:%s", consts.RedisDatabase, consts.REDIS_KEY_AUTH_CODE, data.Telephone)).Result()
	if data.AuthCode != code {
		return nil, errors.New("验证码错误")
	}
	//然后将该用户的密码更新为传过来的新密码
	if err := global.Db.Model(&users.User{}).Where("phone = ?", data.Telephone).Update("password", encodePassword(data.Password)).Error; err != nil {
		return nil, errors.New("update pwd failed:" + err.Error())
	}
	return "修改密码成功", nil
}

// encodePassword 将密码进行加密存储的工具方法
func encodePassword(password string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}

func RefreshToken(ctx *gin.Context) (newToken string, err error) {
	token := ctx.Request.Header.Get("Authorization")
	//先检验有没有token
	if len(token) == 0 {
		return "", errors.New("token过期")
	}

	//然后检验token是否过期
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return "", err
	}

	//更新有效期，重新生成token
	claims.ExpiresAt = time.Now().Add(7 * 24 * time.Hour).Unix()

	unsigned := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
	tokenString, _ := unsigned.SignedString(jwt.Hotkey)

	//重写ctx里的token
	ctx.Request.Header.Set("Authorization", tokenString)
	return tokenString, nil
}
