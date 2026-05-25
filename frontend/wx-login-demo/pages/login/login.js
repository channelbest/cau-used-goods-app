const API_BASE = 'http://127.0.0.1:8080'

Page({
  data: {
    code: '',
    token: '',
    message: ''
  },

  getCode() {
    wx.login({
      success: (res) => {
        if (!res.code) {
          this.setData({ message: 'wx.login 成功但没有返回 code' })
          return
        }

        console.log('wx.login code:', res.code)
        this.setData({
          code: res.code,
          token: '',
          message: '已获取 code，可复制到 Postman，或点击请求后端。'
        })
      },
      fail: (err) => {
        console.error('wx.login failed:', err)
        this.setData({ message: JSON.stringify(err, null, 2) })
      }
    })
  },

  loginBackend() {
    if (!this.data.code) {
      this.setData({ message: '请先获取 code' })
      return
    }

    wx.request({
      url: `${API_BASE}/auth/wechat-login`,
      method: 'POST',
      header: {
        'content-type': 'application/json'
      },
      data: {
        code: this.data.code
      },
      success: (res) => {
        console.log('/auth/wechat-login response:', res)
        const token = res.data && res.data.data ? res.data.data.token : ''
        this.setData({
          token,
          message: JSON.stringify(res.data, null, 2)
        })
      },
      fail: (err) => {
        console.error('/auth/wechat-login failed:', err)
        this.setData({ message: JSON.stringify(err, null, 2) })
      }
    })
  }
})
