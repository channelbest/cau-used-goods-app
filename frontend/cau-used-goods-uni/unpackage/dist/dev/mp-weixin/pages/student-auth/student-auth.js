"use strict";
const api = require("../../api/auth.js");
const auth = require("../../utils/auth.js");

Page({
  data: {
    realName: "",
    studentId: "",
    college: "",
    loading: false
  },

  onRealNameInput(event) {
    const value = (event.detail.value || "").replace(/[^\u4e00-\u9fa5]/g, "");
    this.setData({ realName: value });
  },

  onStudentIdInput(event) {
    const value = (event.detail.value || "").replace(/\D/g, "");
    this.setData({ studentId: value });
  },

  onCollegeInput(event) {
    const value = (event.detail.value || "").replace(/[^\u4e00-\u9fa5]/g, "");
    this.setData({ college: value });
  },

  validate() {
    const chineseReg = /^[\u4e00-\u9fa5]+$/;
    const numberReg = /^\d+$/;
    if (!this.data.realName.trim()) return "请填写姓名";
    if (!chineseReg.test(this.data.realName.trim())) return "姓名只能填写汉字";
    if (!this.data.studentId.trim()) return "请填写学号";
    if (!numberReg.test(this.data.studentId.trim())) return "学号只能填写数字";
    if (!this.data.college.trim()) return "请填写学院";
    if (!chineseReg.test(this.data.college.trim())) return "学院只能填写汉字";
    return "";
  },

  async refreshUserAndBack(title) {
    try {
      const user = await api.getCurrentUser();
      auth.setUser(user);
    } catch (error) {}
    wx.showToast({ title, icon: "success" });
    setTimeout(() => wx.navigateBack(), 600);
  },

  async submitAuth() {
    const message = this.validate();
    if (message) {
      wx.showToast({ title: message, icon: "none" });
      return;
    }
    if (this.data.loading) return;
    this.setData({ loading: true });
    try {
      await api.submitStudentVerification({
        realName: this.data.realName.trim(),
        studentId: this.data.studentId.trim(),
        college: this.data.college.trim()
      });
      await this.refreshUserAndBack("认证已提交");
    } catch (error) {
      const msg = error.message || "提交失败";
      if (msg.includes("already approved") || msg.includes("已经认证") || msg.includes("已认证")) {
        await this.refreshUserAndBack("认证已通过");
        return;
      }
      if (msg.includes("pending")) {
        await this.refreshUserAndBack("认证审核中");
        return;
      }
      wx.showToast({ title: msg, icon: "none" });
    } finally {
      this.setData({ loading: false });
    }
  }
});