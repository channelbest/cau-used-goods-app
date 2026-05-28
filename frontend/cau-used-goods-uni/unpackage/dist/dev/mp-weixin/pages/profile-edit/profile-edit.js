"use strict";
const api = require("../../api/auth.js");
const auth = require("../../utils/auth.js");

const BASE_URL = "http://127.0.0.1:8080";
function normalizeAvatar(url) {
  if (!url) return "";
  if (url.indexOf("http://") === 0 || url.indexOf("https://") === 0) return url;
  if (url.indexOf("/uploads/") === 0) return BASE_URL + url;
  return url;
}

Page({
  data: {
    nickname: "",
    phone: "",
    avatarUrl: "",
    loading: false,
    uploading: false
  },

  onLoad() {
    const user = auth.getUser() || {};
    this.setData({
      nickname: user.nickname || "",
      phone: user.phone || "",
      avatarUrl: normalizeAvatar(user.avatarUrl)
    });
  },

  onNicknameInput(event) { this.setData({ nickname: event.detail.value }); },
  onPhoneInput(event) { this.setData({ phone: event.detail.value }); },

  askUseWechatProfile() {
    wx.showModal({
      title: "使用微信资料",
      content: "是否沿用你的微信昵称和头像？",
      confirmText: "使用",
      cancelText: "不用",
      success: (res) => {
        if (res.confirm) this.useWechatProfile();
      }
    });
  },

  useWechatProfile() {
    if (!wx.getUserProfile) {
      wx.showToast({ title: "当前工具不支持获取微信资料", icon: "none" });
      return;
    }
    wx.getUserProfile({
      desc: "用于完善个人资料",
      success: async (res) => {
        const info = res.userInfo || {};
        this.setData({
          nickname: info.nickName || this.data.nickname,
          avatarUrl: info.avatarUrl || this.data.avatarUrl
        });
        try {
          const user = await api.updateProfile({
            nickname: info.nickName || this.data.nickname,
            avatarUrl: info.avatarUrl || ""
          });
          auth.setUser(user);
          this.setData({ avatarUrl: normalizeAvatar(user.avatarUrl) });
          wx.showToast({ title: "已使用微信资料", icon: "success" });
        } catch (error) {
          wx.showToast({ title: error.message || "保存微信资料失败", icon: "none" });
        }
      },
      fail: () => wx.showToast({ title: "已取消使用微信资料", icon: "none" })
    });
  },

  chooseAvatar() {
    wx.chooseMedia({
      count: 1,
      mediaType: ["image"],
      sourceType: ["album", "camera"],
      success: async (res) => {
        const filePath = res.tempFiles && res.tempFiles[0] && res.tempFiles[0].tempFilePath;
        if (!filePath) return;
        this.setData({ uploading: true });
        wx.showLoading({ title: "上传中" });
        try {
          const data = await api.uploadAvatar(filePath);
          if (data.user) {
            auth.setUser(data.user);
            this.setData({ avatarUrl: normalizeAvatar(data.user.avatarUrl) });
          }
          wx.showToast({ title: "头像已更新", icon: "success" });
        } catch (error) {
          wx.showToast({ title: error.message || "头像上传失败", icon: "none" });
        } finally {
          this.setData({ uploading: false });
          wx.hideLoading();
        }
      }
    });
  },

  async saveProfile() {
    if (this.data.loading) return;
    const nickname = this.data.nickname.trim();
    const phone = this.data.phone.trim();
    if (!nickname && !phone) {
      wx.showToast({ title: "请填写昵称或手机号", icon: "none" });
      return;
    }
    this.setData({ loading: true });
    try {
      const user = await api.updateProfile({ nickname, phone });
      auth.setUser(user);
      wx.showToast({ title: "资料已保存", icon: "success" });
      setTimeout(() => wx.navigateBack(), 600);
    } catch (error) {
      wx.showToast({ title: error.message || "保存失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
    }
  }
});