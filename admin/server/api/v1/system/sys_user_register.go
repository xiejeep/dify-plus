package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	systemService "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserRegisterApi struct{}

var userRegisterApi = UserRegisterApi{}

// SendEmailCode
// @Tags     UserRegister
// @Summary  发送邮箱验证码
// @Produce  application/json
// @Param    data  body      systemReq.SendEmailCode                           true  "邮箱和验证码信息"
// @Success  200   {object}  response.Response{msg=string}                     "发送成功"
// @Router   /user/sendEmailCode [post]
func (a *UserRegisterApi) SendEmailCode(c *gin.Context) {
	var req systemReq.SendEmailCode
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 验证图片验证码
	if req.Captcha == "" || req.CaptchaId == "" {
		response.FailWithMessage("请输入图片验证码", c)
		return
	}

	if !store.Verify(req.CaptchaId, req.Captcha, true) {
		response.FailWithMessage("图片验证码错误", c)
		return
	}

	// 发送邮箱验证码
	err = systemService.UserRegisterServiceApp.SendEmailCode(req, c.ClientIP())
	if err != nil {
		global.GVA_LOG.Error("发送邮箱验证码失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("验证码发送成功", c)
}

// SelfRegister
// @Tags     UserRegister
// @Summary  用户自主注册
// @Produce  application/json
// @Param    data  body      systemReq.SelfRegister                               true  "注册信息"
// @Success  200   {object}  response.Response{data=systemRes.SysUserResponse,msg=string}  "注册成功"
// @Router   /user/selfRegister [post]
func (a *UserRegisterApi) SelfRegister(c *gin.Context) {
	var req systemReq.SelfRegister
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 验证请求参数
	err = utils.Verify(req, utils.RegisterVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 验证图片验证码
	if req.Captcha == "" || req.CaptchaId == "" {
		response.FailWithMessage("请输入图片验证码", c)
		return
	}

	if !store.Verify(req.CaptchaId, req.Captcha, true) {
		response.FailWithMessage("图片验证码错误", c)
		return
	}

	// 用户自主注册
	userReturn, err := systemService.UserRegisterServiceApp.SelfRegister(req)
	if err != nil {
		global.GVA_LOG.Error("用户自主注册失败!", zap.Error(err))
		response.FailWithDetailed(systemRes.SysUserResponse{User: userReturn}, err.Error(), c)
		return
	}

	response.OkWithDetailed(systemRes.SysUserResponse{User: userReturn}, "注册成功", c)
}

// CheckUsername
// @Tags     UserRegister
// @Summary  检查用户名是否可用
// @Produce  application/json
// @Param    username  query     string                                          true  "用户名"
// @Success  200       {object}  response.Response{data=bool,msg=string}         "检查结果"
// @Router   /user/checkUsername [get]
func (a *UserRegisterApi) CheckUsername(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		response.FailWithMessage("用户名不能为空", c)
		return
	}

	var existUser system.SysUser
	err := global.GVA_DB.Where("username = ?", username).First(&existUser).Error
	
	available := err != nil // 如果查询出错（即未找到），说明用户名可用
	
	if available {
		response.OkWithDetailed(gin.H{"available": true}, "用户名可用", c)
	} else {
		response.OkWithDetailed(gin.H{"available": false}, "用户名已存在", c)
	}
}

// CheckEmail
// @Tags     UserRegister
// @Summary  检查邮箱是否可用
// @Produce  application/json
// @Param    email  query     string                                          true  "邮箱"
// @Success  200    {object}  response.Response{data=bool,msg=string}         "检查结果"
// @Router   /user/checkEmail [get]
func (a *UserRegisterApi) CheckEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		response.FailWithMessage("邮箱不能为空", c)
		return
	}

	var existUser system.SysUser
	err := global.GVA_DB.Where("email = ?", email).First(&existUser).Error
	
	available := err != nil // 如果查询出错（即未找到），说明邮箱可用
	
	if available {
		response.OkWithDetailed(gin.H{"available": true}, "邮箱可用", c)
	} else {
		response.OkWithDetailed(gin.H{"available": false}, "邮箱已注册", c)
	}
} 