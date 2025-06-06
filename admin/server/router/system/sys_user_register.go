package system

import (
	"github.com/gin-gonic/gin"
)

type UserRegisterRouter struct{}

func (s *UserRegisterRouter) InitUserRegisterRouter(Router *gin.RouterGroup) {
	userRegisterRouter := Router.Group("user")
	{
		userRegisterRouter.POST("sendEmailCode", userRegisterApi.SendEmailCode)   // 发送邮箱验证码
		userRegisterRouter.POST("selfRegister", userRegisterApi.SelfRegister)     // 用户自主注册
		userRegisterRouter.GET("checkUsername", userRegisterApi.CheckUsername)    // 检查用户名是否可用
		userRegisterRouter.GET("checkEmail", userRegisterApi.CheckEmail)          // 检查邮箱是否可用
	}
} 