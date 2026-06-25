<template>
  <view class="page">
    <view v-if="activeTab === 'data'" class="content">
      <view class="admin-hero">
        <text class="hero-eyebrow">CAU CAMPUS MARKET</text>
        <text class="page-title">数据</text>
        <text class="hero-copy">查看平台交易、商品和风险概况</text>
      </view>

      <view class="todo-summary">
        <view class="todo-summary-title">核心数据</view>
        <view class="todo-summary-grid">
          <view class="todo-summary-item">
            <view class="todo-summary-value">{{ money(orderOverview.totalCompletedAmount) }}</view>
            <view class="todo-summary-label">成交金额</view>
          </view>
          <view class="todo-summary-item">
            <view class="todo-summary-value">{{ orderOverview.completedOrders || 0 }}</view>
            <view class="todo-summary-label">成交订单</view>
          </view>
        </view>
      </view>

      <view class="section-title">平台概况</view>
      <view class="overview-grid">
        <view class="overview-card">
          <view class="overview-head">
            <text class="overview-title">用户情况</text>
            <text class="overview-tag">认证率 {{ percent(userOverview.verifiedUsers, userOverview.totalUsers) }}</text>
          </view>
          <view class="overview-main">{{ userOverview.totalUsers || 0 }}</view>
          <view class="overview-sub">注册用户</view>
          <view class="overview-row">
            <text>已认证</text>
            <text>{{ userOverview.verifiedUsers || 0 }}</text>
          </view>
          <view class="overview-row">
            <text>审核中</text>
            <text>{{ userOverview.pendingUsers || 0 }}</text>
          </view>
        </view>

        <view class="overview-card">
          <view class="overview-head">
            <text class="overview-title">商品流通</text>
            <text class="overview-tag">在售 {{ percent(productOverview.onSaleProducts, productOverview.totalProducts) }}</text>
          </view>
          <view class="overview-main">{{ productOverview.onSaleProducts || 0 }}</view>
          <view class="overview-sub">在售商品</view>
          <view class="overview-row">
            <text>已下架</text>
            <text>{{ productOverview.offShelfProducts || 0 }}</text>
          </view>
          <view class="overview-row">
            <text>已成交</text>
            <text>{{ productOverview.soldProducts || 0 }}</text>
          </view>
        </view>

        <view class="overview-card">
          <view class="overview-head">
            <text class="overview-title">商品热度</text>
            <text class="overview-tag">浏览 / 收藏</text>
          </view>
          <view class="overview-main">{{ productOverview.totalViews || 0 }}</view>
          <view class="overview-sub">浏览量</view>
          <view class="overview-row">
            <text>收藏数</text>
            <text>{{ productOverview.totalFavorites || 0 }}</text>
          </view>
          <view class="overview-row">
            <text>平均价格</text>
            <text>{{ money(productOverview.averagePrice) }}</text>
          </view>
        </view>

        <view class="overview-card">
          <view class="overview-head">
            <text class="overview-title">风控情况</text>
            <text :class="['overview-tag', riskTodoCount ? 'danger' : '']">
              {{ riskTodoCount ? '需处理' : '正常' }}
            </text>
          </view>
          <view :class="['overview-main', riskTodoCount ? 'danger' : '']">
            {{ riskTodoCount }}
          </view>
          <view class="overview-sub">待处理申诉和举报</view>
          <view class="overview-row">
            <text>商品举报/申诉</text>
            <text>{{ productRiskCount }}</text>
          </view>
          <view class="overview-row">
            <text>用户举报/申诉</text>
            <text>{{ userRiskCount }}</text>
          </view>
          <view class="overview-row">
            <text>订单举报/申诉</text>
            <text>{{ orderRiskCount }}</text>
          </view>
        </view>
      </view>

      <view class="section-title">近 7 天发布 / 成交趋势</view>
      <view class="trend-card">
        <view class="trend-legend">
          <view class="legend-item"><text class="legend-dot publish"></text>发布</view>
          <view class="legend-item"><text class="legend-dot deal"></text>成交</view>
        </view>
        <view class="trend-bars">
          <view v-for="item in productTrend" :key="item.date" class="trend-item">
            <view class="trend-bar-wrap">
              <view class="trend-bar publish" :style="trendBarStyle(item.count)"></view>
              <view class="trend-bar deal" :style="trendBarStyle(item.completedCount || 0)"></view>
            </view>
            <view class="trend-count">{{ item.count || 0 }}/{{ item.completedCount || 0 }}</view>
            <view class="trend-date">{{ shortDate(item.date) }}</view>
          </view>
        </view>
        <view v-if="!productTrend.length" class="empty">暂无趋势数据</view>
        <view class="trend-note">成交为当日发布商品中已售出的数量</view>
      </view>

      <view class="section-title category-title">
        <text>{{ selectedPrimaryCategory ? `${selectedPrimaryCategory.name} · 二级分类` : '商品分类统计' }}</text>
        <text v-if="selectedPrimaryCategory" class="category-back" @click="clearSelectedPrimary">返回一级</text>
      </view>
      <view class="category-list">
        <view
          v-for="item in categoryRows"
          :key="item.categoryId || item.categoryName"
          class="category-item"
          @click="openPrimaryCategory(item)"
        >
          <view class="category-head">
            <text class="category-name">{{ item.categoryName }}</text>
            <text class="category-count">{{ item.productCount || 0 }} 件{{ !selectedPrimaryCategory && item.childrenCount ? ' ›' : '' }}</text>
          </view>
          <view class="category-progress">
            <view class="category-bar" :style="categoryBarStyle(item.productCount)"></view>
            <view class="category-bar sale" :style="categoryOnSaleBarStyle(item)"></view>
          </view>
          <view class="category-meta">
            <text>在售 {{ item.onSaleCount || 0 }}，占比 {{ percent(item.onSaleCount, item.productCount) }}</text>
            <text>均价 {{ money(item.averagePrice) }}</text>
          </view>
        </view>
        <view v-if="!categoryRows.length" class="empty">{{ selectedPrimaryCategory ? '暂无二级分类统计' : '暂无分类统计' }}</view>
      </view>

      <view class="section-title">运营提醒</view>
      <view class="panel">
        <view class="row">
          <text>待确认订单</text>
          <text>{{ orderOverview.pendingConfirmOrders || 0 }}</text>
        </view>
        <view class="row">
          <text>待面交订单</text>
          <text>{{ orderOverview.waitMeetOrders || 0 }}</text>
        </view>
        <view class="row">
          <text>已下架商品</text>
          <text>{{ productOverview.offShelfProducts || 0 }}</text>
        </view>
        <view class="row">
          <text>异常关闭订单</text>
          <text>{{ orderOverview.exceptionClosedOrders || 0 }}</text>
        </view>
      </view>
    </view>

    <view v-else-if="activeTab === 'review'" class="content">
      <view class="admin-hero">
        <text class="hero-eyebrow">CAU CAMPUS MARKET</text>
        <text class="page-title">审核</text>
        <text class="hero-copy">优先处理认证、举报和申诉</text>
      </view>

      <view class="todo-panel">
        <view class="todo-card urgent" @click="goPage('/pages/admin-students/admin-students')">
          <view :class="['todo-number', userOverview.pendingUsers ? 'danger' : '']">{{ userOverview.pendingUsers || 0 }}</view>
          <view class="todo-label">待认证审核</view>
        </view>
        <view class="todo-card" @click="goPage('/pages/admin-reports/admin-reports')">
          <view :class="['todo-number', riskTodoCount ? 'danger' : '']">{{ riskTodoCount }}</view>
          <view class="todo-label">待处理申诉和举报</view>
        </view>
      </view>

      <view class="section-title">管理</view>
      <view class="manage-grid">
        <view class="manage-card" @click="goPage('/pages/admin-products/admin-products')">
          <view class="manage-title">商品管理</view>
          <view class="manage-desc">上架 / 下架</view>
        </view>
        <view class="manage-card" @click="goPage('/pages/admin-users/admin-users')">
          <view class="manage-title">用户管理</view>
          <view class="manage-desc">用户与角色</view>
        </view>
        <view class="manage-card" @click="goPage('/pages/admin-sensitive/admin-sensitive')">
          <view class="manage-title">敏感词管理</view>
          <view class="manage-desc">审核词库</view>
        </view>
        <view class="manage-card" @click="goPage('/pages/admin-categories/admin-categories')">
          <view class="manage-title">标签管理</view>
          <view class="manage-desc">商品分类标签</view>
        </view>
        <view class="manage-card" @click="goPage('/pages/admin-announcements/admin-announcements')">
          <view class="manage-title">公告管理</view>
          <view class="manage-desc">发布 / 下线</view>
        </view>
        <view class="manage-card" @click="goPage('/pages/admin-orders/admin-orders')">
          <view class="manage-title">订单管理</view>
          <view class="manage-desc">异常关闭</view>
        </view>
        <view class="manage-card" @click="goPage('/pages/admin-stats/admin-stats')">
          <view class="manage-title">申诉统计</view>
          <view class="manage-desc">申诉数据</view>
        </view>
      </view>

    </view>

    <view v-else class="content">
      <view class="admin-hero">
        <text class="hero-eyebrow">CAU CAMPUS MARKET</text>
        <text class="page-title">我的</text>
        <text class="hero-copy">管理员资料与审计日志</text>
      </view>
      <view class="profile-card">
        <image v-if="adminAvatarUrl" class="avatar" :src="adminAvatarUrl" mode="aspectFill" />
        <view v-else class="avatar placeholder">管</view>
        <view class="profile-main">
          <view class="nickname">{{ currentUser.nickname || '管理员' }}</view>
          <view class="role">管理员</view>
        </view>
      </view>

      <view class="menu-list">
        <view class="menu-item" @click="goPage('/pages/profile-edit/profile-edit')">
          <view class="menu-title">个人资料编辑</view>
          <text class="arrow">›</text>
        </view>
        <view class="menu-item" @click="goPage('/pages/admin-logs/admin-logs')">
          <view class="menu-title">管理员日志</view>
          <text class="arrow">›</text>
        </view>
      </view>

      <button class="logout-button" @click="logout">退出登录</button>
    </view>

    <view class="tabbar">
      <view :class="['tab-item', activeTab === 'data' ? 'active' : '']" @click="switchTab('data')">数据</view>
      <view :class="['tab-item', activeTab === 'review' ? 'active' : '']" @click="switchTab('review')">审核</view>
      <view :class="['tab-item', activeTab === 'mine' ? 'active' : '']" @click="switchTab('mine')">我的</view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import {
  getCategoryDistribution,
  getAdminCategories,
  getAppealOverview,
  getOrderOverview,
  getProductOverview,
  getProductTrend,
  getReportOverview,
  getUserOverview
} from '../../api/admin'
import { getCurrentUser } from '../../api/auth'
import { clearAuth, getUser, setUser } from '../../utils/auth'
import { BASE_URL } from '../../utils/request'

