"use strict";
const auth = require("../../utils/auth.js");
const productApi = require("../../api/product.js");
const productFormat = require("../../utils/product-format.js");

Page({
  data: {
    productId: "",
    goods: null,
    favorited: false,
    loading: false,
    emptyText: "",
    showOrderForm: false,
    showReportForm: false,
    orderMeetTime: "",
    orderMeetLocation: "",
    orderRemark: "",
    reportReasonType: "商品信息异常",
    reportDescription: ""
  },

  onLoad(options) {
    this.setData({ productId: options.id || "" });
    this.loadDetail();
  },

  async loadDetail() {
    if (!this.data.productId) {
      this.setData({ emptyText: "商品不存在" });
      return;
    }
    this.setData({ loading: true, emptyText: "" });
    try {
      const categories = await productApi.listCategories();
      const categoryMap = productFormat.buildCategoryMap(categories);
      const product = await productApi.getProductById(this.data.productId);
      this.setData({ goods: productFormat.formatProduct(product, categoryMap) });
      this.loadFavoriteStatus();
    } catch (error) {
      this.setData({ emptyText: error.message || "商品加载失败" });
      wx.showToast({ title: error.message || "商品加载失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
    }
  },

  async loadFavoriteStatus() {
    if (!auth.getToken()) return;
    try {
      const result = await productApi.checkFavorite(Number(this.data.productId));
      this.setData({ favorited: !!result.favorited });
    } catch (error) {}
  },

  checkVerified() {
    const user = auth.getUser();
    if (!auth.getToken()) {
      wx.showToast({ title: "请先登录", icon: "none" });
      return false;
    }
    if (!user || user.authStatus !== "VERIFIED") {
      wx.showToast({ title: "请先在个人中心完成学生认证", icon: "none" });
      return false;
    }
    return true;
  },

  async toggleFavorite() {
    if (!this.checkVerified()) return;
    try {
      if (this.data.favorited) {
        await productApi.removeFavorite(Number(this.data.productId));
        this.setData({ favorited: false });
        wx.showToast({ title: "已取消收藏", icon: "success" });
      } else {
        await productApi.addFavorite(Number(this.data.productId));
        this.setData({ favorited: true });
        wx.showToast({ title: "已收藏", icon: "success" });
      }
    } catch (error) {
      wx.showToast({ title: error.message || "操作失败", icon: "none" });
    }
  },

  openOrderForm() {
    if (!this.checkVerified()) return;
    this.setData({ showOrderForm: true });
  },

  closeOrderForm() {
    this.setData({ showOrderForm: false });
  },

  onOrderMeetTimeInput(event) { this.setData({ orderMeetTime: event.detail.value }); },
  onOrderMeetLocationInput(event) { this.setData({ orderMeetLocation: event.detail.value }); },
  onOrderRemarkInput(event) { this.setData({ orderRemark: event.detail.value }); },

  async submitOrder() {
    try {
      const payload = { productId: Number(this.data.productId) };
      if (this.data.orderMeetTime.trim()) payload.meetTime = this.data.orderMeetTime.trim();
      if (this.data.orderMeetLocation.trim()) payload.meetLocation = this.data.orderMeetLocation.trim();
      if (this.data.orderRemark.trim()) payload.remark = this.data.orderRemark.trim();
      await productApi.createOrder(payload);
      this.setData({ showOrderForm: false, orderMeetTime: "", orderMeetLocation: "", orderRemark: "" });
      wx.showToast({ title: "预约已提交", icon: "success" });
    } catch (error) {
      wx.showToast({ title: error.message || "预约失败", icon: "none" });
    }
  },

  openReportForm() {
    if (!this.checkVerified()) return;
    this.setData({ showReportForm: true });
  },

  closeReportForm() {
    this.setData({ showReportForm: false });
  },

  onReportReasonInput(event) { this.setData({ reportReasonType: event.detail.value }); },
  onReportDescriptionInput(event) { this.setData({ reportDescription: event.detail.value }); },

  async submitReport() {
    if (!this.data.reportReasonType.trim()) {
      wx.showToast({ title: "请填写举报原因", icon: "none" });
      return;
    }
    try {
      const description = this.data.reportDescription.trim();
      await productApi.createReport({
        targetType: "PRODUCT",
        targetId: Number(this.data.productId),
        reasonType: this.data.reportReasonType.trim(),
        description: description || "用户在商品详情页提交举报",
        images: []
      });
      this.setData({ showReportForm: false, reportReasonType: "商品信息异常", reportDescription: "" });
      wx.showToast({ title: "举报已提交", icon: "success" });
    } catch (error) {
      wx.showToast({ title: error.message || "举报失败", icon: "none" });
    }
  }
});