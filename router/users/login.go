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
		loginRouter.POST("/register", loginController.Register)
		loginRouter.POST("/getAuthCode", loginController.GenerateAuthCode)
		loginRouter.POST("/updatePassword", loginController.UpdatePassword)
		loginRouter.POST("/refreshToken", loginController.RefreshToken)
	}
}