const activeTab = ref('data')
const currentUser = ref(getUser() || {})
const userOverview = ref({})
const productOverview = ref({})
const orderOverview = ref({})
const reportOverview = ref({})
const appealOverview = ref({})
const categoryDistribution = ref([])
const adminCategories = ref([])
const selectedPrimaryId = ref('')
const productTrend = ref([])

const firstValue = (values) => {
  for (const value of values) {
    if (value !== undefined && value !== null && value !== '') return value
  }
  return 0
}

const toNumber = (value) => {
  const num = Number(value || 0)
  return Number.isFinite(num) ? num : 0
}

const trendPublishCount = (item) => {
  return toNumber(firstValue([item?.count, item?.publishCount, item?.publish_count]))
}

const trendCompletedCount = (item) => {
  return toNumber(firstValue([item?.completed_count,item?.completedCount]))
}

const trendItems = (result) => {
  if (Array.isArray(result)) return result
  if (Array.isArray(result?.list)) return result.list
  if (Array.isArray(result?.items)) return result.items
  if (Array.isArray(result?.data)) return result.data
  if (Array.isArray(result?.data?.list)) return result.data.list
  if (Array.isArray(result?.data?.items)) return result.data.items
  return []
}

const normalizeProductTrend = (result) => {
  return trendItems(result).map((raw) => {
    const item = raw || {}
    return {
      ...item,
      date: firstValue([item?.date, item?.trendDate, item?.trend_date]),
      count: trendPublishCount(item),
      completedCount: trendCompletedCount(item)
    }
  })
}

