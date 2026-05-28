"use strict";
const api = require("../../api/auth.js");
const auth = require("../../utils/auth.js");

const BASE_URL = "http://127.0.0.1:8080";
const statusTextMap = {
  UNVERIFIED: "未认证",
  PENDING: "审核中",
  VERIFIED: "已认证",
  REJECTED: "未通过"
};

function normalizeAvatar(url) {
  if (!url) return "";
  if (url.indexOf("http://") === 0 || url.indexOf("https://") === 0) return url;
  if (url.indexOf("/uploads/") === 0) return BASE_URL + url;
  return url;
}

Page({
  data: {
    user: null,
    avatarUrl: "",
    authStatusText: "未认证",
    isAdmin: false
  },

  onShow() {
    const cachedUser = auth.getUser();
    if (cachedUser) this.applyUser(cachedUser);
    this.loadUser();
  },

  applyUser(user) {
    this.setData({
      user,
      avatarUrl: normalizeAvatar(user.avatarUrl),
      authStatusText: statusTextMap[user.authStatus] || "未认证",
      isAdmin: user.role === "ADMIN"
    });
  },

  async loadUser() {
    try {
      const user = await api.getCurrentUser();
      auth.setUser(user);
      this.applyUser(user);
    } catch (error) {
      const token = auth.getToken();
      if (!token) wx.reLaunch({ url: "/pages/login/login" });
    }
  },

  goProfileEdit() {
    wx.navigateTo({ url: "/pages/profile-edit/profile-edit" });
  },

  goStudentAuth() {
    wx.navigateTo({ url: "/pages/student-auth/student-auth" });
  },

  goStudentVerifyList() {
    wx.showToast({ title: "学生认证审核列表待后续页面接入", icon: "none" });
  },

  logout() {
    auth.clearAuth();
    wx.reLaunch({ url: "/pages/login/login" });
  }
});