package system

import (
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"math/big"
	"net/smtp"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	serviceGaia "github.com/flipped-aurora/gin-vue-admin/server/service/gaia"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gofrs/uuid/v5"
	"github.com/jordan-wright/email"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRegisterService struct{}

var UserRegisterServiceApp = new(UserRegisterService)

const (
	EmailCodeTypeRegister = 1 // 注册验证码
	EmailCodeTypeReset    = 2 // 找回密码验证码
)

// SendEmailCode 发送邮箱验证码
func (s *UserRegisterService) SendEmailCode(req systemReq.SendEmailCode, ip string) error {
	// 1. 检查邮箱是否已注册
	var existUser system.SysUser
	if !errors.Is(global.GVA_DB.Where("email = ?", req.Email).First(&existUser).Error, gorm.ErrRecordNotFound) {
		if req.Type == EmailCodeTypeRegister {
			return errors.New("该邮箱已注册")
		}
	} else {
		if req.Type == EmailCodeTypeReset {
			return errors.New("该邮箱未注册")
		}
	}

	// 2. 检查发送频率限制（1分钟内只能发送一次）
	var lastCode system.SysEmailVerification
	oneMinuteAgo := time.Now().Add(-1 * time.Minute)
	err := global.GVA_DB.Where("email = ? AND type = ? AND created_at > ?", req.Email, req.Type, oneMinuteAgo).
		First(&lastCode).Error
	if err == nil {
		// 找到了1分钟内的记录，说明发送过于频繁
		return errors.New("发送过于频繁，请稍后再试")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 数据库查询出错
		global.GVA_LOG.Error("查询邮箱验证码失败", zap.Error(err))
		return errors.New("系统错误，请稍后再试")
	}
	// err == gorm.ErrRecordNotFound 表示没有找到1分钟内的记录，可以发送

	// 3. 检查IP限制（每小时最多发送5次）
	var ipCount int64
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	global.GVA_DB.Model(&system.SysEmailVerification{}).
		Where("ip = ? AND created_at > ?", ip, oneHourAgo).Count(&ipCount)
	if ipCount >= 5 {
		return errors.New("发送次数过多，请稍后再试")
	}

	// 4. 生成6位随机验证码
	code, err := s.generateEmailCode()
	if err != nil {
		return err
	}

	// 5. 保存验证码到数据库
	emailVerification := system.SysEmailVerification{
		Email:     req.Email,
		Code:      code,
		Type:      req.Type,
		ExpiredAt: time.Now().Add(5 * time.Minute), // 5分钟过期
		IP:        ip,
	}

	if err := global.GVA_DB.Create(&emailVerification).Error; err != nil {
		global.GVA_LOG.Error("保存邮箱验证码失败", zap.Error(err))
		return errors.New("发送验证码失败")
	}

	// 6. 发送邮件
	subject := "Dify-Plus 验证码"
	var body string
	if req.Type == EmailCodeTypeRegister {
		body = fmt.Sprintf(`
			<h2>欢迎注册 Dify-Plus</h2>
			<p>您的注册验证码是：<strong style="color: #409EFF; font-size: 18px;">%s</strong></p>
			<p>验证码有效期为 5 分钟，请及时使用。</p>
			<p>如果这不是您的操作，请忽略此邮件。</p>
		`, code)
	} else {
		body = fmt.Sprintf(`
			<h2>Dify-Plus 密码重置</h2>
			<p>您的密码重置验证码是：<strong style="color: #409EFF; font-size: 18px;">%s</strong></p>
			<p>验证码有效期为 5 分钟，请及时使用。</p>
			<p>如果这不是您的操作，请忽略此邮件。</p>
		`, code)
	}

	// 调试：打印邮件配置信息
	global.GVA_LOG.Info("邮件配置信息", 
		zap.String("from", global.GVA_CONFIG.Email.From),
		zap.String("host", global.GVA_CONFIG.Email.Host),
		zap.String("to", req.Email),
		zap.Int("port", global.GVA_CONFIG.Email.Port),
		zap.Bool("isSSL", global.GVA_CONFIG.Email.IsSSL))

	// 直接发送邮件，不依赖插件
	if err := s.sendEmailDirect(req.Email, subject, body); err != nil {
		global.GVA_LOG.Error("发送邮件失败", zap.Error(err))
		return errors.New("发送验证码失败")
	}

	return nil
}

// VerifyEmailCode 验证邮箱验证码
func (s *UserRegisterService) VerifyEmailCode(req systemReq.VerifyEmailCode) error {
	var emailVerification system.SysEmailVerification
	err := global.GVA_DB.Where("email = ? AND code = ? AND type = ? AND used = ?", 
		req.Email, req.Code, req.Type, false).First(&emailVerification).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("验证码错误或已使用")
	}
	
	if err != nil {
		return err
	}

	if !emailVerification.IsValid() {
		return errors.New("验证码已过期")
	}

	// 标记验证码为已使用
	emailVerification.Used = true
	global.GVA_DB.Save(&emailVerification)

	return nil
}

