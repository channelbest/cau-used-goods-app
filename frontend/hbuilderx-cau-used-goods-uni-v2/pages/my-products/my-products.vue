<template>
  <view class="page">
    <view class="summary">
      <view>
        <view class="title">我发布的</view>
        <view class="muted">管理自己的闲置商品</view>
      </view>
      <view class="summary-actions">
        <button v-if="canBatchPutOnSale" class="batch-on-sale" @click="batchPutOnSale">一键上架</button>
        <button class="publish" @click="goPublish">发布</button>
      </view>
    </view>

    <view v-if="products.length" class="list">
      <view v-for="item in products" :key="item.id" class="product">
        <view class="product-main" @click="openDetail(item.id)">
          <image v-if="item.coverImage" class="cover" :src="item.coverImage" mode="aspectFill" />
          <view v-else class="cover placeholder">暂无图片</view>
          <view class="content">
            <view class="top">
              <view class="name">{{ item.title }}</view>
              <text class="status" :class="productStatus(item)">{{ getStatusText(productStatus(item)) }}</text>
            </view>
            <view class="meta">{{ item.category }} · {{ item.conditionText }}</view>
            <view v-if="offShelfReason(item)" class="off-shelf-reason">下架原因：{{ offShelfReason(item) }}</view>
            <text class="price">￥{{ item.priceText }}</text>
          </view>
        </view>
        <view v-if="canEdit(item)" class="actions">
          <button v-if="productStatus(item) === 'ON_SALE'" class="action muted-action" @click="changeStatus(item, 'OFF_SHELF')">下架</button>
          <button v-if="canPutOnSale(item)" class="action primary-action" @click="changeStatus(item, 'ON_SALE')">上架</button>
          <text v-else-if="offShelfTip(item)" class="off-shelf-tip">{{ offShelfTip(item) }}</text>
          <button v-if="canEdit(item)" class="action" @click="editProduct(item)">编辑</button>
          <button v-if="canEdit(item)" class="action danger-action" @click="deleteMyProduct(item)">删除</button>
        </view>
      </view>
    </view>

    <view v-else-if="!loading" class="empty">
      <view class="empty-title">还没有发布商品</view>
      <view class="muted">把闲置教材、数码和生活用品发布出去吧</view>
      <button class="empty-button" @click="goPublish">去发布</button>
    </view>
    <view class="load-state">{{ loading ? '正在加载...' : '' }}</view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onPullDownRefresh, onShow } from '@dcloudio/uni-app'
import { batchPutOnSaleProducts, deleteProduct, listCategories, listMyProducts, updateProductStatus } from '../../api/product'
import { getCurrentUser } from '../../api/auth'
import { getUser, setUser } from '../../utils/auth'
import { buildCategoryMap, formatProduct, getStatusText } from '../../utils/product-format'
import { accountStatusOf, userTradeRestrictionMessage } from '../../utils/user-format'

const categories = ref([])
const rawProducts = ref([])
const loading = ref(false)
const products = computed(() => rawProducts.value.map((item) => formatProduct(item, buildCategoryMap(categories.value))))
const restorableOffShelfSources = ['USER', 'ACCOUNT_STATUS']
const canBatchPutOnSale = computed(() => products.value.some((item) => canPutOnSale(item)))

const loadData = async () => {
  loading.value = true
  try {
    const [categoryList, productList] = await Promise.all([listCategories(), listMyProducts()])
    categories.value = categoryList || []
    rawProducts.value = productList || []
  } catch (error) {
    uni.showToast({ title: error.message, icon: 'none' })
  } finally {
    loading.value = false
    uni.stopPullDownRefresh()
  }
}

const openDetail = (id) => uni.navigateTo({ url: `/pages/detail/detail?id=${id}` })
const goPublish = () => {
  uni.removeStorageSync('PUBLISH_EDIT_PRODUCT_ID')
  uni.removeStorageSync('PUBLISH_EDIT_PRODUCT_DATA')
  uni.removeStorageSync('PUBLISH_EDIT_INTENT')
  uni.switchTab({ url: '/pages/publish/publish' })
}
const normalizeCode = (value) => String(value || '').trim().toUpperCase()
const productStatus = (item) => normalizeCode(item.status || item.productStatus || item.product_status)
const canEdit = (item) => ['ON_SALE', 'OFF_SHELF'].includes(productStatus(item))
const offShelfBy = (item) => normalizeCode(item.offShelfBy || item.off_shelf_by)
const offShelfReason = (item) => String(item.offShelfReason || item.off_shelf_reason || '').trim()
const canPutOnSale = (item) => productStatus(item) === 'OFF_SHELF' && restorableOffShelfSources.includes(offShelfBy(item))
const offShelfTip = (item) => {
  if (productStatus(item) !== 'OFF_SHELF') return ''
  const source = offShelfBy(item)
  if (source === 'ADMIN') return '管理员已下架'
  if (source === 'SYSTEM') return '系统已下架'
  if (!source) return '暂不可上架'
  return ''
}

