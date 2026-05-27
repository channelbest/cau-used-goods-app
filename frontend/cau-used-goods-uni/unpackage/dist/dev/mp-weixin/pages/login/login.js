"use strict";
const common_vendor = require("../../common/vendor.js");
const api_auth = require("../../api/auth.js");
const utils_auth = require("../../utils/auth.js");
const _sfc_main = {
  __name: "login",
  setup(__props) {
    const loading = common_vendor.ref(false);
    const goNext = (user) => {
      if ((user == null ? void 0 : user.authStatus) === "VERIFIED") {
        common_vendor.index.navigateTo({
          url: "/pages/home/home"
        });
        return;
      }
      common_vendor.index.navigateTo({
        url: "/pages/student-auth/student-auth"
      });
    };
    const handleDevLogin = async () => {
      if (loading.value)
        return;
      loading.value = true;
      try {
        const result = await api_auth.devLogin({
          openid: "frontend_a_dev_user",
          role: "USER"
        });
        utils_auth.saveLoginResult(result);
        common_vendor.index.showToast({
          title: "Login success",
          icon: "success"
        });
        goNext(result.user);
      } catch (error) {
        common_vendor.index.showToast({
          title: error.message || "Login failed",
          icon: "none"
        });
      } finally {
        loading.value = false;
      }
    };
    const previewStudentAuth = () => {
      common_vendor.index.navigateTo({
        url: "/pages/student-auth/student-auth"
      });
    };
    return (_ctx, _cache) => {
      return {
        a: loading.value,
        b: common_vendor.o(handleDevLogin, "3e"),
        c: common_vendor.o(previewStudentAuth, "6d")
      };
    };
  }
};
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["__scopeId", "data-v-e4e4508d"]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/login/login.js.map
