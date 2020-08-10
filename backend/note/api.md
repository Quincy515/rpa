1. 复制 `config/config.example.toml` 新建 `config/config.toml` 文件 配置信息
2. mysql 数据库配置 `note/db.sql`  
`mysql – u 用户名 – p 密码 < sql脚本` 
3. swagger 文档生产 在 `backend` 目录下执行
`swag init -g http/controller/router.go`

### 1. 邮件注册
POST http://localhost:7788/api/v1/email-register

`{"email":"2@1.com","password":"Aa123456", "confirm_password":"Aa123456"}`

#### 返回
返回的session过期时间是30天，注册的生成的token是默认7天过期

```text
POST http://localhost:7788/api/v1/email-register

HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Set-Cookie: mysession=MTU5MzcxODc2OXxEdi1CQkFFQ180SUFBUkFCRUFBQV84SF9nZ0FCQm5OMGNtbHVad3dIQUFWMGIydGxiZ1p6ZEhKcGJtY01fNk1BXzZCbGVVcG9Za2RqYVU5cFNrbFZla2t4VG1sSmMwbHVValZqUTBrMlNXdHdXRlpEU2prdVpYbEtNV015Vm5sVk1qUnBUMnBKTWs1cVVURk9WRWw2VDFSbk1VMTZUVEJOYW1ONVRFTktiR1ZJUVdsUGFrVXhUMVJSZWsxcVRURk9hbXR6U1cxc2VtTjVTVFpKYm14d1kyMUdkVWx1TUM1WWNXcEhha296U0U5WlQwSTNhMFZ4YjFSdGJqQXlRMFJsWmxOS1dEUkxSVmxGU0ZCaWFHTm1TVVpyfJn0uqq8nB8stVdQfjbYsSGlHh1XW-oVOKzJ_WYjFM04; Path=/; Expires=Sat, 04 Aug 1787 00:53:39 GMT; Max-Age=2592000000000000
Date: Thu, 02 Jul 2020 19:39:29 GMT
Content-Length: 247
```

```json
{
  "code": 1000,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyU24iOjI2NjQ1NTIzOTg1MzM0MjcyLCJleHAiOjE1OTQzMjM1NjksImlzcyI6InlpcmFuIn0.XqjGjJ3HOYOB7kEqoTmn02CDefSJX4KEYEHPbhcfIFk",
    "userSn": 26645523985334272
  }
}
```


### 2. 邮件登录
POST http://localhost:7788/api/v1/email-login

`{"email":"1@1.com","password":"Aa123456"}`

#### 返回
返回的session过期时间是30天，登录(选择记住我)生成的token过期时间也是30天

```text
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Set-Cookie: mysession=MTU5MzcxODUwNXxEdi1CQkFFQ180SUFBUkFCRUFBQV84SF9nZ0FCQm5OMGNtbHVad3dIQUFWMGIydGxiZ1p6ZEhKcGJtY01fNk1BXzZCbGVVcG9Za2RqYVU5cFNrbFZla2t4VG1sSmMwbHVValZqUTBrMlNXdHdXRlpEU2prdVpYbEtNV015Vm5sVk1qUnBUMnBKTWs1cVRYbE9SRkV6VG1wbk5FOVVXVEpOVkZFd1RFTktiR1ZJUVdsUGFrVXhUMVJOTTAxVVp6Rk5SRlZ6U1cxc2VtTjVTVFpKYm14d1kyMUdkVWx1TUM1YVRuaHpUVWQzUjBFdGRuQlJSRkpIYjB4dWFXb3RUR0oyTFVvellUTlZhMjAzTTJJM1ozQm1kR2c0fCcpOTVt02ht8aYHRppkzna3yZZUJhLMIfA0GelBKBZ6; Path=/; Expires=Sat, 04 Aug 1787 00:49:14 GMT; Max-Age=2592000000000000
Date: Thu, 02 Jul 2020 19:35:05 GMT
Content-Length: 247
```

```json
{"code":1000,"message":"success","data":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyU24iOjI2NjMyNDQ3Njg4OTY2MTQ0LCJleHAiOjE1OTM3MTg1MDUsImlzcyI6InlpcmFuIn0.ZNxsMGwGA-vpQDRGoLnij-Lbv-J3a3Ukm73b7gpfth8","userSn":26632447688966144}}
```


### 3. 上传文件
POST http://localhost:7788/api/v1/upload

Content-Type: multipart/form-data; boundary=WebAppBoundary

Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyU24iOjI0MDc4OTY5MTM3NTI4ODMyLCJleHAiOjE1OTIzNzUwMzgsImlzcyI6IllhS3RvQui_kOiQpeekvuWMuiJ9.EQVMzoArpPbLRO7U0LTHG5s83mC9-1a_E2bMRbNLgi8

--WebAppBoundary

Content-Disposition: form-data; name="file"; filename="custer1.png"

< /Users/tianxiaoqiang/Pictures/1.jpg

--WebAppBoundary--