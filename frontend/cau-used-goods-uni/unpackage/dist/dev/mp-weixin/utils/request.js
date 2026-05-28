"use strict";
const common_vendor = require("../common/vendor.js");
const utils_auth = require("./auth.js");
const BASE_URL = "http://127.0.0.1:8080";
const request = ({
  url,
  method = "GET",
  data = {},
  header = {},
  auth = true
}) => {
  const token = utils_auth.getToken();
  const requestHeader = {
    "Content-Type": "application/json",
    ...header
  };
  if (auth && token) {
    requestHeader.Authorization = `Bearer ${token}`;
  }
  return new Promise((resolve, reject) => {
    common_vendor.index.request({
      url: `${BASE_URL}${url}`,
      method,
      data,
      header: requestHeader,
      success: (res) => {
        const body = res.data || {};
        if (res.statusCode === 401) {
          utils_auth.clearAuth();
          reject(new Error(body.message || "登录已过期，请重新登录"));
          return;
        }
        if (res.statusCode < 200 || res.statusCode >= 300 || body.code !== 0) {
          reject(new Error(body.message || "请求失败"));
          return;
        }
        resolve(body.data);
      },
      fail: () => {
        reject(new Error("无法连接服务器，请确认后端服务已启动"));
      }
    });
  });
};
exports.request = request;
