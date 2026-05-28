"use strict";
const productApi = require("../../api/product.js");
const productFormat = require("../../utils/product-format.js");

Page({
  data: {
    loading: false,
    categories: [],
    goodsList: [],
    emptyText: ""
  },

  onShow() {
    this.loadHomeData();
  },

  async loadHomeData() {
    this.setData({ loading: true, emptyText: "" });
    try {
      const categories = await productApi.listCategories();
      const categoryMap = productFormat.buildCategoryMap(categories);
      const result = await productApi.listProducts({ status: "ON_SALE", sort: "newest", page: 1, pageSize: 10 });
      const list = (result.list || []).map(item => productFormat.formatProduct(item, categoryMap));
      this.setData({
        categories: (categories || []).slice(0, 6),
        goodsList: list,
        emptyText: list.length ? "" : "暂无在售商品"
      });
    } catch (error) {
      this.setData({ emptyText: error.message || "商品加载失败" });
      wx.showToast({ title: error.message || "商品加载失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
    }
  },

  goSearch() {
    wx.navigateTo({ url: "/pages/search/search" });
  },

  goMine() {
    wx.navigateTo({ url: "/pages/index/index" });
  },

  goCategory(event) {
    const id = event.currentTarget.dataset.id;
    const name = event.currentTarget.dataset.name;
    wx.navigateTo({ url: `/pages/category/category?categoryId=${id}&name=${encodeURIComponent(name)}` });
  },

  goDetail(event) {
    wx.navigateTo({ url: `/pages/detail/detail?id=${event.currentTarget.dataset.id}` });
  }
});