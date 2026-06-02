<template>
  <view class="page">
    <view class="hero">
      <text class="eyebrow">CAU CAMPUS MARKET</text>
      <text class="headline">让闲置，在校园里重新发光</text>
      <view class="search" @click="goSearch">搜索教材、数码、生活用品</view>
    </view>

    <view class="section-head">
      <text class="section-title">逛分类</text>
      <text class="muted">快速找到需要的好物</text>
    </view>
    <scroll-view scroll-x class="category-scroll">
      <view class="category-row">
        <view v-for="item in visibleCategories" :key="item.id" class="category" @click="goCategory(item.id)">
          <view class="category-icon">{{ item.name.slice(0, 2) }}</view>
          <text>{{ item.name }}</text>
        </view>
      </view>
    </scroll-view>

    <view class="section-head">
      <text class="section-title">新鲜发布</text>
      <text class="more" @click="goSearch">筛选排序 ›</text>
    </view>
    <view v-if="products.length" class="grid">
      <ProductCard v-for="item in products" :key="item.id" :product="item" />
    </view>
    <view v-else-if="!loading" class="empty-state">暂时没有在售商品</view>
    <view class="load-state">{{ loading ? '正在加载...' : finished ? '已经到底啦' : '' }}</view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onPullDownRefresh, onReachBottom, onShow } from '@dcloudio/uni-app'
import ProductCard from '../../components/ProductCard.vue'
import { listCategories, listProducts } from '../../api/product'
import { buildCategoryMap, formatProduct } from '../../utils/product-format'

const categories = ref([])
const rawProducts = ref([])
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const products = computed(() => rawProducts.value.map((item) => formatProduct(item, buildCategoryMap(categories.value))))
const visibleCategories = computed(() => categories.value.filter((item) => Number(item.id) !== 0))
const finished = computed(() => rawProducts.value.length >= total.value && total.value > 0)

const loadCategories = async () => {
  categories.value = await listCategories()
}

const loadProducts = async (reset = false) => {
  if (loading.value || (!reset && finished.value)) return
  loading.value = true
  try {
    const nextPage = reset ? 1 : page.value
    const result = await listProducts({ page: nextPage, pageSize: 8, sort: 'newest' })
    rawProducts.value = reset ? result.list : rawProducts.value.concat(result.list)
    total.value = result.total
    page.value = nextPage + 1
  } catch (error) {
    uni.showToast({ title: error.message, icon: 'none' })
  } finally {
    loading.value = false
  }
}

const refresh = async () => {
  try {
    await loadCategories()
    await loadProducts(true)
  } finally {
    uni.stopPullDownRefresh()
  }
}

const goSearch = () => uni.navigateTo({ url: '/pages/search/search' })
const goCategory = (categoryId) => uni.navigateTo({ url: `/pages/category/category?categoryId=${categoryId}` })

onShow(() => {
  if (!rawProducts.value.length || uni.getStorageSync('PRODUCT_LIST_DIRTY')) {
    uni.removeStorageSync('PRODUCT_LIST_DIRTY')
    refresh()
  }
})
onReachBottom(() => loadProducts())
onPullDownRefresh(refresh)
</script>

<style scoped>
.page { min-height: 100vh; padding-bottom: 36rpx; }
.hero { padding: 92rpx 30rpx 34rpx; border-radius: 0 0 40rpx 40rpx; background: linear-gradient(145deg, #1f6a49, #328660); color: #fff; }
.eyebrow, .headline { display: block; }
.eyebrow { color: rgba(255,255,255,.7); font-size: 20rpx; letter-spacing: 3rpx; }
.headline { margin-top: 14rpx; font-size: 38rpx; font-weight: 700; }
.search { margin-top: 30rpx; padding: 24rpx; border-radius: 20rpx; background: #fff; color: #9ca7a3; }
.section-head { display: flex; align-items: center; justify-content: space-between; padding: 32rpx 28rpx 18rpx; }
.section-title { color: #26342f; font-size: 32rpx; font-weight: 700; }
.muted, .load-state { color: #929c98; font-size: 23rpx; }
.more { color: #23734f; font-size: 25rpx; }
.category-scroll { width: 100%; white-space: nowrap; }
.category-row { display: inline-flex; gap: 18rpx; padding: 2rpx 28rpx 8rpx; }
.category { width: 128rpx; color: #65706c; font-size: 23rpx; text-align: center; }
.category-icon { display: flex; width: 96rpx; height: 96rpx; margin: 0 auto 10rpx; align-items: center; justify-content: center; border-radius: 28rpx; background: #e6f3eb; color: #23734f; font-weight: 700; }
.grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 20rpx; padding: 0 28rpx; }
.load-state { padding: 28rpx; text-align: center; }
</style>