const riskTodoCount = computed(() => {
  return Number(reportOverview.value.pendingReports || 0) + Number(appealOverview.value.pendingAppeals || 0)
})

const productRiskCount = computed(() => {
  return Number(reportOverview.value.productReports || 0) + Number(appealOverview.value.productAppeals || 0)
})

const userRiskCount = computed(() => {
  return Number(reportOverview.value.userReports || 0) + Number(appealOverview.value.userAppeals || 0)
})

const orderRiskCount = computed(() => {
  return Number(reportOverview.value.orderReports || 0) + Number(appealOverview.value.orderAppeals || 0)
})

const adminAvatarUrl = computed(() => {
  const url = currentUser.value?.avatarUrl || currentUser.value?.avatar_url || ''
  if (!url) return ''
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  if (url.startsWith('/uploads/')) return BASE_URL + url
  return url
})

const idOf = (item) => item?.id || item?.categoryId || item?.category_id || ''
const nameOf = (item) => item?.name || item?.categoryName || item?.category_name || '分类'
const parentIdOf = (item) => Number(item?.parentId || item?.parent_id || 0)
const statCategoryId = (item) => item?.categoryId || item?.category_id || item?.id || ''
const categoryMap = computed(() => {
  const map = {}
  adminCategories.value.forEach((item) => {
    const id = idOf(item)
    if (id) map[String(id)] = item
  })
  return map
})
const primaryCategories = computed(() => adminCategories.value.filter((item) => parentIdOf(item) === 0))
const selectedPrimaryCategory = computed(() => {
  if (!selectedPrimaryId.value) return null
  const category = categoryMap.value[String(selectedPrimaryId.value)]
  return category ? { ...category, id: idOf(category), name: nameOf(category) } : null
})

