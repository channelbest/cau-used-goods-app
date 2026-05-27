"use strict";
const common_vendor = require("../../common/vendor.js");
const _sfc_main = {
  __name: "home",
  setup(__props) {
    const goodsList = [
      {
        id: 1,
        title: "高等数学教材",
        condition: "八成新",
        price: 20
      },
      {
        id: 2,
        title: "蓝牙耳机",
        condition: "九成新",
        price: 68
      },
      {
        id: 3,
        title: "台灯",
        condition: "七成新",
        price: 15
      }
    ];
    const goSearch = () => {
      common_vendor.index.navigateTo({
        url: "/pages/search/search"
      });
    };
    const goCategory = (category) => {
      common_vendor.index.navigateTo({
        url: `/pages/category/category?category=${category}`
      });
    };
    const goDetail = (id) => {
      common_vendor.index.navigateTo({
        url: `/pages/detail/detail?id=${id}`
      });
    };
    return (_ctx, _cache) => {
      return {
        a: common_vendor.o(goSearch, "a0"),
        b: common_vendor.o(($event) => goCategory("book"), "ed"),
        c: common_vendor.o(($event) => goCategory("digital"), "f7"),
        d: common_vendor.o(($event) => goCategory("daily"), "dd"),
        e: common_vendor.o(($event) => goCategory("sport"), "08"),
        f: common_vendor.f(goodsList, (item, k0, i0) => {
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
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["__scopeId", "data-v-07e72d3c"]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/home/home.js.map
