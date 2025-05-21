# 移动端注册API调用指南

本文档详细说明了如何从移动端应用调用Dify-Plus系统的注册API。

## 注册流程概述

注册流程基于邮箱验证机制，主要包含以下步骤：

1. 用户提供邮箱地址
2. 系统发送验证码到该邮箱
3. 用户验证邮箱所有权
4. 用户设置账号信息完成注册
5. 用户使用新账号登录系统

## API端点详解

### 1. 发送注册验证码

这是注册流程的第一步，系统会向用户邮箱发送验证码。

**请求：**

```http
POST /console/api/register
Content-Type: application/json

{
  "email": "用户邮箱",
  "language": "zh-Hans"  // 可选，默认为en-US
}
```

**成功响应：**

```json
{
  "result": "success",
  "data": "token字符串"
}
```

**失败响应（账号已存在）：**

```json
{
  "result": "fail",
  "code": "account_already_exists",
  "message": "Account already exists"
}
```

### 2. 验证邮箱验证码

用户收到验证码后，需要提交验证码进行验证。

**请求：**

```http
POST /console/api/register/validity
Content-Type: application/json

{
  "email": "用户邮箱",
  "code": "验证码",
  "token": "上一步返回的token"
}
```

**成功响应：**

```json
{
  "is_valid": true,
  "email": "用户邮箱"
}
```

### 3. 完成注册

验证成功后，用户需要设置账号信息完成注册。

**请求：**

```http
POST /console/api/register/complete
Content-Type: application/json

{
  "token": "上一步的token",
  "name": "用户姓名",
  "password": "密码",
  "password_confirm": "确认密码"
}
```

**成功响应：**

```json
{
  "result": "success"
}
```

### 4. 登录获取访问令牌

注册完成后，用户可以使用新账号登录系统。

**请求：**

```http
POST /console/api/login
Content-Type: application/json

{
  "email": "用户邮箱",
  "password": "密码"
}
```

**成功响应：**

```json
{
  "result": "success",
  "data": {
    "access_token": "访问令牌",
    "refresh_token": "刷新令牌"
  }
}
```

## 移动端实现示例

### Android (Kotlin) 实现

```kotlin
import okhttp3.*
import org.json.JSONObject
import java.io.IOException

class DifyAuthService {
    private val client = OkHttpClient()
    private val baseUrl = "https://您的服务器地址/console/api"
    private val mediaType = MediaType.parse("application/json; charset=utf-8")

    // 第一步：发送注册验证码
    fun sendRegisterCode(email: String, callback: (Boolean, String?, String?) -> Unit) {
        val json = JSONObject().apply {
            put("email", email)
            put("language", "zh-Hans")
        }

        val requestBody = RequestBody.create(mediaType, json.toString())
        val request = Request.Builder()
            .url("$baseUrl/register")
            .post(requestBody)
            .build()

        client.newCall(request).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback(false, null, e.message)
            }

            override fun onResponse(call: Call, response: Response) {
                val responseBody = response.body()?.string()
                val jsonResponse = JSONObject(responseBody)

                val result = jsonResponse.getString("result")

                if (result == "success") {
                    val token = jsonResponse.getString("data")
                    callback(true, token, null)
                } else {
                    val code = if (jsonResponse.has("code")) jsonResponse.getString("code") else null
                    val message = if (jsonResponse.has("message")) jsonResponse.getString("message") else "注册失败"
                    callback(false, null, message)
                }
            }
        })
    }

    // 第二步：验证注册验证码
    fun verifyRegisterCode(email: String, code: String, token: String, callback: (Boolean, String?) -> Unit) {
        val json = JSONObject().apply {
            put("email", email)
            put("code", code)
            put("token", token)
        }

        val requestBody = RequestBody.create(mediaType, json.toString())
        val request = Request.Builder()
            .url("$baseUrl/register/validity")
            .post(requestBody)
            .build()

        client.newCall(request).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback(false, e.message)
            }

            override fun onResponse(call: Call, response: Response) {
                val responseBody = response.body()?.string()
                val jsonResponse = JSONObject(responseBody)

                val isValid = if (jsonResponse.has("is_valid")) jsonResponse.getBoolean("is_valid") else false

                callback(isValid, null)
            }
        })
    }

    // 第三步：完成注册
    fun completeRegistration(token: String, name: String, password: String, passwordConfirm: String, callback: (Boolean, String?) -> Unit) {
        val json = JSONObject().apply {
            put("token", token)
            put("name", name)
            put("password", password)
            put("password_confirm", passwordConfirm)
        }

        val requestBody = RequestBody.create(mediaType, json.toString())
        val request = Request.Builder()
            .url("$baseUrl/register/complete")
            .post(requestBody)
            .build()

        client.newCall(request).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback(false, e.message)
            }

            override fun onResponse(call: Call, response: Response) {
                val responseBody = response.body()?.string()
                val jsonResponse = JSONObject(responseBody)

                val result = jsonResponse.getString("result")

                callback(result == "success", null)
            }
        })
    }

    // 第四步：登录
    fun login(email: String, password: String, callback: (Boolean, String?, String?, String?) -> Unit) {
        val json = JSONObject().apply {
            put("email", email)
            put("password", password)
        }

        val requestBody = RequestBody.create(mediaType, json.toString())
        val request = Request.Builder()
            .url("$baseUrl/login")
            .post(requestBody)
            .build()

        client.newCall(request).enqueue(object : Callback {
            override fun onFailure(call: Call, e: IOException) {
                callback(false, null, null, e.message)
            }

            override fun onResponse(call: Call, response: Response) {
                val responseBody = response.body()?.string()
                val jsonResponse = JSONObject(responseBody)

                val result = jsonResponse.getString("result")

                if (result == "success") {
                    val data = jsonResponse.getJSONObject("data")
                    val accessToken = data.getString("access_token")
                    val refreshToken = data.getString("refresh_token")

                    callback(true, accessToken, refreshToken, null)
                } else {
                    callback(false, null, null, "登录失败")
                }
            }
        })
    }
}
```

