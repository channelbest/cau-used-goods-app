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
          reject(new Error(body.message || "Login expired, please login again"));
          return;
        }
        if (res.statusCode < 200 || res.statusCode >= 300 || body.code !== 0) {
          reject(new Error(body.message || "Request failed"));
          return;
        }
        resolve(body.data);
      },
      fail: () => {
        reject(new Error("Cannot connect to server. Please make sure backend is running."));
      }
    });
  });
};
exports.request = request;
//# sourceMappingURL=../../.sourcemap/mp-weixin/utils/request.js.map
