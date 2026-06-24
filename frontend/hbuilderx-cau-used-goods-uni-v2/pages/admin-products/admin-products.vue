<template>
  <view class="page">
    <view class="search-box">
      <input
        v-model="searchKeyword"
        class="search-input"
        placeholder="搜索商品名称"
        placeholder-class="search-placeholder"
        confirm-type="search"
      />
      <button v-if="searchKeyword" class="clear-btn" size="mini" @click="searchKeyword = ''">清空</button>
    </view>

    <view v-if="filteredProducts.length === 0" class="empty">{{ products.length ? '没有匹配商品' : '暂无商品' }}</view>
    <view v-for="item in filteredProducts" :key="item.id" class="card" @click="goDetail(item.id)">
      <image v-if="coverImage(item)" class="cover" :src="coverImage(item)" mode="aspectFill" />
      <view v-else class="cover placeholder">商品</view>
      <view class="card-main">
        <view class="name">{{ item.title }}</view>
        <view class="desc">￥{{ item.price }} · {{ statusText(item.status) }}</view>
      </view>
      <view class="arrow">›</view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getAdminProducts } from '../../api/admin'
import { normalizeImage } from '../../utils/product-format'

const products = ref([])
const searchKeyword = ref('')

const filteredProducts = computed(() => {
  const keyword = searchKeyword.value.trim().toLowerCase()
  if (!keyword) return products.value
  return products.value.filter((item) => {
    const title = String(item?.title || '').toLowerCase()
    return title.includes(keyword)
  })
})

const load = async () => {
  try {
    const result = await getAdminProducts()
    products.value = result?.list || result?.items || []
  } catch (error) {
    uni.showToast({ title: error.message || '加载失败', icon: 'none' })
  }
}

onShow(load)

const statusText = (status) => {
  const map = { ON_SALE: '在售', OFF_SHELF: '已下架', LOCKED: '交易锁定', SOLD: '已售出', DELETED: '已删除' }
  return map[status] || status || '未知'
}

const coverImage = (item) => {
  return normalizeImage(item?.images?.[0] || '')
}

const goDetail = (id) => {
  uni.navigateTo({
    url: `/pages/detail/detail?id=${id}&adminView=1&readonly=1`
  })
}
</script>

<style scoped>
.page { min-height: 100vh; padding: 24rpx; background: #f5f6f8; box-sizing: border-box; }
.search-box { display: flex; align-items: center; gap: 14rpx; margin-bottom: 18rpx; }
.search-input { flex: 1; height: 76rpx; padding: 0 24rpx; border-radius: 16rpx; background: #fff; color: #1f2933; font-size: 26rpx; box-sizing: border-box; }
.search-placeholder { color: #98a2b3; }
.clear-btn { flex-shrink: 0; height: 76rpx; line-height: 76rpx; margin: 0; padding: 0 24rpx; border-radius: 16rpx; background: #eef2f6; color: #667085; font-size: 24rpx; }
.card, .empty { padding: 28rpx; border-radius: 16rpx; background: #fff; margin-bottom: 18rpx; }
.card { display: flex; align-items: center; }
.cover { width: 112rpx; height: 112rpx; margin-right: 20rpx; border-radius: 12rpx; background: #eef2f6; flex-shrink: 0; }
.placeholder { display: flex; align-items: center; justify-content: center; color: #98a2b3; font-size: 22rpx; }
.card-main { flex: 1; min-width: 0; }
.name { font-size: 30rpx; font-weight: 700; color: #1f2933; }
.desc, .empty { margin-top: 8rpx; color: #667085; font-size: 26rpx; }
.arrow { margin-left: 20rpx; color: #b2bdca; font-size: 42rpx; }
</style>