const normalizeStatusError = (error, status = 'ON_SALE') => {
  const message = String(error?.message || '')
  if (/BANNED|PERM_BANNED|PERMANENT_BANNED|\u5c01\u7981/i.test(message)) return userTradeRestrictionMessage('BANNED', 'sale')
  if (/DISABLED|\u7981\u7528/i.test(message)) return userTradeRestrictionMessage('DISABLED', 'sale')
  if (/CANCELED|CANCELLED|\u6ce8\u9500/i.test(message)) return userTradeRestrictionMessage('CANCELED', 'sale')
  if (/管理员下架|admin/i.test(message)) return '管理员下架的商品不能自行上架'
  if (/系统下架|system/i.test(message)) return '系统下架的商品不能自行上架'
  if (/403|FORBIDDEN|PERMISSION|VERIFY|VERIFIED|\u6743\u9650/i.test(message)) {
    return status === 'ON_SALE' ? '\u8d26\u53f7\u72b6\u6001\u4e0d\u6ee1\u8db3\u4e0a\u67b6\u6761\u4ef6' : '\u8d26\u53f7\u72b6\u6001\u4e0d\u6ee1\u8db3\u4e0b\u67b6\u6761\u4ef6'
  }
  return message || (status === 'ON_SALE' ? '\u4e0a\u67b6\u5931\u8d25\uff0c\u8bf7\u7a0d\u540e\u91cd\u8bd5' : '\u4e0b\u67b6\u5931\u8d25\uff0c\u8bf7\u7a0d\u540e\u91cd\u8bd5')
}

const ensureCanPutOnSale = async () => {
  try {
    const cached = getUser() || {}
    const current = { ...cached, ...(await getCurrentUser()) }
    setUser(current)
    const restrictionMessage = userTradeRestrictionMessage(accountStatusOf(current), 'sale')
    if (restrictionMessage) {
      uni.showToast({ title: restrictionMessage, icon: 'none' })
      return false
    }
  } catch (error) {
    const cachedRestriction = userTradeRestrictionMessage(accountStatusOf(getUser() || {}), 'sale')
    uni.showToast({ title: cachedRestriction || normalizeStatusError(error, 'ON_SALE'), icon: 'none' })
    return false
  }
  return true
}

const editProduct = (item) => {
  uni.setStorageSync('PUBLISH_EDIT_PRODUCT_ID', item.id)
  uni.setStorageSync('PUBLISH_EDIT_PRODUCT_DATA', item)
  uni.setStorageSync('PUBLISH_EDIT_INTENT', '1')
  uni.switchTab({ url: '/pages/publish/publish' })
}

const changeStatus = (item, status) => {
  const title = status === 'ON_SALE' ? '确认重新上架？' : '确认下架商品？'
  uni.showModal({
    title,
    content: status === 'ON_SALE' ? '商品上架后将重新公开展示，请确认商品无违规、物品仍可正常交易。' : '下架后买家暂时看不到该商品。',
    success: async ({ confirm }) => {
      if (!confirm) return
      try {
        if (status === 'ON_SALE' && !(await ensureCanPutOnSale())) return
        await updateProductStatus(item.id, status)
        uni.showToast({ title: status === 'ON_SALE' ? '已上架' : '已下架', icon: 'success' })
        await loadData()
      } catch (error) {
        uni.showToast({ title: normalizeStatusError(error, status), icon: 'none' })
      }
    }
  })
}

const batchPutOnSale = () => {
  uni.showModal({
    title: '确认一键上架？',
    content: '商品上架后将重新公开展示，请确认商品无违规、物品仍可正常交易。',
    success: async ({ confirm }) => {
      if (!confirm) return
      try {
        if (!(await ensureCanPutOnSale())) return
        const result = await batchPutOnSaleProducts()
        const count = Number(result?.count || 0)
        uni.showToast({ title: count > 0 ? `已上架${count}件商品` : '暂无可上架商品', icon: count > 0 ? 'success' : 'none' })
        await loadData()
      } catch (error) {
        uni.showToast({ title: normalizeStatusError(error, 'ON_SALE'), icon: 'none' })
      }
    }
  })
}