function emptyCategoryRow(categoryId, categoryName) {
  return {
    categoryId,
    categoryName,
    productCount: 0,
    onSaleCount: 0,
    averagePrice: 0,
    _totalPrice: 0,
    childrenCount: 0
  }
}

function addCategoryStat(target, item) {
  const productCount = Number(item?.productCount || item?.count || 0)
  const onSaleCount = Number(item?.onSaleCount || item?.on_sale_count || 0)
  const averagePrice = Number(item?.averagePrice || item?.average_price || 0)
  target.productCount += productCount
  target.onSaleCount += onSaleCount
  target._totalPrice += averagePrice * productCount
  target.averagePrice = target.productCount ? target._totalPrice / target.productCount : 0
}

const primaryCategoryRows = computed(() => {
  const rows = new Map()
  categoryDistribution.value.forEach((item) => {
    const categoryId = statCategoryId(item)
    const category = categoryMap.value[String(categoryId)]
    const primary = category && parentIdOf(category) ? categoryMap.value[String(parentIdOf(category))] : category
    const primaryId = idOf(primary) || categoryId || nameOf(item)
    const primaryName = primary ? nameOf(primary) : (item.categoryName || item.name || '分类')
    if (!rows.has(String(primaryId))) rows.set(String(primaryId), emptyCategoryRow(primaryId, primaryName))
    const row = rows.get(String(primaryId))
    if (category && parentIdOf(category)) row.childrenCount += 1
    addCategoryStat(row, item)
  })
  return Array.from(rows.values()).sort((a, b) => Number(b.productCount || 0) - Number(a.productCount || 0))
})

