"use strict";
const utils_request = require("../utils/request.js");
const utils_auth = require("../utils/auth.js");
const BASE_URL = "http://127.0.0.1:8080";

const devLogin = (payload = {}) => {
  return utils_request.request({
    url: "/auth/dev-login",
    method: "POST",
    data: payload,
    auth: false
  });
};

const wechatLogin = (code) => {
  return utils_request.request({
    url: "/auth/wechat-login",
    method: "POST",
    data: { code },
    auth: false
  });
};

const getCurrentUser = () => {
  return utils_request.request({
    url: "/users/me"
  });
};

const updateProfile = (payload) => {
  return utils_request.request({
    url: "/users/profile",
    method: "PUT",
    data: payload
  });
};

const submitStudentVerification = (payload) => {
  return utils_request.request({
    url: "/users/student-verify",
    method: "POST",
    data: payload
  });
};

const uploadAvatar = (filePath) => {
  const token = utils_auth.getToken();
  return new Promise((resolve, reject) => {
    wx.uploadFile({
      url: `${BASE_URL}/users/avatar`,
      filePath,
      name: "avatar",
      header: {
        Authorization: `Bearer ${token}`
      },
      success: (res) => {
        let body = {};
        try {
          body = JSON.parse(res.data || "{}");
        } catch (error) {
          reject(new Error("头像上传响应解析失败"));
          return;
        }
        if (res.statusCode < 200 || res.statusCode >= 300 || body.code !== 0) {
          reject(new Error(body.message || "头像上传失败"));
          return;
        }
        resolve(body.data);
      },
      fail: () => reject(new Error("无法连接服务器，请确认后端服务已启动"))
    });
  });
};

exports.devLogin = devLogin;
exports.wechatLogin = wechatLogin;
exports.getCurrentUser = getCurrentUser;
exports.updateProfile = updateProfile;
exports.submitStudentVerification = submitStudentVerification;
exports.uploadAvatar = uploadAvatar;
