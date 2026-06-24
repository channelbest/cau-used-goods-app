import { request } from '../utils/request'

export const getPublicProfile = (id) => request({
  url: `/users/${id}/public`
})

export const getPublicHomepage = (id, params = {}) => request({
  url: `/users/${id}/homepage`,
  data: params
})

export const getSellerReviews = (id, params = {}) => request({
  url: `/sellers/${id}/reviews`,
  data: params
})
