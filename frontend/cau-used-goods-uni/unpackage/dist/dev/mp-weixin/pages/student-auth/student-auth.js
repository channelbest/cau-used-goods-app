"use strict";
const common_vendor = require("../../common/vendor.js");
const api_auth = require("../../api/auth.js");
const _sfc_main = {
  __name: "student-auth",
  setup(__props) {
    const loading = common_vendor.ref(false);
    const form = common_vendor.reactive({
      realName: "",
      studentId: "",
      college: ""
    });
    const validateForm = () => {
      if (!form.realName.trim())
        return "Enter real name";
      if (!form.studentId.trim())
        return "Enter student ID";
      if (!form.college.trim())
        return "Enter college";
      return "";
    };
    const goHome = () => {
      common_vendor.index.navigateTo({
        url: "/pages/home/home"
      });
    };
    const handleSubmit = async () => {
      const message = validateForm();
      if (message) {
        common_vendor.index.showToast({
          title: message,
          icon: "none"
        });
        return;
      }
      if (loading.value)
        return;
      loading.value = true;
      try {
        await api_auth.submitStudentVerification({
          studentId: form.studentId.trim(),
          realName: form.realName.trim(),
          college: form.college.trim()
        });
        common_vendor.index.showToast({
          title: "Submit success",
          icon: "success"
        });
        goHome();
      } catch (error) {
        common_vendor.index.showToast({
          title: error.message || "Submit failed",
          icon: "none"
        });
      } finally {
        loading.value = false;
      }
    };
    return (_ctx, _cache) => {
      return {
        a: form.realName,
        b: common_vendor.o(($event) => form.realName = $event.detail.value, "2f"),
        c: form.studentId,
        d: common_vendor.o(($event) => form.studentId = $event.detail.value, "c3"),
        e: form.college,
        f: common_vendor.o(($event) => form.college = $event.detail.value, "0d"),
        g: loading.value,
        h: common_vendor.o(handleSubmit, "b0"),
        i: common_vendor.o(goHome, "53")
      };
    };
  }
};
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["__scopeId", "data-v-34b2aeae"]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/student-auth/student-auth.js.map
