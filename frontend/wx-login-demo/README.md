# 微信登录最小测试小程序

这个目录是后端登录联调用的最小微信小程序，不是正式前端。

## 使用方式

1. 打开微信开发者工具。
2. 选择“导入项目”。
3. 项目目录选择：

```text
frontend/wx-login-demo
```

4. AppID 填你们的小程序 AppID。
5. 打开后点击页面里的“获取 wx.login code”。
6. 页面会显示 code，也会打印到控制台。

## 直接请求后端

如果后端已启动：

```powershell
cd backend
go run .\cmd\server
```

默认请求地址是：

```text
http://127.0.0.1:8080/auth/wechat-login
```

在微信开发者工具里需要勾选：

```text
详情 -> 本地设置 -> 不校验合法域名、web-view、TLS 版本以及 HTTPS 证书
```

否则本地 `http://127.0.0.1:8080` 请求可能会被拦截。

