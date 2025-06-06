import service from '@/utils/request'

// 发送邮箱验证码
export const sendEmailCode = (data) => {
  return service({
    url: '/user/sendEmailCode',
    method: 'post',
    data
  })
}

// 用户自主注册
export const selfRegister = (data) => {
  return service({
    url: '/user/selfRegister',
    method: 'post',
    data
  })
}

// 检查用户名是否可用
export const checkUsername = (params) => {
  return service({
    url: '/user/checkUsername',
    method: 'get',
    params
  })
}

// 检查邮箱是否可用
export const checkEmail = (params) => {
  return service({
    url: '/user/checkEmail',
    method: 'get',
    params
  })
} 