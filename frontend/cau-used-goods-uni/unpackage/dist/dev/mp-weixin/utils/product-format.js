"use strict";
const BASE_URL = "http://127.0.0.1:8080";

function normalizeImage(url) {
  if (!url) return "";
  if (url.indexOf("http://") === 0 || url.indexOf("https://") === 0) return url;
  if (url.indexOf("/uploads/") === 0) return BASE_URL + url;
  return url;
}

function formatPrice(price) {
  const n = Number(price || 0);
  return Number.isInteger(n) ? String(n) : n.toFixed(2);
}

function formatProduct(item, categoryMap) {
  const categoryName = categoryMap && categoryMap[item.categoryId] ? categoryMap[item.categoryId] : "未分类";
  const image = item.images && item.images.length ? normalizeImage(item.images[0]) : "";
  return {
    ...item,
    category: categoryName,
    priceText: formatPrice(item.price),
    conditionText: item.conditionLevel || "成色未填写",
    timeText: item.createTime || "",
    coverImage: image
  };
}

function buildCategoryMap(categories) {
  const map = {};
  (categories || []).forEach((item) => {
    map[item.id] = item.name;
  });
  return map;
}

exports.normalizeImage = normalizeImage;
exports.formatPrice = formatPrice;
exports.formatProduct = formatProduct;
exports.buildCategoryMap = buildCategoryMap;