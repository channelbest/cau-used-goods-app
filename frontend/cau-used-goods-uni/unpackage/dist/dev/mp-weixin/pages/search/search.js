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
    keyword: "",
    categoryId: "",
    condition: "all",
    minPrice: "",
    maxPrice: "",
    sortType: "latest",
    categories: [],
    conditions: ["全部", "九成新", "八成新", "七成新"],
    goodsList: [],
    emptyText: ""
  },

  onLoad() {
    this.loadCategoriesAndProducts();
  },

  async loadCategoriesAndProducts() {
    try {
      const categories = await productApi.listCategories();
      this.setData({ categories: [{ id: "", name: "全部" }].concat(categories || []) });
      this.refreshList();
    } catch (error) {
      wx.showToast({ title: error.message || "分类加载失败", icon: "none" });
    }
  },

  onKeywordInput(event) { this.setData({ keyword: event.detail.value || "" }); this.refreshList(); },
  onMinPriceInput(event) { this.setData({ minPrice: (event.detail.value || "").replace(/\D/g, "") }); this.refreshList(); },
  onMaxPriceInput(event) { this.setData({ maxPrice: (event.detail.value || "").replace(/\D/g, "") }); this.refreshList(); },
  setCategory(event) { this.setData({ categoryId: event.currentTarget.dataset.id }); this.refreshList(); },
  setCondition(event) { this.setData({ condition: event.currentTarget.dataset.condition }); this.refreshList(); },
  setSort(event) { this.setData({ sortType: event.currentTarget.dataset.sort }); this.refreshList(); },

  resetFilters() {
    this.setData({ keyword: "", categoryId: "", condition: "all", minPrice: "", maxPrice: "", sortType: "latest" });
    this.refreshList();
  },

  async refreshList() {
    try {
      const categoryMap = productFormat.buildCategoryMap(this.data.categories.filter(item => item.id));
      const result = await productApi.listProducts({
        keyword: this.data.keyword.trim(),
        categoryId: this.data.categoryId,
        status: "ON_SALE",
        minPrice: this.data.minPrice,
        maxPrice: this.data.maxPrice,
        sort: sortMap[this.data.sortType] || "newest",
        page: 1,
        pageSize: 50
      });
      let list = (result.list || []).map(item => productFormat.formatProduct(item, categoryMap));
      if (this.data.condition !== "all" && this.data.condition !== "全部") {
        list = list.filter(item => item.conditionLevel === this.data.condition);
      }
      this.setData({ goodsList: list, emptyText: list.length ? "" : "暂无符合条件的商品" });
    } catch (error) {
      this.setData({ emptyText: error.message || "搜索失败" });
    }
  },

  goDetail(event) {
    wx.navigateTo({ url: `/pages/detail/detail?id=${event.currentTarget.dataset.id}` });
  }
});