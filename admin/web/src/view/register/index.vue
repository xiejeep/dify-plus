<template>
  <div class="register-container">
    <div class="register-box">
      <div class="register-header">
        <h2>用户注册</h2>
        <p>欢迎注册 Dify-Plus</p>
      </div>
      
      <el-form
        ref="registerFormRef"
        :model="registerForm"
        :rules="registerRules"
        label-width="80px"
        class="register-form"
      >
        <!-- 用户名 -->
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="registerForm.username"
            placeholder="请输入用户名"
            maxlength="20"
            show-word-limit
            @blur="checkUsernameAvailable"
          />
        </el-form-item>

        <!-- 昵称 -->
        <el-form-item label="昵称" prop="nickName">
          <el-input
            v-model="registerForm.nickName"
            placeholder="请输入昵称"
            maxlength="20"
            show-word-limit
          />
        </el-form-item>

        <!-- 邮箱 -->
        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="registerForm.email"
            placeholder="请输入邮箱地址"
            @blur="checkEmailAvailable"
          />
        </el-form-item>

        <!-- 手机号 -->
        <el-form-item label="手机号" prop="phone">
          <el-input
            v-model="registerForm.phone"
            placeholder="请输入手机号（可选）"
            maxlength="11"
          />
        </el-form-item>

        <!-- 密码 -->
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="registerForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>

        <!-- 确认密码 -->
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="registerForm.confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            show-password
          />
        </el-form-item>

        <!-- 图片验证码 -->
        <el-form-item label="验证码" prop="captcha">
          <div class="captcha-container">
            <el-input
              v-model="registerForm.captcha"
              placeholder="请输入图片验证码"
              style="width: 200px"
            />
            <img
              :src="captchaData.picPath"
              alt="验证码"
              class="captcha-img"
              @click="getCaptcha"
            />
          </div>
        </el-form-item>

        <!-- 邮箱验证码 -->
        <el-form-item label="邮箱验证码" prop="emailCode">
          <div class="email-code-container">
            <el-input
              v-model="registerForm.emailCode"
              placeholder="请输入邮箱验证码"
              style="width: 200px"
            />
            <el-button
              :disabled="emailCodeDisabled"
              :loading="emailCodeLoading"
              type="primary"
              @click="sendEmailCodeHandler"
            >
              {{ emailCodeText }}
            </el-button>
          </div>
        </el-form-item>

        <!-- 注册按钮 -->
        <el-form-item>
          <el-button
            :loading="registerLoading"
            type="primary"
            size="large"
            style="width: 100%"
            @click="handleRegister"
          >
            注册
          </el-button>
        </el-form-item>

        <!-- 返回登录 -->
        <el-form-item>
          <div class="back-to-login">
            <span>已有账号？</span>
            <el-button type="text" @click="goToLogin">返回登录</el-button>
          </div>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { captcha } from '@/api/user'
import { sendEmailCode, selfRegister, checkUsername, checkEmail } from '@/api/userRegister'

const router = useRouter()

// 表单引用
const registerFormRef = ref()

// 表单数据
const registerForm = reactive({
  username: '',
  nickName: '',
  email: '',
  phone: '',
  password: '',
  confirmPassword: '',
  captcha: '',
  captchaId: '',
  emailCode: ''
})

// 验证码数据
const captchaData = reactive({
  picPath: '',
  captchaId: '',
  openCaptcha: false
})

// 邮箱验证码相关
const emailCodeDisabled = ref(false)
const emailCodeLoading = ref(false)
const emailCodeText = ref('发送验证码')
const emailCodeCountdown = ref(0)

// 注册loading
const registerLoading = ref(false)

// 表单验证规则
const registerRules = reactive({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  nickName: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 1, max: 20, message: '昵称长度在 1 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
    { 
      pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d@$!%*?&]{6,}$/, 
      message: '密码必须包含大小写字母和数字', 
      trigger: 'blur' 
    }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== registerForm.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  captcha: [
    { required: true, message: '请输入图片验证码', trigger: 'blur' }
  ],
  emailCode: [
    { required: true, message: '请输入邮箱验证码', trigger: 'blur' },
    { len: 6, message: '邮箱验证码为6位数字', trigger: 'blur' }
  ]
})