const deleteMyProduct = (item) => {
  uni.showModal({
    title: '确认删除商品？',
    content: '删除后商品将不再展示，且无法由买家预约。',
    success: async ({ confirm }) => {
      if (!confirm) return
      try {
        await deleteProduct(item.id)
        uni.showToast({ title: '已删除商品', icon: 'success' })
        await loadData()
      } catch (error) {
        uni.showToast({ title: '删除失败，请稍后重试', icon: 'none' })
      }
    }
  })
}

onShow(loadData)
onPullDownRefresh(loadData)
</script>

<style scoped>
.page { min-height: 100vh; padding: 28rpx; background: #f5f8f6; box-sizing: border-box; }
.summary { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24rpx; padding: 28rpx; border-radius: 26rpx; background: #fff; box-shadow: 0 8rpx 24rpx rgba(28, 68, 52, .05); }
.title { color: #26342f; font-size: 38rpx; font-weight: 700; }
.muted, .load-state { color: #929c98; font-size: 23rpx; }
.summary-actions { display: flex; flex-shrink: 0; align-items: center; gap: 12rpx; margin-left: auto; }
.publish, .empty-button, .batch-on-sale { height: 66rpx; border-radius: 999rpx; font-size: 26rpx; line-height: 66rpx; }
.publish, .empty-button { background: #23734f; color: #fff; }
.publish { width: 128rpx; margin: 0; flex-shrink: 0; }
.batch-on-sale { min-width: 160rpx; padding: 0 20rpx; background: #fff7e6; color: #a96500; }
.list { display: flex; flex-direction: column; gap: 22rpx; }
.product { display: flex; flex-direction: column; gap: 18rpx; padding: 22rpx; border-radius: 24rpx; background: #fff; box-shadow: 0 8rpx 24rpx rgba(28, 68, 52, .05); }
.product-main { display: flex; min-width: 0; }
.cover { display: flex; width: 176rpx; height: 176rpx; margin-right: 18rpx; flex-shrink: 0; align-items: center; justify-content: center; border-radius: 14rpx; background: #e8efeb; color: #9aa5a1; font-size: 23rpx; }
.content { display: flex; min-width: 0; flex: 1; flex-direction: column; justify-content: space-between; }
.top, .actions { display: flex; align-items: center; }
.top { justify-content: space-between; gap: 12rpx; }
.name { display: -webkit-box; overflow: hidden; flex: 1; color: #26342f; font-size: 29rpx; font-weight: 700; line-height: 1.4; -webkit-box-orient: vertical; -webkit-line-clamp: 2; }
.status { flex-shrink: 0; padding: 6rpx 12rpx; border-radius: 999rpx; background: #e7f4ec; color: #23734f; font-size: 21rpx; }
.status.SOLD { background: #f4e9e7; color: #b85d45; }
.status.LOCKED { background: #fff2d8; color: #9a6a1d; }
.status.OFF_SHELF { background: #eef0f0; color: #7a817e; }
.meta { color: #89938f; font-size: 23rpx; }
.off-shelf-reason { overflow: hidden; color: #8a6a35; font-size: 23rpx; line-height: 1.35; text-overflow: ellipsis; white-space: nowrap; }
.price { color: #e36a3e; font-size: 34rpx; font-weight: 700; }
.actions { justify-content: flex-end; gap: 12rpx; padding-top: 16rpx; border-top: 1rpx solid #edf2ef; }
.action { min-width: 104rpx; height: 58rpx; padding: 0 18rpx; border-radius: 999rpx; background: #edf4f1; color: #23734f; font-size: 24rpx; line-height: 58rpx; }
.muted-action { background: #f4f1ed; color: #9a7745; }
.primary-action { background: #23734f; color: #fff; }
.danger-action { background: #fff1ef; color: #d85c45; }
.off-shelf-tip { display: flex; min-width: 132rpx; height: 58rpx; align-items: center; justify-content: center; padding: 0 18rpx; border-radius: 999rpx; background: #f3f0ec; color: #9a7745; font-size: 23rpx; }
.empty { margin-top: 140rpx; padding: 44rpx 28rpx; border-radius: 18rpx; background: #fff; text-align: center; }
.empty-title { margin-bottom: 12rpx; color: #26342f; font-size: 32rpx; font-weight: 700; }
.empty-button { width: 180rpx; margin-top: 28rpx; }
.load-state { padding: 28rpx; text-align: center; }
</style>
