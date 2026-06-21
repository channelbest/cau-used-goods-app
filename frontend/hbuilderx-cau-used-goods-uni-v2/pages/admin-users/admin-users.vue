<template>
  <view class="page">
    <view class="section-title">用户管理</view>

    <view class="toolbar">
      <input v-model="filters.keyword" class="search" placeholder="搜索昵称、学号、姓名、手机号" confirm-type="search" @confirm="reload" />
      <button class="search-btn" :loading="loading" @click="reload">搜索</button>
    </view>

    <view class="filters">
      <view
        v-for="item in accountOptions"
        :key="item.value"
        :class="['chip', filters.accountStatus === item.value ? 'active' : '']"
        @click="setFilter('accountStatus', item.value)"
      >
        {{ item.label }}
      </view>
    </view>
    <view class="filters">
      <view
        v-for="item in authOptions"
        :key="item.value"
        :class="['chip', filters.authStatus === item.value ? 'active' : '']"
        @click="setFilter('authStatus', item.value)"
      >
        {{ item.label }}
      </view>
    </view>
    <view v-if="isSuperAdmin" class="filters">
      <view
        v-for="item in roleOptions"
        :key="item.value"
        :class="['chip', filters.role === item.value ? 'active' : '']"
        @click="setFilter('role', item.value)"
      >
        {{ item.label }}
      </view>
    </view>

    <view v-if="!loading && users.length === 0" class="empty">暂无用户</view>

    <view v-for="item in users" :key="item.id" class="card">
      <view class="row">
        <view class="avatar">{{ firstChar(item) }}</view>
        <view class="main" @click="openUserHome(item)">
          <view class="name">{{ userName(item) }}</view>
          <view class="desc">用户ID：{{ item.id }} · {{ roleText(item.role) }} · {{ accountText(item.accountStatus) }}</view>
          <view class="desc">认证：{{ authText(item.authStatus) }} · 学院：{{ item.college || '未认证' }}</view>
          <view class="desc">手机：{{ item.phone || '未填写' }}</view>
        </view>
        <view class="home-link" @click="openUserHome(item)">主页</view>
      </view>

      <view class="actions">
        <button class="mini" @click="toggleRelated(item)">关联信息</button>
        <button class="mini" @click="goStudentAuth(item)">认证状态</button>
        <button v-if="item.accountStatus === 'NORMAL'" class="mini warn" @click="changeStatus(item, 'DISABLED')">禁用</button>
        <button v-if="item.accountStatus === 'NORMAL' || item.accountStatus === 'DISABLED'" class="mini danger" @click="changeStatus(item, 'BANNED')">封禁</button>
        <button v-if="canRecoverUser(item)" class="mini ok" @click="changeStatus(item, 'NORMAL')">{{ recoverButtonText(item) }}</button>
        <button v-if="canChangeRole(item)" class="mini role" @click="changeRole(item)">{{ roleActionText(item) }}</button>
      </view>

      <view v-if="expandedId === item.id" class="related">
        <view class="stats" v-if="detail.stats">
          <view>商品 {{ detail.stats.productCount || 0 }}</view>
          <view>订单 {{ detail.stats.orderCount || 0 }}</view>
          <view>举报 {{ detail.stats.reportSubmittedCount || 0 }}/被举报 {{ detail.stats.reportedCount || 0 }}</view>
        </view>
        <view class="related-title">发布商品</view>
        <view v-if="related.products.length === 0" class="muted">暂无商品</view>
        <view v-for="product in related.products" :key="product.id" class="related-line" @click="openProduct(product)">
          {{ product.title }} · {{ productStatusText(product.status) }} · ¥{{ product.price }}
        </view>

        <view class="related-title">相关订单</view>
        <view v-if="related.orders.length === 0" class="muted">暂无订单</view>
        <view v-for="order in related.orders" :key="order.id" class="related-line">
          {{ order.productTitleSnapshot }} · {{ orderStatusText(order.status) }}
        </view>

        <view class="related-title">相关举报</view>
        <view v-if="related.reports.length === 0" class="muted">暂无举报</view>
        <view v-for="report in related.reports" :key="report.id" class="related-line">
          {{ report.relation === 'SUBMITTED' ? '发起' : '被举报' }} · {{ reportReasonText(report.reasonType) }} · {{ reportStatusText(report.status) }}
        </view>
      </view>
    </view>

    <view class="pager">
      <button class="page-btn" :disabled="page <= 1 || loading" @click="prevPage">上一页</button>
      <text class="page-text">{{ page }} / {{ totalPages }}</text>
      <button class="page-btn" :disabled="page >= totalPages || loading" @click="nextPage">下一页</button>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import {
  getAdminUsers,
  getAdminUserDetail,
  getAdminUserProducts,
  getAdminUserOrders,
  getAdminUserReports,
  updateAdminUserStatus,
  updateAdminUserRole
} from '../../api/admin'
import { getUser } from '../../utils/auth'