// 获取图片验证码
const getCaptcha = async () => {
  try {
    const res = await captcha()
    if (res.code === 0) {
      captchaData.picPath = res.data.picPath
      captchaData.captchaId = res.data.captchaId
      captchaData.openCaptcha = res.data.openCaptcha
      registerForm.captchaId = res.data.captchaId
    }
  } catch (error) {
    console.error('获取验证码失败:', error)
  }
}

// 检查用户名是否可用
const checkUsernameAvailable = async () => {
  if (!registerForm.username) return
  
  try {
    const res = await checkUsername({ username: registerForm.username })
    if (res.code === 0 && !res.data.available) {
      ElMessage.warning('用户名已存在，请换一个')
    }
  } catch (error) {
    console.error('检查用户名失败:', error)
  }
}

// 检查邮箱是否可用
const checkEmailAvailable = async () => {
  if (!registerForm.email) return
  
  try {
    const res = await checkEmail({ email: registerForm.email })
    if (res.code === 0 && !res.data.available) {
      ElMessage.warning('邮箱已注册，请换一个或直接登录')
    }
  } catch (error) {
    console.error('检查邮箱失败:', error)
  }
}

// 发送邮箱验证码
const sendEmailCodeHandler = async () => {
  // 先验证邮箱和图片验证码
  if (!registerForm.email) {
    ElMessage.warning('请先输入邮箱地址')
    return
  }
  
  if (!registerForm.captcha) {
    ElMessage.warning('请先输入图片验证码')
    return
  }

  try {
    emailCodeLoading.value = true
    
    const res = await sendEmailCode({
      email: registerForm.email,
      type: 1, // 注册验证码
      captchaId: registerForm.captchaId,
      captcha: registerForm.captcha
    })
    
    if (res.code === 0) {
      ElMessage.success('验证码发送成功，请查收邮件')
      
      // 开始倒计时
      emailCodeDisabled.value = true
      emailCodeCountdown.value = 60
      
      const timer = setInterval(() => {
        emailCodeCountdown.value--
        emailCodeText.value = `${emailCodeCountdown.value}s后重新发送`
        
        if (emailCodeCountdown.value <= 0) {
          clearInterval(timer)
          emailCodeDisabled.value = false
          emailCodeText.value = '发送验证码'
        }
      }, 1000)
    } else {
      ElMessage.error(res.msg || '发送失败')
      // 重新获取图片验证码
      getCaptcha()
    }
  } catch (error) {
    ElMessage.error('发送失败，请重试')
    getCaptcha()
  } finally {
    emailCodeLoading.value = false
  }
}

// 处理注册
const handleRegister = async () => {
  try {
    // 表单验证
    await registerFormRef.value.validate()
    
    registerLoading.value = true
    
    const res = await selfRegister({
      username: registerForm.username,
      nickName: registerForm.nickName,
      email: registerForm.email,
      phone: registerForm.phone,
      password: registerForm.password,
      emailCode: registerForm.emailCode,
      captchaId: registerForm.captchaId,
      captcha: registerForm.captcha
    })
    
    if (res.code === 0) {
      ElMessage.success('注册成功！请使用您的账号登录')
      router.push('/login')
    } else {
      ElMessage.error(res.msg || '注册失败')
      // 重新获取图片验证码
      getCaptcha()
    }
  } catch (error) {
    console.error('注册失败:', error)
    ElMessage.error('注册失败，请重试')
    getCaptcha()
  } finally {
    registerLoading.value = false
  }
}

// 返回登录页
const goToLogin = () => {
  router.push('/login')
}

// 页面加载时获取验证码
onMounted(() => {
  getCaptcha()
})
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-box {
  width: 480px;
  padding: 40px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.register-header h2 {
  margin: 0 0 8px 0;
  color: #333;
  font-size: 28px;
  font-weight: 600;
}

.register-header p {
  margin: 0;
  color: #666;
  font-size: 14px;
}

.register-form {
  margin-top: 20px;
}

.captcha-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.captcha-img {
  width: 120px;
  height: 40px;
  cursor: pointer;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}

.email-code-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.back-to-login {
  text-align: center;
  color: #666;
  font-size: 14px;
}

.back-to-login .el-button {
  padding: 0;
  font-size: 14px;
}
</style> 