### iOS (Swift) 实现

```swift
import Foundation

class DifyAuthService {
    private let baseURL = "https://您的服务器地址/console/api"

    // 第一步：发送注册验证码
    func sendRegisterCode(email: String, completion: @escaping (Bool, String?, String?) -> Void) {
        let url = URL(string: "\(baseURL)/register")!
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")

        let parameters: [String: Any] = [
            "email": email,
            "language": "zh-Hans"
        ]

        do {
            request.httpBody = try JSONSerialization.data(withJSONObject: parameters)
        } catch {
            completion(false, nil, "请求数据序列化失败")
            return
        }

        URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                completion(false, nil, error.localizedDescription)
                return
            }

            guard let data = data else {
                completion(false, nil, "没有返回数据")
                return
            }

            do {
                if let json = try JSONSerialization.jsonObject(with: data) as? [String: Any] {
                    let result = json["result"] as? String

                    if result == "success" {
                        let token = json["data"] as? String
                        completion(true, token, nil)
                    } else {
                        let message = json["message"] as? String ?? "注册失败"
                        completion(false, nil, message)
                    }
                } else {
                    completion(false, nil, "解析响应失败")
                }
            } catch {
                completion(false, nil, "解析响应失败: \(error.localizedDescription)")
            }
        }.resume()
    }

    // 第二步：验证注册验证码
    func verifyRegisterCode(email: String, code: String, token: String, completion: @escaping (Bool, String?) -> Void) {
        let url = URL(string: "\(baseURL)/register/validity")!
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")

        let parameters: [String: Any] = [
            "email": email,
            "code": code,
            "token": token
        ]

        do {
            request.httpBody = try JSONSerialization.data(withJSONObject: parameters)
        } catch {
            completion(false, "请求数据序列化失败")
            return
        }

        URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                completion(false, error.localizedDescription)
                return
            }

            guard let data = data else {
                completion(false, "没有返回数据")
                return
            }

            do {
                if let json = try JSONSerialization.jsonObject(with: data) as? [String: Any] {
                    let isValid = json["is_valid"] as? Bool ?? false
                    completion(isValid, nil)
                } else {
                    completion(false, "解析响应失败")
                }
            } catch {
                completion(false, "解析响应失败: \(error.localizedDescription)")
            }
        }.resume()
    }

    // 第三步：完成注册
    func completeRegistration(token: String, name: String, password: String, passwordConfirm: String, completion: @escaping (Bool, String?) -> Void) {
        let url = URL(string: "\(baseURL)/register/complete")!
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")

        let parameters: [String: Any] = [
            "token": token,
            "name": name,
            "password": password,
            "password_confirm": passwordConfirm
        ]

        do {
            request.httpBody = try JSONSerialization.data(withJSONObject: parameters)
        } catch {
            completion(false, "请求数据序列化失败")
            return
        }

        URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                completion(false, error.localizedDescription)
                return
            }

            guard let data = data else {
                completion(false, "没有返回数据")
                return
            }

            do {
                if let json = try JSONSerialization.jsonObject(with: data) as? [String: Any] {
                    let result = json["result"] as? String
                    completion(result == "success", nil)
                } else {
                    completion(false, "解析响应失败")
                }
            } catch {
                completion(false, "解析响应失败: \(error.localizedDescription)")
            }
        }.resume()
    }
}
```

## 注意事项

1. 密码必须符合系统要求：至少8个字符，包含字母和数字
2. 验证码有效期通常为10分钟，用户需要及时验证
3. 所有API请求应使用HTTPS以确保安全
4. 在生产环境中，应妥善处理各种错误情况
5. 移动端应用应保存登录后获取的访问令牌和刷新令牌，用于后续API调用

## 新旧API对比

| 功能 | 旧API | 新API |
|------|-------|-------|
| 发送验证码 | `/forgot-password` | `/register` |
| 验证验证码 | `/forgot-password/validity` | `/register/validity` |
| 完成注册 | `/reset-password` | `/register/complete` |

新API相比旧API的优势：
1. 更清晰的API命名，直接表明是注册功能
2. 不再依赖于重置密码流程，逻辑更加清晰
3. 支持设置用户姓名等更多信息
4. 不再检查`allowRegistration`设置，始终允许注册

## 常见问题

### Q: 验证码没有收到怎么办？
A: 请检查邮箱垃圾箱，或尝试重新发送验证码。

### Q: 密码设置有什么要求？
A: 密码必须至少8个字符，且同时包含字母和数字。

### Q: 注册后如何管理账号？
A: 注册成功后，可以登录系统进入个人设置页面管理账号信息。

### Q: 访问令牌有效期是多久？
A: 访问令牌通常有效期为60分钟，过期后需使用刷新令牌获取新的访问令牌。

### Q: 新旧API是否会同时支持？
A: 是的，为了保持兼容性，系统将同时支持新旧API，但建议新开发的应用使用新API。
