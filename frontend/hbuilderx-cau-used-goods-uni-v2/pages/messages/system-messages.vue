<template>
  <view class="page">
    <view class="header">
      <view class="header-copy">
        <text class="title">系统消息</text>
        <text class="subtitle">订单进度、举报处理和平台通知</text>
      </view>
      <button class="read-all" :disabled="!unreadCount || marking" @click.stop="markAllRead">一键已读</button>
    </view>

    <view v-if="messages.length" class="list">
      <view v-for="item in messages" :key="item.id" class="swipe-wrap">
        <view class="delete-action" :class="{ disabled: deletingId === item.id }" @click.stop="remove(item)">删除</view>
        <view
          class="card"
          :class="{ swiped: swipedId === item.id }"
          @touchstart="touchStart($event, item.id)"
          @touchend="touchEnd"
          @click="open(item)"
        >
          <view class="card-main">
            <view class="icon" :class="{ read: item.read }">{{ item.read ? '✓' : '!' }}</view>
            <view class="body">
              <view class="head">
                <text class="message-title">{{ displayMessageTitle(item) }}</text>
                <text class="tag" :class="{ read: item.read }">{{ item.read ? '已读' : '未读' }}</text>
              </view>
              <text class="time">{{ item.createdAt }}</text>
              <text v-if="displayMessageContent(item)" class="content">{{ displayMessageContent(item) }}</text>
            </view>
          </view>
        </view>
      </view>
    </view>

    <EmptyState v-else title="暂无系统消息" detail="有新的交易进度或平台通知时会显示在这里" />
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import { tradeService } from '../../services/trade'
import { navigate, showError } from '../../utils/navigation'

const ORDER_MESSAGE_TYPES = ['ORDER_CREATED', 'ORDER_CONFIRMED', 'ORDER_CANCELED', 'ORDER_TIMEOUT', 'ORDER_EXCEPTION_CLOSED']
const SYSTEM_TYPES = [...ORDER_MESSAGE_TYPES, 'REPORT_HANDLED', 'SYSTEM_NOTICE']
const messages = ref([])
const marking = ref(false)
const deletingId = ref('')
const swipedId = ref('')
let startX = 0
let touchId = ''
const unreadCount = computed(() => messages.value.filter((item) => !item.read).length)

onShow(load)

async function load() {
  try {
    swipedId.value = ''
    const list = await tradeService.getMessages()
    messages.value = list.filter((item) => SYSTEM_TYPES.includes(item.type || item.messageType))
    updateBadge()
  } catch (error) {
    showError(error)
  }
}

function updateBadge() {
  const total = unreadCount.value
  if (total > 0) uni.setTabBarBadge({ index: 2, text: total > 99 ? '99+' : String(total) })
  else uni.removeTabBarBadge({ index: 2 })
}

async function open(item) {
  if (swipedId.value === item.id) {
    swipedId.value = ''
    return
  }
  if (isOrderProgressMessage(item)) {
    const orderId = orderTargetId(item)
    if (orderId) {
      await markOneRead(item)
      navigate('/pages/order/detail', { id: orderId, fromMessage: 1 })
      return
    }
  }
  navigate('/pages/interaction/message-detail', { id: item.id })
}

function messageTypeOf(item) {
  return item.type || item.messageType || ''
}

function isOrderProgressMessage(item) {
  return item.targetType === 'ORDER'
    || item.relatedType === 'ORDER'
    || ORDER_MESSAGE_TYPES.includes(messageTypeOf(item))
}

function orderTargetId(item) {
  const order = item.order || item.relatedOrder
  return order?.id || item.targetId || item.relatedId || item.orderId
}

async function markOneRead(item) {
  if (!item?.id || item.read) return
  try {
    await tradeService.markMessageRead(item.id)
    item.read = true
    updateBadge()
  } catch (error) {
    // Keep order navigation available even if read status update fails.
  }
}

function isStudentAuthResultMessage(item) {
  return item.title === '学生认证审核结果'
    || String(item.content || '').includes('学生认证已通过')
    || String(item.content || '').includes('学生认证未通过')
}

function displayMessageTitle(item) {
  if (!isStudentAuthResultMessage(item)) return item.title
  return String(item.content || '').includes('未通过')
    ? '学生认证审核未通过'
    : '学生认证审核通过'
}

function displayMessageContent(item) {
  if (isStudentAuthResultMessage(item)) return ''
  return item.content
}

function touchStart(event, id) {
  startX = event.changedTouches?.[0]?.clientX || 0
  touchId = id
  swipedId.value = swipedId.value === id ? '' : swipedId.value
}

