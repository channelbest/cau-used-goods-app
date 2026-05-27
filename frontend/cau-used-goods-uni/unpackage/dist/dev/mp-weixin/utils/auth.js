"use strict";
const common_vendor = require("../common/vendor.js");
const TOKEN_KEY = "CAU_USED_GOODS_TOKEN";
const USER_KEY = "CAU_USED_GOODS_USER";
const getToken = () => {
  return common_vendor.index.getStorageSync(TOKEN_KEY) || "";
};
const setToken = (token) => {
  common_vendor.index.setStorageSync(TOKEN_KEY, token);
};
const setUser = (user) => {
  common_vendor.index.setStorageSync(USER_KEY, user);
};
const clearAuth = () => {
  common_vendor.index.removeStorageSync(TOKEN_KEY);
  common_vendor.index.removeStorageSync(USER_KEY);
};
const saveLoginResult = (result) => {
  if (result == null ? void 0 : result.token) {
    setToken(result.token);
  }
  if (result == null ? void 0 : result.user) {
    setUser(result.user);
  }
};
exports.clearAuth = clearAuth;
exports.getToken = getToken;
exports.saveLoginResult = saveLoginResult;
//# sourceMappingURL=../../.sourcemap/mp-weixin/utils/auth.js.map
