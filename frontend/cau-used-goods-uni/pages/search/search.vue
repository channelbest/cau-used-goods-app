<template>
  <view class="page">
    <view class="search-row">
      <input v-model.trim="filters.keyword" class="search-input" confirm-type="search" placeholder="搜索商品标题或描述" @confirm="search" />
      <button class="search-button" @click="search">搜索</button>
    </view>
    <scroll-view scroll-x class="category-scroll">
      <view class="category-row">
        <text v-for="item in categories" :key="item.id" class="pill" :class="{ active: filters.categoryId === item.id }" @click="selectCategory(item.id)">{{ item.name }}</text>
      </view>
    </scroll-view>
    <view class="filters">
      <view class="filter-line">
        <text>价格</text>
        <input v-model="filters.minPrice" type="digit" placeholder="最低价" />
        <text>-</text>
        <input v-model="filters.maxPrice" type="digit" placeholder="最高价" />
      </view>
      <view class="filter-line">
        <text>成色</text>
        <picker :range="conditionOptions" @change="changeCondition"><view class="picker">{{ filters.conditionLevel || '不限成色' }} ›</view></picker>
      </view>
      <view class="filter-line">
        <text>排序</text>
        <picker :range="sortOptions.map((item) => item.label)" @change="changeSort"><view class="picker">{{ currentSortLabel }} ›</view></picker>
      </view>
      <button class="apply" @click="search">应用筛选</button>
    </view>
    <view class="result-head"><text class="title">搜索结果</text><text class="muted">共 {{ total }} 件</text></view>
    <view v-if="products.length" class="grid"><ProductCard v-for="item in products" :key="item.id" :product="item" /></view>
    <view v-else-if="!loading" class="empty-state">暂时没有找到合适的商品</view>
    <view class="load-state">{{ loading ? '正在加载...' : finished && products.length ? '已经到底啦' : '' }}</view>
  </view>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { onLoad, onReachBottom } from '@dcloudio/uni-app'
import ProductCard from '../../components/ProductCard.vue'
import { listCategories, listProducts } from '../../api/product'
import { buildCategoryMap, formatProduct, withAllCategory } from '../../utils/product-format'

const categories = ref([])
const rawProducts = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const conditionOptions = ['不限成色', '全新', '九成新', '八成新', '七成新', '有明显使用痕迹']
const sortOptions = [{ label: '最新发布', value: 'newest' }, { label: '价格从低到高', value: 'price_asc' }, { label: '价格从高到低', value: 'price_desc' }, { label: '收藏最多', value: 'popular' }]
const filters = reactive({ keyword: '', categoryId: 0, minPrice: '', maxPrice: '', conditionLevel: '', sort: 'newest' })
const products = computed(() => rawProducts.value.map((item) => formatProduct(item, buildCategoryMap(categories.value))))
const finished = computed(() => rawProducts.value.length >= total.value && total.value > 0)
const currentSortLabel = computed(() => sortOptions.find((item) => item.value === filters.sort)?.label || '最新发布')

const loadProducts = async (reset = false) => {
  if (loading.value || (!reset && finished.value)) return
  loading.value = true
  try {
    const nextPage = reset ? 1 : page.value
    const result = await listProducts({ ...filters, page: nextPage, pageSize: 8 })
    rawProducts.value = reset ? result.list : rawProducts.value.concat(result.list)
    total.value = result.total
    page.value = nextPage + 1
  } catch (error) {
    uni.showToast({ title: error.message, icon: 'none' })
  } finally {
    loading.value = false
  }
}
const search = () => loadProducts(true)
const selectCategory = (categoryId) => { filters.categoryId = categoryId; search() }
const changeCondition = ({ detail }) => { filters.conditionLevel = conditionOptions[detail.value] === '不限成色' ? '' : conditionOptions[detail.value] }
const changeSort = ({ detail }) => { filters.sort = sortOptions[detail.value].value; search() }

onLoad(async (options) => {
  filters.keyword = options.keyword || ''
  filters.categoryId = Number(options.categoryId || 0)
  try { categories.value = withAllCategory(await listCategories()); search() } catch (error) { uni.showToast({ title: error.message, icon: 'none' }) }
})
onReachBottom(() => loadProducts())
</script>

<style scoped>
.page { min-height: 100vh; padding: 28rpx; box-sizing: border-box; }
.search-row, .filter-line, .result-head { display: flex; align-items: center; }
.search-row { gap: 14rpx; }
.search-input { flex: 1; height: 76rpx; padding: 0 22rpx; border-radius: 16rpx; background: #fff; }
.search-button, .apply { border-radius: 14rpx; background: #23734f; color: #fff; font-size: 27rpx; }
.search-button { width: 132rpx; height: 76rpx; line-height: 76rpx; }
.category-scroll { width: 100%; margin: 22rpx 0; white-space: nowrap; }
.category-row { display: inline-flex; gap: 12rpx; }
.pill { padding: 12rpx 20rpx; border-radius: 999rpx; background: #fff; color: #65706c; font-size: 24rpx; }
.pill.active { background: #23734f; color: #fff; }
.filters { padding: 10rpx 20rpx 20rpx; border-radius: 18rpx; background: #fff; }
.filter-line { min-height: 72rpx; gap: 16rpx; border-bottom: 2rpx solid #f0f2f1; color: #53605b; }
.filter-line input { width: 160rpx; height: 54rpx; padding: 0 12rpx; border-radius: 12rpx; background: #f6f8f5; text-align: center; }
.picker { color: #23734f; }
.apply { margin-top: 18rpx; height: 66rpx; line-height: 66rpx; background: #e7f4ec; color: #23734f; }
.result-head { justify-content: space-between; padding: 30rpx 0 18rpx; }
.title { font-size: 32rpx; font-weight: 700; }
.muted, .load-state { color: #929c98; font-size: 23rpx; }
.grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 18rpx; }
.load-state { padding: 28rpx; text-align: center; }
</style>