function touchEnd(event) {
  const endX = event.changedTouches?.[0]?.clientX || 0
  const distance = endX - startX
  if (distance < -42) {
    swipedId.value = touchId
    return
  }
  if (distance > 20) swipedId.value = ''
}

async function markAllRead() {
  if (!unreadCount.value || marking.value) return
  marking.value = true
  try {
    const unread = messages.value.filter((item) => !item.read)
    await Promise.all(unread.map((item) => tradeService.markMessageRead(item.id).catch(() => null)))
    messages.value = messages.value.map((item) => ({ ...item, read: true }))
    updateBadge()
    uni.showToast({ title: '已全部标为已读', icon: 'success' })
  } catch (error) {
    showError(error)
  } finally {
    marking.value = false
  }
}

function remove(item) {
  if (deletingId.value === item.id) return
  uni.showModal({
    title: '删除系统消息',
    content: '确认删除这条系统消息吗？',
    confirmColor: '#d92d20',
    success: async ({ confirm }) => {
      if (!confirm || deletingId.value) return
      deletingId.value = item.id
      try {
        await tradeService.deleteMessage(item.id)
        messages.value = messages.value.filter((message) => String(message.id) !== String(item.id))
        swipedId.value = ''
        updateBadge()
        uni.showToast({ title: '已删除', icon: 'success' })
      } catch (error) {
        showError(error)
      } finally {
        deletingId.value = ''
        swipedId.value = ''
      }
    }
  })
}
</script>

<style scoped>
.page { min-height: 100vh; padding: 30rpx 28rpx 48rpx; background: #f5f8f6; box-sizing: border-box; }
.header { display: flex; align-items: center; justify-content: space-between; gap: 20rpx; margin-bottom: 24rpx; padding: 28rpx 26rpx; border-radius: 28rpx; background: linear-gradient(135deg, #23734f, #3e9b72); box-shadow: 0 12rpx 32rpx rgba(35, 115, 79, .18); }
.header-copy { flex: 1; min-width: 0; }
.title, .subtitle { display: block; }
.title { color: #fff; font-size: 42rpx; font-weight: 800; }
.subtitle { margin-top: 10rpx; color: rgba(255,255,255,.78); font-size: 24rpx; }
.read-all { flex-shrink: 0; min-width: 136rpx; height: 56rpx; margin: 0 0 0 auto; padding: 0 18rpx; border-radius: 999rpx; background: rgba(255,255,255,.94); color: #23734f; font-size: 24rpx; line-height: 56rpx; }
.read-all[disabled], .delete-action.disabled { opacity: .56; }
.list { display: flex; flex-direction: column; gap: 22rpx; }
.swipe-wrap { position: relative; overflow: hidden; border-radius: 28rpx; }
.delete-action { position: absolute; top: 0; right: 0; bottom: 0; width: 136rpx; display: flex; align-items: center; justify-content: center; background: #f04444; color: #fff; font-size: 28rpx; }
.card { position: relative; z-index: 1; display: flex; align-items: stretch; overflow: hidden; border-radius: 28rpx; background: #fff; box-shadow: 0 10rpx 30rpx rgba(28, 68, 52, .06); transition: transform .18s ease; }
.card.swiped { transform: translateX(-136rpx); }
.card-main { display: flex; flex: 1; min-width: 0; gap: 22rpx; padding: 30rpx 20rpx 30rpx 28rpx; }
.icon { display: flex; width: 54rpx; height: 54rpx; flex: 0 0 54rpx; align-items: center; justify-content: center; border-radius: 50%; background: #fff1e8; color: #e36a3e; font-size: 26rpx; font-weight: 800; }
.icon.read { background: #eef4f1; color: #7c8a84; }
.body { flex: 1; min-width: 0; }
.head { display: flex; align-items: center; justify-content: space-between; gap: 18rpx; }
.message-title { flex: 1; min-width: 0; overflow: hidden; color: #26342f; font-size: 31rpx; font-weight: 800; text-overflow: ellipsis; white-space: nowrap; }
.tag { flex-shrink: 0; padding: 6rpx 14rpx; border-radius: 999rpx; background: #fff3dd; color: #bd7a16; font-size: 21rpx; }
.tag.read { background: #eef2f0; color: #8b9691; }
.time { display: block; margin-top: 8rpx; color: #9ba5a0; font-size: 22rpx; }
.content { display: block; margin-top: 18rpx; color: #5c6862; font-size: 26rpx; line-height: 1.7; }
</style>