const childCategoryRows = computed(() => {
  if (!selectedPrimaryId.value) return []
  const rows = []
  const children = adminCategories.value.filter((item) => Number(parentIdOf(item)) === Number(selectedPrimaryId.value))
  children.forEach((child) => {
    const childId = idOf(child)
    const row = emptyCategoryRow(childId, nameOf(child))
    categoryDistribution.value
      .filter((item) => String(statCategoryId(item)) === String(childId))
      .forEach((item) => addCategoryStat(row, item))
    if (row.productCount > 0) rows.push(row)
  })
  const directRow = emptyCategoryRow(`${selectedPrimaryId.value}-direct`, '未细分')
  categoryDistribution.value
    .filter((item) => String(statCategoryId(item)) === String(selectedPrimaryId.value))
    .forEach((item) => addCategoryStat(directRow, item))
  if (directRow.productCount > 0) rows.unshift(directRow)
  return rows.sort((a, b) => Number(b.productCount || 0) - Number(a.productCount || 0))
})

const categoryRows = computed(() => selectedPrimaryId.value ? childCategoryRows.value : primaryCategoryRows.value)

const loadAdminData = async () => {
  try {
    const [me, users, productStats, orders, reportStats, appealStats, categories, categoryTree, trend] = await Promise.all([
      getCurrentUser(),
      getUserOverview(),
      getProductOverview(),
      getOrderOverview(),
      getReportOverview(),
      getAppealOverview(),
      getCategoryDistribution(),
      getAdminCategories().catch(() => []),
      getProductTrend(7)
    ])
    currentUser.value = me || {}
    setUser(me || {})
    userOverview.value = users || {}
    productOverview.value = productStats || {}
    orderOverview.value = orders || {}
    reportOverview.value = reportStats || {}
    appealOverview.value = appealStats || {}
    categoryDistribution.value = categories?.list || categories || []
    adminCategories.value = categoryTree || []
    productTrend.value = normalizeProductTrend(trend)
  } catch (error) {
    uni.showToast({ title: error.message || '后台数据加载失败', icon: 'none' })
  }
}

onShow(loadAdminData)

const switchTab = (tab) => {
  activeTab.value = tab
}

const goPage = (url) => {
  uni.navigateTo({ url })
}

const openPrimaryCategory = (item) => {
  if (selectedPrimaryId.value || !item?.childrenCount) return
  selectedPrimaryId.value = item.categoryId
}

const clearSelectedPrimary = () => {
  selectedPrimaryId.value = ''
}

const logout = () => {
  clearAuth()
  uni.reLaunch({ url: '/pages/login/login' })
}

const money = (value) => {
  return `￥${Number(value || 0).toFixed(2)}`
}

const percent = (value, total) => {
  const num = Number(value || 0)
  const den = Number(total || 0)
  if (!den) return '0%'
  return `${Math.round((num / den) * 100)}%`
}

const trendBarStyle = (count) => {
  const max = Math.max(...productTrend.value.flatMap((item) => [
    Number(item.count || 0),
    Number(item.completedCount || 0)
  ]), 1)
  const height = Math.max(12, Math.round((Number(count || 0) / max) * 120))
  return `height: ${height}rpx;`
}

const shortDate = (date) => {
  if (!date) return ''
  const value = String(date)
  const datePart = value.includes('T') ? value.split('T')[0] : value
  const parts = datePart.split('-')
  if (parts.length < 3) return value
  return `${parts[1]}.${parts[2]}`
}

const categoryBarStyle = (count) => {
  const max = Math.max(...categoryRows.value.map((item) => Number(item.productCount || 0)), 1)
  const width = Math.max(8, Math.round((Number(count || 0) / max) * 100))
  return `width: ${width}%;`
}

const categoryOnSaleBarStyle = (item) => {
  const total = Number(item?.productCount || 0)
  const onSale = Number(item?.onSaleCount || 0)
  if (!total) return 'width: 0%;'
  return `width: ${Math.round((onSale / total) * 100)}%;`
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding-bottom: 140rpx;
  background: #f3f8f5;
  box-sizing: border-box;
}

