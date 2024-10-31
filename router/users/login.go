package users

import (
	"github.com/gin-gonic/gin"
	"gomall/controllers/users"
)

type LoginRouter struct {
}

// locahost:9090/sso/register
func (s *LoginRouter) InitLoginRouter(Router *gin.RouterGroup) {
	loginRouter := Router.Group("sso").Use()
	{
		loginController := new(users.LoginController)
		loginRouter.POST("/login", loginController.Login)
		loginRouter.GET("/info", loginController.Info) //前端执行完login之后会立马执行一次info
		loginRouter.POST("/register", loginController.Register)
		loginRouter.POST("/getAuthCode", loginController.GenerateAuthCode)
		loginRouter.POST("/updatePassword", loginController.UpdatePassword)
		loginRouter.POST("/refreshToken", loginController.RefreshToken)
	}
}
