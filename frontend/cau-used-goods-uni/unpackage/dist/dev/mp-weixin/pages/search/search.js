"use strict";
const common_vendor = require("../../common/vendor.js");
const _sfc_main = {
  __name: "search",
  setup(__props) {
    const keyword = common_vendor.ref("");
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
      }
    ];
    const goDetail = (id) => {
      common_vendor.index.navigateTo({
        url: `/pages/detail/detail?id=${id}`
      });
    };
    return (_ctx, _cache) => {
      return {
        a: keyword.value,
        b: common_vendor.o(($event) => keyword.value = $event.detail.value, "88"),
        c: common_vendor.f(goodsList, (item, k0, i0) => {
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
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["__scopeId", "data-v-c10c040c"]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/search/search.js.map
