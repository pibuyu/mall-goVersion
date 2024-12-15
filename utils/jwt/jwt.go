package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gomall/consts"
	"gomall/global"
	"gomall/models/users"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Hotkey 密钥
var Hotkey = []byte("haifengonline.top")

// Claims  TOKEN 的结构体
type Claims struct {
	UserID    int64
	LoginTime string
	User      users.User `json:"user"`
	jwt.StandardClaims
}

// NextToken 就是返回一个token
func NextToken(uid int64) string {

	expireTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &Claims{
		UserID:    uid,
		LoginTime: time.Now().Format(time.DateTime),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "root",       // 签名颁发者
			Subject:   "user token", //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(Hotkey)
	if err != nil {
		log.Fatalf("签发token时出错:%v", err)
		return ""
	}

	//签发token的同时往redis里放一份;然后在中间件查看收到的token和redis中的最新token是否一致
	key := fmt.Sprintf("%d_%s", uid, consts.TokenString)
	global.RedisDb.Set(key, tokenString, 7*24*time.Hour)
	global.Logger.Infof("给用户 %d 签发并向redis投递了token：%s", uid, tokenString)
	return tokenString
}

// ParseToken 解析 Token
func ParseToken(tokenStr string) (*Claims, error) {
	//todo:前端放进去的token的格式为：response.data.tokenHead+response.data.token;也就是在token前面加上了  Bearer 前缀 ,解析的时候去掉就行
	tokenStr = tokenStr[6:]
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return Hotkey, nil
	})
	if err != nil {
		global.Logger.Errorf("token parse err : " + err.Error())
		return nil, errors.New("token parse err : " + err.Error())
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func GetMemberIdFromCtx(ctx *gin.Context) (int64, error) {
	//注意：用postman进行测试的时候，token是Bearer+token的格式。不然解析会出问题
	token := ctx.Request.Header.Get("Authorization")

	claims, err := ParseToken(token)
	if err != nil {
		return 0, errors.New("从ctx中获取MemberId出错:" + err.Error())
	}
	return claims.UserID, nil
}