const users = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = 10
const total = ref(0)
const expandedId = ref(null)
const detail = ref({})
const related = ref({ products: [], orders: [], reports: [] })
const filters = ref({ keyword: '', accountStatus: '', authStatus: '', role: '' })

const accountOptions = [
  { label: '全部状态', value: '' },
  { label: '正常', value: 'NORMAL' },
  { label: '禁用', value: 'DISABLED' },
  { label: '封禁', value: 'BANNED' }
]
const authOptions = [
  { label: '全部认证', value: '' },
  { label: '未认证', value: 'UNVERIFIED' },
  { label: '审核中', value: 'PENDING' },
  { label: '已认证', value: 'VERIFIED' },
  { label: '已驳回', value: 'REJECTED' }
]
const roleOptions = [
  { label: '全部角色', value: '' },
  { label: '普通用户', value: 'USER' },
  { label: '管理员', value: 'ADMIN' }
]
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))
const currentUserRole = computed(() => String(getUser()?.role || '').toUpperCase())
const isSuperAdmin = computed(() => currentUserRole.value === 'SUPER_ADMIN')

function setFilter(key, value) {
  filters.value[key] = value
  page.value = 1
  loadUsers()
}

function reload() {
  page.value = 1
  loadUsers()
}

async function loadUsers() {
  loading.value = true
  try {
    const query = { ...filters.value }
    if (!isSuperAdmin.value) {
      delete query.role
    }
    const result = await getAdminUsers({ ...query, page: page.value, pageSize })
    const rawUsers = result?.items || []
    users.value = rawUsers
    total.value = Number(result?.total || rawUsers.length)
  } catch (error) {
    uni.showToast({ title: error.message || '用户加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

function userName(item) {
  return item.nickname || item.realName || '微信用户'
}

function firstChar(item) {
  return userName(item).slice(0, 1)
}

function roleText(role) {
  return role === 'ADMIN' ? '管理员' : role === 'SUPER_ADMIN' ? '超级管理员' : '普通用户'
}

function authText(status) {
  return { UNVERIFIED: '未认证', PENDING: '审核中', VERIFIED: '已认证', REJECTED: '已驳回' }[status] || status || '未知'
}

function accountText(status) {
  return { NORMAL: '正常', DISABLED: '禁用', BANNED: '封禁', CANCELED: '已注销' }[status] || status || '未知'
}

function canRecoverUser(item) {
  if (!item) return false
  if (item.accountStatus === 'DISABLED') return true
  return false
}

function recoverButtonText(item) {
  return item?.accountStatus === 'BANNED' ? '解封' : '恢复'
}

function canChangeRole(item) {
  return isSuperAdmin.value && (item?.role === 'USER' || item?.role === 'ADMIN')
}

function nextRole(item) {
  return item?.role === 'ADMIN' ? 'USER' : 'ADMIN'
}

function roleActionText(item) {
  return nextRole(item) === 'ADMIN' ? '调整为管理员' : '调整为普通用户'
}

function roleTargetText(item) {
  return nextRole(item) === 'ADMIN' ? '管理员' : '普通用户'
}

function productStatusText(status) {
  return { ON_SALE: '在售', OFF_SHELF: '已下架', LOCKED: '已预约', SOLD: '已售出', DELETED: '已删除' }[status] || status || '未知'
}

function orderStatusText(status) {
  return { PENDING_CONFIRM: '待确认', WAIT_MEET: '待面交', COMPLETED: '已完成', CANCELED: '已取消', EXCEPTION_CLOSED: '异常关闭' }[status] || status || '未知'
}

function reportStatusText(status) {
  return { PENDING: '待处理', PROCESSING: '处理中', APPROVED: '已处理', REJECTED: '已驳回', CLOSED: '已关闭' }[status] || status || '未知'
}

function reportReasonText(reasonType) {
  return {
    FAKE_PRODUCT: '虚假或违规商品',
    INAPPROPRIATE_CONTENT: '不当内容',
    SCAM: '欺诈风险',
    TRADE_DISPUTE: '交易纠纷',
    FAKE: '虚假信息',
    FRAUD: '疑似诈骗',
    PROHIBITED: '违规商品',
    INAPPROPRIATE: '不当内容',
    HARASSMENT: '骚扰行为',
    OTHER: '其他原因'
  }[reasonType] || reasonType || '举报'
}

function openUserHome(item) {
  uni.navigateTo({ url: `/pages/user-profile/user-profile?id=${item.id}` })
}

function openProduct(item) {
  uni.navigateTo({ url: `/pages/detail/detail?id=${item.id}` })
}

function goStudentAuth(item) {
  uni.showToast({ title: `学生认证：${authText(item.authStatus)}`, icon: 'none' })
}

async function toggleRelated(item) {
  if (expandedId.value === item.id) {
    expandedId.value = null
    return
  }
  expandedId.value = item.id
  detail.value = {}
  related.value = { products: [], orders: [], reports: [] }
  try {
    const [userDetail, products, orders, reports] = await Promise.all([
      getAdminUserDetail(item.id),
      getAdminUserProducts(item.id),
      getAdminUserOrders(item.id),
      getAdminUserReports(item.id)
    ])
    detail.value = userDetail || {}
    related.value = {
      products: products?.items || [],
      orders: orders?.items || [],
      reports: reports?.items || []
    }
  } catch (error) {
    uni.showToast({ title: error.message || '关联信息加载失败', icon: 'none' })
  }
}

function changeStatus(item, status) {
  const actionText = accountText(status)
  uni.showModal({
    title: `${actionText}用户`,
    editable: true,
    placeholderText: '请输入处理原因',
    success: async (res) => {
      if (!res.confirm) return
      const reason = (res.content || '').trim()
      if (!reason) {
        uni.showToast({ title: '请填写处理原因', icon: 'none' })
        return
      }
      try {
        await updateAdminUserStatus(item.id, { accountStatus: status, reason })
        uni.showToast({ title: '操作成功', icon: 'success' })
        loadUsers()
      } catch (error) {
        uni.showToast({ title: error.message || '操作失败', icon: 'none' })
      }
    }
  })
}

function changeRole(item) {
  if (!canChangeRole(item)) return
  const targetRole = nextRole(item)
  const targetText = roleTargetText(item)
  uni.showModal({
    title: '调整角色',
    content: `确认将 ${userName(item)} 的角色调整为 ${targetText}？`,
    success: async (res) => {
      if (!res.confirm) return
      try {
        await updateAdminUserRole(item.id, { role: targetRole, reason: `超级管理员调整角色为${targetText}` })
        uni.showToast({ title: '角色调整成功', icon: 'success' })
        loadUsers()
      } catch (error) {
        uni.showToast({ title: error.message || '角色调整失败', icon: 'none' })
      }
    }
  })
}

function prevPage() {
  if (page.value <= 1) return
  page.value -= 1
  loadUsers()
}

function nextPage() {
  if (page.value >= totalPages.value) return
  page.value += 1
  loadUsers()
}

onShow(loadUsers)
</script>

<style scoped>
.page { min-height: 100vh; padding: 24rpx; background: #f5f6f8; box-sizing: border-box; }
.section-title { margin: 24rpx 0 18rpx; font-size: 34rpx; font-weight: 800; color: #1f2933; }
.toolbar { display: flex; gap: 16rpx; margin-bottom: 16rpx; }
.search { flex: 1; height: 76rpx; padding: 0 22rpx; border-radius: 14rpx; background: #fff; font-size: 26rpx; box-sizing: border-box; }
.search-btn { width: 140rpx; height: 76rpx; line-height: 76rpx; border-radius: 14rpx; background: #207f55; color: #fff; font-size: 26rpx; }
.filters { display: flex; flex-wrap: wrap; gap: 12rpx; margin-bottom: 14rpx; }
.chip { padding: 12rpx 20rpx; border-radius: 999rpx; background: #fff; color: #667085; font-size: 24rpx; }
.chip.active { background: #e7f6ee; color: #207f55; font-weight: 700; }
.card, .empty { padding: 24rpx; border-radius: 16rpx; background: #fff; margin-bottom: 18rpx; }
.row { display: flex; align-items: center; gap: 18rpx; }
.avatar { width: 72rpx; height: 72rpx; line-height: 72rpx; border-radius: 50%; background: #e7f6ee; color: #207f55; text-align: center; font-size: 30rpx; font-weight: 800; }
.main { flex: 1; min-width: 0; }
.name { font-size: 30rpx; font-weight: 800; color: #1f2933; }
.desc, .empty, .muted { margin-top: 8rpx; color: #667085; font-size: 24rpx; line-height: 1.45; }
.home-link { color: #207f55; font-size: 25rpx; font-weight: 700; }
.actions { display: flex; flex-wrap: wrap; gap: 12rpx; margin-top: 20rpx; }
.mini { margin: 0; padding: 0 18rpx; height: 56rpx; line-height: 56rpx; border-radius: 12rpx; background: #f2f4f7; color: #344054; font-size: 24rpx; }
.mini.warn { background: #fff7e6; color: #b66a00; }
.mini.danger { background: #fff1f0; color: #cf1322; }
.mini.ok { background: #e7f6ee; color: #207f55; }
.related { margin-top: 20rpx; padding-top: 18rpx; border-top: 1rpx solid #eef0f2; }
.stats { display: flex; flex-wrap: wrap; gap: 14rpx; color: #344054; font-size: 24rpx; }
.related-title { margin-top: 18rpx; color: #1f2933; font-size: 26rpx; font-weight: 800; }
.related-line { margin-top: 10rpx; padding: 12rpx 14rpx; border-radius: 10rpx; background: #f7f8fa; color: #475467; font-size: 24rpx; line-height: 1.45; }
.pager { display: flex; align-items: center; justify-content: center; gap: 20rpx; padding: 20rpx 0 40rpx; }
.page-btn { width: 150rpx; height: 60rpx; line-height: 60rpx; border-radius: 12rpx; background: #fff; color: #344054; font-size: 24rpx; }
.page-text { color: #667085; font-size: 24rpx; }
</style>
