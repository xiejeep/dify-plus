package system

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// SysEmailVerification 邮箱验证码表
type SysEmailVerification struct {
	global.GVA_MODEL
	Email     string    `json:"email" gorm:"index;comment:邮箱地址"`
	Code      string    `json:"code" gorm:"comment:验证码"`
	Type      int       `json:"type" gorm:"default:1;comment:验证码类型 1:注册 2:找回密码"`
	ExpiredAt time.Time `json:"expired_at" gorm:"comment:过期时间"`
	Used      bool      `json:"used" gorm:"default:false;comment:是否已使用"`
	IP        string    `json:"ip" gorm:"comment:请求IP"`
}

func (SysEmailVerification) TableName() string {
	return "sys_email_verifications"
}

// IsExpired 检查验证码是否过期
func (e *SysEmailVerification) IsExpired() bool {
	return time.Now().After(e.ExpiredAt)
}

// IsValid 检查验证码是否有效（未过期且未使用）
func (e *SysEmailVerification) IsValid() bool {
	return !e.Used && !e.IsExpired()
} 