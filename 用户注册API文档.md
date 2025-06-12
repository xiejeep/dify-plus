# 用户注册API文档

## 概述

本文档描述了用户注册的完整流程，包含三个步骤：
1. 获取图形验证码
2. 发送邮箱验证码  
3. 完成用户注册

**基础URL**: `http://app.classhorse.top:8081/api`

## API接口列表

### 1. 获取图形验证码

#### 接口信息
- **接口地址**: `/base/captcha`
- **请求方式**: `POST`
- **接口描述**: 获取图形验证码，用于后续发送邮箱验证码和用户注册时的安全验证

#### 请求参数
无需传入参数，直接POST请求即可。

#### 请求示例
```bash
curl -X POST http://app.classhorse.top:8081/api/base/captcha \
  -H "Content-Type: application/json"
```

#### 响应参数
| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| code | int | 是 | 响应状态码，7000表示成功 |
| data | object | 是 | 验证码数据 |
| msg | string | 是 | 响应消息 |

**data对象结构**:
| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| captchaId | string | 是 | 验证码ID，后续请求需要携带 |
| picPath | string | 是 | 验证码图片的base64编码 |
| captchaLength | int | 是 | 验证码长度 |
| openCaptcha | bool | 是 | 是否开启验证码验证 |

#### 响应示例
```json
{
  "code": 7000,
  "data": {
    "captchaId": "nP6eUXrZTBjsXNEJ",
    "picPath": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAMgAAABkCAYAAADDhn8LAAA...",
    "captchaLength": 5,
    "openCaptcha": true
  },
  "msg": "验证码获取成功"
}
```

### 2. 发送邮箱验证码

#### 接口信息
- **接口地址**: `/user/sendEmailCode`
- **请求方式**: `POST`
- **接口描述**: 向指定邮箱发送验证码，需要先通过图形验证码验证

#### 请求参数
| 参数名 | 类型 | 必填 | 描述 | 示例 |
|--------|------|------|------|------|
| email | string | 是 | 邮箱地址 | user@example.com |
| type | int | 是 | 验证码类型，1:注册 2:找回密码 | 1 |
| captchaId | string | 是 | 图形验证码ID | nP6eUXrZTBjsXNEJ |
| captcha | string | 是 | 图形验证码内容 | 12345 |

#### 请求示例
```bash
curl -X POST http://app.classhorse.top:8081/api/user/sendEmailCode \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "type": 1,
    "captchaId": "nP6eUXrZTBjsXNEJ",
    "captcha": "12345"
  }'
```

#### 响应参数
| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| code | int | 是 | 响应状态码，7000表示成功 |
| data | object | 否 | 响应数据（成功时为空） |
| msg | string | 是 | 响应消息 |

#### 响应示例
**成功响应**:
```json
{
  "code": 7000,
  "data": null,
  "msg": "验证码发送成功"
}
```

**失败响应**:
```json
{
  "code": 7001,
  "data": null,
  "msg": "图片验证码错误"
}
```

### 3. 用户注册

#### 接口信息
- **接口地址**: `/user/selfRegister`
- **请求方式**: `POST`
- **接口描述**: 用户自主注册，需要提供完整的注册信息和验证码

#### 请求参数
| 参数名 | 类型 | 必填 | 描述 | 约束 | 示例 |
|--------|------|------|------|------|------|
| username | string | 是 | 用户名 | 3-20个字符 | johndoe |
| password | string | 是 | 密码 | 最少6位 | password123 |
| nickName | string | 是 | 昵称 | 1-20个字符 | 约翰 |
| email | string | 是 | 邮箱地址 | 有效邮箱格式 | user@example.com |
| phone | string | 否 | 手机号 | | 13800138000 |
| emailCode | string | 是 | 邮箱验证码 | 6位数字 | 123456 |
| captchaId | string | 是 | 图形验证码ID | | nP6eUXrZTBjsXNEJ |
| captcha | string | 是 | 图形验证码内容 | | 12345 |

#### 请求示例
```bash
curl -X POST http://app.classhorse.top:8081/api/user/selfRegister \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "password123",
    "nickName": "约翰",
    "email": "user@example.com",
    "phone": "13800138000",
    "emailCode": "123456",
    "captchaId": "nP6eUXrZTBjsXNEJ",
    "captcha": "12345"
  }'
```