.content {
  padding: 34rpx 28rpx;
}

.admin-hero {
  margin-bottom: 24rpx;
  padding: 30rpx 30rpx;
  border-radius: 28rpx;
  background: linear-gradient(145deg, #23734f, #2f8b62);
  color: #fff;
  box-shadow: 0 14rpx 34rpx rgba(31, 106, 73, .08);
}

.hero-eyebrow,
.page-title,
.hero-copy {
  display: block;
}

.hero-eyebrow {
  color: rgba(255, 255, 255, .72);
  font-size: 20rpx;
  letter-spacing: 3rpx;
}

.page-title {
  margin: 12rpx 0 0;
  font-size: 38rpx;
  line-height: 46rpx;
  font-weight: 700;
  color: #fff;
}

.hero-copy {
  margin-top: 12rpx;
  font-size: 24rpx;
  line-height: 36rpx;
  color: rgba(255, 255, 255, .78);
}

.section-title {
  margin: 34rpx 0 18rpx;
  font-size: 32rpx;
  font-weight: 700;
  color: #20352b;
}

.todo-summary {
  padding: 28rpx;
  border-radius: 26rpx;
  background: #fff;
  box-shadow: 0 12rpx 32rpx rgba(32, 53, 43, 0.05);
  box-sizing: border-box;
}

.todo-summary-title {
  font-size: 28rpx;
  font-weight: 700;
  color: #20352b;
}

.todo-summary-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20rpx;
  margin-top: 22rpx;
}

.todo-summary-item {
  min-height: 150rpx;
  padding: 24rpx;
  border-radius: 24rpx;
  background: #f3f8f5;
  box-sizing: border-box;
}

.todo-summary-value {
  font-size: 54rpx;
  line-height: 60rpx;
  font-weight: 700;
  color: #23734f;
}

.todo-summary-value.warning {
  color: #fde68a;
}

.todo-summary-value.danger {
  color: #fca5a5;
}

.todo-summary-label {
  margin-top: 14rpx;
  font-size: 26rpx;
  color: #667085;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20rpx;
}

.overview-card {
  min-height: 244rpx;
  padding: 26rpx;
  border-radius: 26rpx;
  background: #fff;
  box-shadow: 0 12rpx 32rpx rgba(32, 53, 43, 0.05);
  box-sizing: border-box;
}

.overview-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10rpx;
}

.overview-title {
  font-size: 26rpx;
  font-weight: 700;
  color: #20352b;
}

.overview-tag {
  max-width: 150rpx;
  padding: 6rpx 12rpx;
  border-radius: 999rpx;
  background: #e8f8ef;
  color: #18a45a;
  font-size: 20rpx;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.overview-tag.danger {
  background: #fee2e2;
  color: #ef4444;
}

.overview-main {
  margin-top: 20rpx;
  font-size: 52rpx;
  line-height: 58rpx;
  font-weight: 700;
  color: #20352b;
}

.overview-main.danger {
  color: #ef4444;
}

.overview-sub {
  margin-top: 8rpx;
  margin-bottom: 14rpx;
  font-size: 24rpx;
  color: #8a96a8;
}

.overview-row {
  display: flex;
  justify-content: space-between;
  padding-top: 10rpx;
  font-size: 24rpx;
  color: #667085;
}

.quick-list {
  border-radius: 26rpx;
  background: #fff;
  box-shadow: 0 12rpx 32rpx rgba(32, 53, 43, 0.05);
  overflow: hidden;
}

.quick-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 116rpx;
  padding: 24rpx 28rpx;
  border-bottom: 1rpx solid #eef0f3;
  box-sizing: border-box;
}

.quick-item:last-child {
  border-bottom: 0;
}

.quick-title {
  font-size: 30rpx;
  font-weight: 700;
  color: #1f2933;
}

.quick-desc {
  margin-top: 8rpx;
  font-size: 24rpx;
  color: #8a96a8;
}

