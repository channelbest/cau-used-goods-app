"use strict";
const requestUtil = require("../utils/request.js");

function buildQuery(params) {
  const parts = [];
  Object.keys(params || {}).forEach((key) => {
    const value = params[key];
    if (value === undefined || value === null || value === "") return;
    parts.push(`${encodeURIComponent(key)}=${encodeURIComponent(value)}`);
  });
  return parts.length ? `?${parts.join("&")}` : "";
}

const listCategories = () => requestUtil.request({
  url: "/categories",
  auth: false
});

const listProducts = (params = {}) => requestUtil.request({
  url: `/products${buildQuery(params)}`,
  auth: false
});

const getProductById = (id) => requestUtil.request({
  url: `/products/${id}`,
  auth: false
});

const addFavorite = (productId) => requestUtil.request({
  url: "/favorites",
  method: "POST",
  data: { productId }
});

const removeFavorite = (productId) => requestUtil.request({
  url: `/favorites/${productId}`,
  method: "DELETE"
});

const checkFavorite = (productId) => requestUtil.request({
  url: `/favorites/check${buildQuery({ productId })}`
});

const createOrder = (payload) => requestUtil.request({
  url: "/orders",
  method: "POST",
  data: payload
});

const createReport = (payload) => requestUtil.request({
  url: "/reports",
  method: "POST",
  data: payload
});

exports.listCategories = listCategories;
exports.listProducts = listProducts;
exports.getProductById = getProductById;
exports.addFavorite = addFavorite;
exports.removeFavorite = removeFavorite;
exports.checkFavorite = checkFavorite;
exports.createOrder = createOrder;
exports.createReport = createReport;