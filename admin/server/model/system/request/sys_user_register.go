package request

// SelfRegister 用户自主注册请求
type SelfRegister struct {
	Username    string `json:"username" binding:"required,min=3,max=20" example:"用户名"`
	Password    string `json:"password" binding:"required,min=6" example:"密码"`
	NickName    string `json:"nickName" binding:"required,min=1,max=20" example:"昵称"`
	Email       string `json:"email" binding:"required,email" example:"邮箱地址"`
	Phone       string `json:"phone" example:"手机号"`
	EmailCode   string `json:"emailCode" binding:"required,len=6" example:"邮箱验证码"`
	CaptchaId   string `json:"captchaId" binding:"required" example:"图片验证码ID"`
	Captcha     string `json:"captcha" binding:"required" example:"图片验证码"`
}

// SendEmailCode 发送邮箱验证码请求
type SendEmailCode struct {
	Email     string `json:"email" binding:"required,email" example:"邮箱地址"`
	Type      int    `json:"type" binding:"required,oneof=1 2" example:"验证码类型 1:注册 2:找回密码"`
	CaptchaId string `json:"captchaId" binding:"required" example:"图片验证码ID"`
	Captcha   string `json:"captcha" binding:"required" example:"图片验证码"`
}

// VerifyEmailCode 验证邮箱验证码请求
type VerifyEmailCode struct {
	Email string `json:"email" binding:"required,email" example:"邮箱地址"`
	Code  string `json:"code" binding:"required,len=6" example:"验证码"`
	Type  int    `json:"type" binding:"required,oneof=1 2" example:"验证码类型 1:注册 2:找回密码"`
} 