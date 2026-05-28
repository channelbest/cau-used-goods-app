"use strict";
const TOKEN_KEY = "CAU_USED_GOODS_TOKEN";
const USER_KEY = "CAU_USED_GOODS_USER";
const getToken = () => wx.getStorageSync(TOKEN_KEY) || "";
const setToken = (token) => wx.setStorageSync(TOKEN_KEY, token);
const getUser = () => wx.getStorageSync(USER_KEY) || null;
const setUser = (user) => wx.setStorageSync(USER_KEY, user);
const clearAuth = () => {
  wx.removeStorageSync(TOKEN_KEY);
  wx.removeStorageSync(USER_KEY);
};
const saveLoginResult = (result) => {
  if (result && result.token) setToken(result.token);
  if (result && result.user) setUser(result.user);
};
exports.clearAuth = clearAuth;
exports.getToken = getToken;
exports.getUser = getUser;
exports.saveLoginResult = saveLoginResult;
exports.setUser = setUser;
