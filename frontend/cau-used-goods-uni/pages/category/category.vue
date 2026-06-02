<template>
  <view class="page">
    <view class="header"><text class="title">{{ categoryName }}</text><text class="muted">共 {{ total }} 件商品</text></view>
    <view class="sort-row">
      <text v-for="item in sortOptions" :key="item.value" class="sort" :class="{ active: sort === item.value }" @click="changeSort(item.value)">{{ item.label }}</text>
    </view>
    <view v-if="products.length" class="grid"><ProductCard v-for="item in products" :key="item.id" :product="item" /></view>
    <view v-else-if="!loading" class="empty-state">这个分类暂时没有在售商品</view>
    <view class="load-state">{{ loading ? '正在加载...' : finished && products.length ? '已经到底啦' : '' }}</view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad, onReachBottom } from '@dcloudio/uni-app'
import ProductCard from '../../components/ProductCard.vue'
import { listCategories, listProducts } from '../../api/product'
import { buildCategoryMap, formatProduct } from '../../utils/product-format'

const categories = ref([])
const rawProducts = ref([])
const categoryId = ref(0)
const sort = ref('newest')
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const sortOptions = [{ label: '最新发布', value: 'newest' }, { label: '价格最低', value: 'price_asc' }, { label: '价格最高', value: 'price_desc' }]
const categoryMap = computed(() => buildCategoryMap(categories.value))
const categoryName = computed(() => categoryMap.value[categoryId.value] || '分类商品')
const products = computed(() => rawProducts.value.map((item) => formatProduct(item, categoryMap.value)))
const finished = computed(() => rawProducts.value.length >= total.value && total.value > 0)

const loadProducts = async (reset = false) => {
  if (loading.value || (!reset && finished.value)) return
  loading.value = true
  try {
    const nextPage = reset ? 1 : page.value
    const result = await listProducts({ categoryId: categoryId.value, sort: sort.value, page: nextPage, pageSize: 8 })
    rawProducts.value = reset ? result.list : rawProducts.value.concat(result.list)
    total.value = result.total
    page.value = nextPage + 1
  } catch (error) { uni.showToast({ title: error.message, icon: 'none' }) } finally { loading.value = false }
}
const changeSort = (value) => { sort.value = value; loadProducts(true) }
onLoad(async (options) => {
  categoryId.value = Number(options.categoryId || 0)
  try { categories.value = await listCategories(); loadProducts(true) } catch (error) { uni.showToast({ title: error.message, icon: 'none' }) }
})
onReachBottom(() => loadProducts())
</script>

<style scoped>
.page { min-height: 100vh; padding: 28rpx; box-sizing: border-box; }
.header, .sort-row { display: flex; align-items: center; justify-content: space-between; }
.title { font-size: 38rpx; font-weight: 700; }.muted, .load-state { color: #929c98; font-size: 23rpx; }
.sort-row { gap: 14rpx; padding: 24rpx 0; }.sort { flex: 1; padding: 14rpx 0; border-radius: 14rpx; background: #fff; color: #65706c; text-align: center; }.sort.active { background: #23734f; color: #fff; }
.grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 18rpx; }.load-state { padding: 28rpx; text-align: center; }
</style>