.quick-count {
  min-width: 44rpx;
  color: #17a84b;
  font-size: 34rpx;
  font-weight: 700;
  text-align: right;
}

.quick-count.danger {
  color: #ef4444;
}

.trend-card {
  padding: 28rpx 22rpx;
  border-radius: 26rpx;
  background: #fff;
  box-shadow: 0 12rpx 32rpx rgba(32, 53, 43, 0.05);
}

.trend-legend {
  display: flex;
  justify-content: flex-end;
  gap: 24rpx;
  margin-bottom: 18rpx;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8rpx;
  font-size: 22rpx;
  color: #667085;
}

.legend-dot {
  width: 14rpx;
  height: 14rpx;
  border-radius: 50%;
}

.legend-dot.publish {
  background: #17a84b;
}

.legend-dot.deal {
  background: #3b82f6;
}

.trend-bars {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  min-height: 220rpx;
}

.trend-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 0;
}

.trend-bar-wrap {
  display: flex;
  align-items: flex-end;
  gap: 6rpx;
  height: 132rpx;
}

.trend-bar {
  width: 22rpx;
  min-height: 12rpx;
  border-radius: 999rpx 999rpx 8rpx 8rpx;
}

.trend-bar.publish {
  background: linear-gradient(180deg, #22c55e 0%, #16a34a 100%);
}

.trend-bar.deal {
  background: linear-gradient(180deg, #60a5fa 0%, #2563eb 100%);
}

.trend-count {
  margin-top: 12rpx;
  font-size: 20rpx;
  color: #1f2933;
  font-weight: 700;
}

.trend-date {
  margin-top: 6rpx;
  font-size: 20rpx;
  color: #98a2b3;
}

.trend-note {
  margin-top: 18rpx;
  padding-top: 18rpx;
  border-top: 1rpx solid #eef0f3;
  font-size: 22rpx;
  color: #98a2b3;
}

.category-list {
  border-radius: 26rpx;
  background: #fff;
  box-shadow: 0 12rpx 32rpx rgba(32, 53, 43, 0.05);
  overflow: hidden;
}

.category-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20rpx;
}

.category-back {
  flex-shrink: 0;
  color: #23734f;
  font-size: 24rpx;
  font-weight: 600;
}

.category-item {
  padding: 24rpx 28rpx;
  border-bottom: 1rpx solid #eef0f3;
}

.category-item:last-child {
  border-bottom: 0;
}

.category-head,
.category-meta {
  display: flex;
  justify-content: space-between;
  gap: 20rpx;
}

.category-name {
  font-size: 28rpx;
  font-weight: 700;
  color: #1f2933;
}

.category-count {
  font-size: 26rpx;
  color: #17a84b;
}

.category-progress {
  position: relative;
  height: 12rpx;
  margin: 18rpx 0 14rpx;
  border-radius: 999rpx;
  background: #eef2f6;
  overflow: hidden;
}

.category-bar {
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  border-radius: 999rpx;
  background: #17a84b;
}

.category-bar.sale {
  background: #facc15;
}

.category-meta {
  font-size: 22rpx;
  color: #8a96a8;
}

.todo-panel {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20rpx;
}

.todo-card {
  min-height: 150rpx;
  padding: 30rpx;
  border-radius: 26rpx;
  background: #fff;
  border: 2rpx solid transparent;
  box-shadow: 0 12rpx 32rpx rgba(32, 53, 43, 0.05);
  box-sizing: border-box;
}

.todo-card.urgent {
  border-color: #79ddb0;
}

.todo-number {
  font-size: 50rpx;
  line-height: 56rpx;
  font-weight: 700;
  color: #17a84b;
}

.todo-number.danger {
  color: #ef4444;
}

.todo-label {
  margin-top: 14rpx;
  font-size: 26rpx;
  color: #667085;
}

.panel,
.menu-list,
.profile-card {
  border-radius: 26rpx;
  background: #fff;
  box-shadow: 0 12rpx 32rpx rgba(32, 53, 43, 0.05);
  overflow: hidden;
}

.row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20rpx;
  padding: 26rpx 30rpx;
  border-bottom: 1rpx solid #eef0f3;
  font-size: 28rpx;
  color: #475467;
}