#### 响应参数
| 参数名 | 类型 | 必填 | 描述 |
|--------|------|------|------|
| code | int | 是 | 响应状态码，7000表示成功 |
| data | object | 是 | 用户信息（成功时） |
| msg | string | 是 | 响应消息 |

**data对象结构**（成功时）:
| 参数名 | 类型 | 描述 |
|--------|------|------|
| user | object | 用户基本信息 |

#### 响应示例
**成功响应**:
```json
{
  "code": 7000,
  "data": {
    "user": {
      "id": 1,
      "username": "johndoe",
      "nickName": "约翰",
      "email": "user@example.com",
      "phone": "13800138000",
      "createdAt": "2024-01-01T10:00:00Z"
    }
  },
  "msg": "注册成功"
}
```

**失败响应**:
```json
{
  "code": 7001,
  "data": null,
  "msg": "用户名已存在"
}
```

## 完整注册流程示例

### 步骤1: 获取图形验证码
```javascript
// 获取图形验证码
const getCaptcha = async () => {
  const response = await fetch('http://app.classhorse.top:8081/api/base/captcha', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  const result = await response.json();
  
  if (result.code === 7000) {
    // 保存captchaId，显示验证码图片
    const captchaId = result.data.captchaId;
    const captchaImage = result.data.picPath;
    // 显示验证码图片让用户输入
  }
};
```

### 步骤2: 发送邮箱验证码
```javascript
// 发送邮箱验证码
const sendEmailCode = async (email, captchaId, captcha) => {
  const response = await fetch('http://app.classhorse.top:8081/api/user/sendEmailCode', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      email: email,
      type: 1, // 注册类型
      captchaId: captchaId,
      captcha: captcha
    }),
  });
  const result = await response.json();
  
  if (result.code === 7000) {
    // 验证码发送成功，提示用户查收邮件
    alert('验证码已发送到您的邮箱');
  }
};
```

### 步骤3: 完成注册
```javascript
// 用户注册
const register = async (userData) => {
  const response = await fetch('http://app.classhorse.top:8081/api/user/selfRegister', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      username: userData.username,
      password: userData.password,
      nickName: userData.nickName,
      email: userData.email,
      phone: userData.phone,
      emailCode: userData.emailCode,
      captchaId: userData.captchaId,
      captcha: userData.captcha
    }),
  });
  const result = await response.json();
  
  if (result.code === 7000) {
    // 注册成功
    alert('注册成功！');
    // 可以跳转到登录页面或自动登录
  } else {
    // 注册失败，显示错误信息
    alert(result.msg);
  }
};
```

## 错误码说明

| 错误码 | 描述 | 解决方案 |
|--------|------|----------|
| 7000 | 成功 | - |
| 7001 | 请求失败 | 检查请求参数和格式 |
| 400 | 参数错误 | 检查必填参数是否完整 |
| 429 | 请求过于频繁 | 等待一段时间后重试 |
| 500 | 服务器内部错误 | 联系技术支持 |

## 注意事项

1. **图形验证码有效期**: 图形验证码有时间限制，建议在获取后尽快使用
2. **邮箱验证码有效期**: 邮箱验证码通常有效期为5-10分钟
3. **IP限制**: 同一IP发送邮箱验证码有频率限制，避免短时间内重复请求
4. **用户名唯一性**: 用户名和邮箱都需要保证唯一性
5. **密码安全**: 建议密码包含字母、数字，长度至少6位
6. **请求头**: 所有请求都需要设置`Content-Type: application/json`

## 移动端集成建议

1. **UI流程设计**:
   - 第一步：显示注册表单，包含图形验证码
   - 第二步：用户输入邮箱后，点击"发送验证码"按钮
   - 第三步：用户输入邮箱验证码，完成注册

2. **错误处理**:
   - 对各种错误码进行友好的错误提示
   - 网络请求失败时提供重试机制

3. **用户体验优化**:
   - 添加loading状态指示器
   - 验证码倒计时功能
   - 表单验证和实时提示

4. **安全考虑**:
   - 不要在客户端存储敏感信息
   - 使用HTTPS确保数据传输安全
   - 添加重复提交防护 