"use strict";
const productApi = require("../../api/product.js");
const productFormat = require("../../utils/product-format.js");

const sortMap = {
  latest: "newest",
  priceAsc: "price_asc",
  priceDesc: "price_desc"
};

Page({
  data: {
    categoryId: "",
    title: "分类商品",
    sortType: "latest",
    goodsList: [],
    emptyText: "",
    loading: false
  },

  onLoad(options) {
    this.setData({
      categoryId: options.categoryId || "",
      title: options.name ? decodeURIComponent(options.name) : "分类商品"
    });
    this.loadProducts();
  },

  setSort(event) {
    this.setData({ sortType: event.currentTarget.dataset.sort });
    this.loadProducts();
  },

  async loadProducts() {
    this.setData({ loading: true, emptyText: "" });
    try {
      const categories = await productApi.listCategories();
      const categoryMap = productFormat.buildCategoryMap(categories);
      const result = await productApi.listProducts({
        categoryId: this.data.categoryId,
        status: "ON_SALE",
        sort: sortMap[this.data.sortType] || "newest",
        page: 1,
        pageSize: 50
      });
      const list = (result.list || []).map(item => productFormat.formatProduct(item, categoryMap));
      this.setData({ goodsList: list, emptyText: list.length ? "" : "该分类暂无商品" });
    } catch (error) {
      this.setData({ emptyText: error.message || "商品加载失败" });
      wx.showToast({ title: error.message || "商品加载失败", icon: "none" });
    } finally {
      this.setData({ loading: false });
    }
  },

  goDetail(event) {
    wx.navigateTo({ url: `/pages/detail/detail?id=${event.currentTarget.dataset.id}` });
  }
});