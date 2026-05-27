"use strict";
const common_vendor = require("../../common/vendor.js");
const _sfc_main = {
  __name: "category",
  setup(__props) {
    const goodsList = [
      {
        id: 1,
        title: "教材书籍示例商品",
        condition: "八成新",
        price: 20
      },
      {
        id: 2,
        title: "分类下的商品",
        condition: "九成新",
        price: 35
      }
    ];
    const goDetail = (id) => {
      common_vendor.index.navigateTo({
        url: `/pages/detail/detail?id=${id}`
      });
    };
    return (_ctx, _cache) => {
      return {
        a: common_vendor.f(goodsList, (item, k0, i0) => {
          return {
            a: common_vendor.t(item.title),
            b: common_vendor.t(item.condition),
            c: common_vendor.t(item.price),
            d: item.id,
            e: common_vendor.o(($event) => goDetail(item.id), item.id)
          };
        })
      };
    };
  }
};
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["__scopeId", "data-v-8145b772"]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/category/category.js.map