// SelfRegister 用户自主注册
func (s *UserRegisterService) SelfRegister(req systemReq.SelfRegister) (userInter system.SysUser, err error) {
	// 1. 验证邮箱验证码
	if err = s.VerifyEmailCode(systemReq.VerifyEmailCode{
		Email: req.Email,
		Code:  req.EmailCode,
		Type:  EmailCodeTypeRegister,
	}); err != nil {
		return userInter, err
	}

	// 2. 检查用户名是否已存在
	var existUser system.SysUser
	if !errors.Is(global.GVA_DB.Where("username = ?", req.Username).First(&existUser).Error, gorm.ErrRecordNotFound) {
		return userInter, errors.New("用户名已存在")
	}

	// 3. 检查邮箱是否已注册
	if !errors.Is(global.GVA_DB.Where("email = ?", req.Email).First(&existUser).Error, gorm.ErrRecordNotFound) {
		return userInter, errors.New("邮箱已注册")
	}

	// 4. 密码强度验证
	if err = serviceGaia.IsUserPasswordValid(req.Password); err != nil {
		return userInter, err
	}

	// 5. 创建用户
	user := system.SysUser{
		Username:    req.Username,
		NickName:    req.NickName,
		Password:    req.Password,
		Email:       req.Email,
		Phone:       req.Phone,
		AuthorityId: system.DefaultGroupID, // 默认为普通用户角色
		Enable:      system.UserActive,     // 默认启用
	}

	// 设置默认头像
	user.HeaderImg = "https://qmplusimg.henrongyi.top/1576554439myAvatar.png"

	// 6. 同步到Dify平台（使用空token，因为是自主注册）
	if err = serviceGaia.RegisterUser(user, ""); err != nil {
		return userInter, errors.New("同步到Dify平台失败:" + err.Error())
	}

	// 7. 密码加密并保存到数据库
	user.Password = utils.BcryptHash(user.Password)
	user.UUID = uuid.Must(uuid.NewV4())

	// 设置用户权限
	var authorities []system.SysAuthority
	authorities = append(authorities, system.SysAuthority{
		AuthorityId: system.DefaultGroupID,
	})
	user.Authorities = authorities

	err = global.GVA_DB.Create(&user).Error
	if err != nil {
		global.GVA_LOG.Error("用户注册失败", zap.Error(err))
		return userInter, errors.New("注册失败")
	}

	global.GVA_LOG.Info("用户自主注册成功", zap.String("username", user.Username), zap.String("email", user.Email))
	
	return user, nil
}

// generateEmailCode 生成6位随机数字验证码
func (s *UserRegisterService) generateEmailCode() (string, error) {
	code := ""
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += n.String()
	}
	return code, nil
}

// CleanExpiredCodes 清理过期的验证码（可以作为定时任务）
func (s *UserRegisterService) CleanExpiredCodes() error {
	return global.GVA_DB.Where("expired_at < ?", time.Now()).Delete(&system.SysEmailVerification{}).Error
}

// sendEmailDirect 直接发送邮件，使用系统配置
func (s *UserRegisterService) sendEmailDirect(to, subject, body string) error {
	config := global.GVA_CONFIG.Email
	
	// 使用resend的SMTP配置
	auth := smtp.PlainAuth("", "resend", config.Secret, config.Host)
	
	e := email.NewEmail()
	if config.Nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", config.Nickname, config.From)
	} else {
		e.From = config.From
	}
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(body)
	
	hostAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	
	global.GVA_LOG.Info("准备发送邮件", 
		zap.String("from", e.From),
		zap.Strings("to", e.To),
		zap.String("hostAddr", hostAddr))
	
	var err error
	if config.IsSSL {
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: config.Host})
	} else {
		err = e.Send(hostAddr, auth)
	}
	
	return err
} 