package users

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gomall/consts"
	controller "gomall/controllers"
	receive "gomall/interaction/receive/users"
	"gomall/logic/userCache"
	"gomall/logic/users"
	userModels "gomall/models/users"
	"gomall/utils/jwt"
)

type LoginController struct {
	controller.BaseControllers
}

func (c *LoginController) Login(ctx *gin.Context) {
	var rec receive.UserLoginReceiveStruct
	rec.Username = ctx.PostForm("username")
	rec.Password = ctx.PostForm("password")
	//if err := ctx.ShouldBind(rec); err != nil {
	//	c.Response(ctx, "用户登录时，绑定参数失败", nil, err)
	//	global.Logger.Errorf("用户登录时，绑定参数失败:%v", err)
	//	return
	//}
	//先是去缓存里查询了以下username对应的用户信息，然后比对密码
	user, err := userCache.LoadUserByUsername(rec.Username)
	if err != nil {
		c.Response(ctx, "登录失败", nil, errors.New("query userinfo failed:"+err.Error()))
		return
	}
	if user == nil {
		c.Response(ctx, "登录失败", nil, errors.New("查询用户信息失败"))
		return
	}
	//然后生成token
	token := jwt.NextToken(user.Id)
	if token == "" {
		c.Response(ctx, "token生成失败", nil, errors.New("token generate failed"))
	}
	//然后向当前context头里添加了token信息
	ctx.Request.Header.Set("Authorization", token)
	responseData := map[string]string{
		"token":     token,
		"tokenHead": consts.TOKEN_HEAD,
	}
	//最后返回token
	c.Response(ctx, "操作成功", responseData, nil)
}

func (c *LoginController) Info(ctx *gin.Context) {
	member := &userModels.User{}
	memberId, err := jwt.GetMemberIdFromCtx(ctx)
	if err != nil {
		c.Response(ctx, "从ctx获取memberId出错", nil, err)
		return
	}
	if err := member.GetMemberById(memberId); err != nil {
		c.Response(ctx, "获取用户信息失败", nil, err)
		return
	}
	c.Response(ctx, "获取用户信息成功", member, nil)
}

func (c *LoginController) Register(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.UserRegisterStruct)); err == nil {
		results, err := users.Register(rec)
		if err != nil {
			c.Response(ctx, "注册失败", nil, errors.New("register failed: "+err.Error()))
		}
		c.Response(ctx, "注册成功", results, nil)
	}
}

func (c *LoginController) GenerateAuthCode(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.UserGetAuthCodeStruct)); err == nil {
		code, err := users.GetAuthCode(rec)
		if err != nil {
			c.Response(ctx, "获取验证码错误", nil, err)
		}
		c.Response(ctx, "获取验证码成功", code, nil)
	}
}

func (c *LoginController) UpdatePassword(ctx *gin.Context) {
	if rec, err := controller.ShouldBind(ctx, new(receive.UserUpdatePasswordStruct)); err == nil {
		results, err := users.UpdatePassword(rec)
		if err != nil {
			c.Response(ctx, "密码修改失败", nil, err)
		}
		c.Response(ctx, "密码修改成功", results, nil)
	}
}

// RefreshToken 当token没有过期的时候，可以刷新token
func (c *LoginController) RefreshToken(ctx *gin.Context) {
	//从context头里取出token，验证是否过期
	newToken, err := users.RefreshToken(ctx)
	if err != nil {
		c.Response(ctx, "token过期", nil, err)
	}

	responseData := map[string]string{
		"token":     newToken,
		"tokenHead": consts.TOKEN_HEAD,
	}
	c.Response(ctx, "token刷新成功", responseData, nil)
}
