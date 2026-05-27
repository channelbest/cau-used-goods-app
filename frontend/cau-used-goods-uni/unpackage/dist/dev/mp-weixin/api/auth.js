"use strict";
const utils_request = require("../utils/request.js");
const devLogin = (payload = {}) => {
  return utils_request.request({
    url: "/auth/dev-login",
    method: "POST",
    data: payload,
    auth: false
  });
};
const submitStudentVerification = (payload) => {
  return utils_request.request({
    url: "/users/student-verify",
    method: "POST",
    data: payload
  });
};
exports.devLogin = devLogin;
exports.submitStudentVerification = submitStudentVerification;
//# sourceMappingURL=../../.sourcemap/mp-weixin/api/auth.js.map