.row text:first-child {
  flex: 1;
}

.row text:last-child {
  min-width: 52rpx;
  height: 44rpx;
  padding: 0 16rpx;
  border-radius: 999rpx;
  background: #f0fdf4;
  color: #16a34a;
  font-weight: 700;
  line-height: 44rpx;
  text-align: center;
}

.row:last-child {
  border-bottom: 0;
}

.empty {
  padding: 40rpx 28rpx;
  font-size: 28rpx;
  color: #98a2b3;
}

.menu-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 118rpx;
  padding: 24rpx 28rpx;
  border-bottom: 1rpx solid #eef0f3;
  box-sizing: border-box;
}

.menu-item:last-child {
  border-bottom: 0;
}

.menu-item.primary {
  min-height: 132rpx;
}

.menu-title {
  font-size: 30rpx;
  font-weight: 700;
  color: #1f2933;
}

.menu-desc {
  margin-top: 8rpx;
  font-size: 24rpx;
  color: #8a96a8;
}

.menu-side {
  display: flex;
  align-items: center;
  gap: 14rpx;
}

.badge {
  min-width: 38rpx;
  height: 38rpx;
  padding: 0 12rpx;
  border-radius: 999rpx;
  background: #fee2e2;
  color: #ef4444;
  font-size: 24rpx;
  line-height: 38rpx;
  text-align: center;
}

.manage-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20rpx;
}

.manage-card {
  min-height: 156rpx;
  padding: 28rpx;
  border-radius: 26rpx;
  background: #fff;
  box-shadow: 0 12rpx 32rpx rgba(32, 53, 43, 0.05);
  box-sizing: border-box;
}

.manage-title {
  font-size: 30rpx;
  font-weight: 700;
  color: #20352b;
}

.manage-desc {
  margin-top: 14rpx;
  font-size: 24rpx;
  color: #8a96a8;
}

.hint-panel {
  padding: 24rpx 28rpx;
  border-radius: 26rpx;
  background: #fff;
}

.hint-row {
  display: flex;
  align-items: flex-start;
  gap: 14rpx;
  margin-bottom: 16rpx;
  font-size: 24rpx;
  line-height: 36rpx;
  color: #667085;
}

.hint-row:last-child {
  margin-bottom: 0;
}

.hint-dot {
  width: 12rpx;
  height: 12rpx;
  margin-top: 12rpx;
  border-radius: 50%;
  background: #17a84b;
  flex-shrink: 0;
}

.arrow {
  font-size: 42rpx;
  color: #b2bdca;
}

.profile-card {
  display: flex;
  align-items: center;
  padding: 34rpx;
  margin-bottom: 24rpx;
}

.avatar {
  width: 116rpx;
  height: 116rpx;
  border-radius: 58rpx;
  margin-right: 24rpx;
  background: #e7edf3;
}

.placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #17a84b;
  font-size: 36rpx;
  font-weight: 700;
}

.profile-main {
  flex: 1;
  min-width: 0;
}

.nickname {
  font-size: 34rpx;
  font-weight: 700;
  color: #1f2933;
}

.role {
  margin-top: 10rpx;
  font-size: 26rpx;
  color: #667085;
}

.logout-button {
  margin-top: 28rpx;
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 22rpx;
  background: #fff;
  color: #ef4444;
  font-size: 30rpx;
  box-shadow: 0 12rpx 32rpx rgba(32, 53, 43, 0.05);
}

.tabbar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 10;
  display: flex;
  height: 116rpx;
  padding-bottom: env(safe-area-inset-bottom);
  background: #fff;
  border-top: 1rpx solid #eef0f3;
  box-shadow: 0 -12rpx 28rpx rgba(32, 53, 43, 0.06);
}

.tab-item {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28rpx;
  color: #667085;
}

.tab-item.active {
  color: #17a84b;
  font-weight: 700;
}
</style>
