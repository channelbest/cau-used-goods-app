<template>
  <view class="page">
    <view class="top">
      <view>
        <text class="title">消息</text>
        <text class="subtitle">系统通知和用户私信都会在这里</text>
      </view>
    </view>

    <view v-if="systemEntry" class="conversation pinned" @click="openSystemMessages">
      <view class="avatar system-avatar">系</view>
      <view class="body">
        <view class="head">
          <text class="name">系统消息</text>
          <text class="time">{{ systemEntry.time }}</text>
        </view>
        <text class="preview">{{ systemEntry.preview }}</text>
      </view>
      <view v-if="systemEntry.unread" class="badge">{{ systemEntry.unread > 99 ? '99+' : systemEntry.unread }}</view>
    </view>

    <view v-if="visibleConversations.length" class="list">
      <view v-for="item in visibleConversations" :key="item.id" class="swipe-wrap">
        <view class="delete-action" @click="hideConversation(item.id)">删除</view>
        <view
          class="conversation user-conversation"
          :class="{ swiped: swipedId === item.id }"
          @touchstart="touchStart($event, item.id)"
          @touchend="touchEnd"
          @click="openChat(item)"
        >
          <view v-if="item.targetBanned" class="avatar banned-avatar">禁</view>
          <image v-else-if="item.avatar" class="avatar image-avatar" :src="item.avatar" mode="aspectFill" />
          <view class="body">
            <view class="head">
              <text class="name">{{ conversationDisplayName(item) }}</text>
              <text class="time">{{ formatTime(item.lastMessageTime || item.updateTime) }}</text>
            </view>
            <text class="product-line">商品：{{ conversationProductTitle(item) }}</text>
            <text class="preview">{{ conversationPreview(item) }}</text>
          </view>
          <view v-if="item.unreadCount" class="badge">{{ item.unreadCount > 99 ? '99+' : item.unreadCount }}</view>
        </view>
      </view>
    </view>

    <view v-if="!loading && !systemEntry && !visibleConversations.length" class="empty-wrap">
      <EmptyState title="暂无消息" detail="订单进度、系统通知和聊天消息会显示在这里" />
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import EmptyState from '../../components/EmptyState.vue'
import { hideConversation as hideConversationApi, listConversations } from '../../api/chat'
import { getPublicProfile } from '../../api/user'
import { tradeService } from '../../services/trade'
import { getUser } from '../../utils/auth'
import { BASE_URL } from '../../utils/request'
import { navigate, showError } from '../../utils/navigation'
import { accountStatusOf, isBannedUserStatus } from '../../utils/user-format'

const SYSTEM_TYPES = ['ORDER_CREATED', 'ORDER_CONFIRMED', 'ORDER_CANCELED', 'ORDER_TIMEOUT', 'REPORT_HANDLED', 'SYSTEM_NOTICE']
const loading = ref(false)
const messages = ref([])
const conversations = ref([])
const swipedId = ref('')
let startX = 0
let touchId = ''

const normalizeImage = (url) => {
  if (!url) return ''
  return /^https?:\/\//.test(url) ? url : `${BASE_URL}${url}`
}

function localizeHttpImage(url) {
  if (!/^http:\/\//.test(url || '')) return Promise.resolve(url || '')
  return new Promise((resolve) => {
    uni.downloadFile({
      url,
      success: (res) => resolve(res.tempFilePath || url),
      fail: () => resolve(url)
    })
  })
}

const pick = (...values) => values.find((value) => value !== undefined && value !== null && value !== '') || ''
const sameId = (a, b) => a !== undefined && a !== null && a !== '' && b !== undefined && b !== null && b !== '' && String(a) === String(b)

function currentUserId() {
  const user = getUser() || {}
  return pick(user.id, user.userId, user.user_id)
}

function getBuyerId(item) {
  return pick(item.buyerId, item.buyerUserId, item.buyer_id, item.buyer?.id, item.buyer?.userId)
}

function getSellerId(item) {
  return pick(item.sellerId, item.sellerUserId, item.seller_id, item.seller?.id, item.seller?.userId, item.product?.sellerId, item.product?.seller?.id)
}

function getPeerId(item) {
  const mine = currentUserId()
  const buyerId = getBuyerId(item)
  const sellerId = getSellerId(item)
  if (sameId(mine, buyerId) && sellerId) return sellerId
  if (sameId(mine, sellerId) && buyerId) return buyerId
  return pick(
    item.targetUserId,
    item.targetId,
    item.otherUserId,
    item.userId,
    item.targetUser?.id,
    item.otherUser?.id,
    item.user?.id,
    sellerId,
    buyerId
  )
}

function getPeerNickname(item) {
  const peerId = getPeerId(item)
  const buyerId = getBuyerId(item)
  const sellerId = getSellerId(item)
  if (sameId(peerId, buyerId)) {
    return pick(item.buyerNickname, item.buyerName, item.buyerUserNickname, item.buyer?.nickname, item.buyer?.name)
  }
  if (sameId(peerId, sellerId)) {
    return pick(item.sellerNickname, item.sellerName, item.sellerUserNickname, item.seller?.nickname, item.seller?.name, item.product?.seller?.nickname)
  }
  return ''
}

function getPeerAvatar(item) {
  const peerId = getPeerId(item)
  const buyerId = getBuyerId(item)
  const sellerId = getSellerId(item)
  if (sameId(peerId, buyerId)) {
    return normalizeImage(pick(item.buyerAvatarUrl, item.buyerAvatar, item.buyerUserAvatarUrl, item.buyer?.avatarUrl, item.buyer?.avatar, item.buyer?.avatar_url))
  }
  if (sameId(peerId, sellerId)) {
    return normalizeImage(pick(item.sellerAvatarUrl, item.sellerAvatar, item.sellerUserAvatarUrl, item.seller?.avatarUrl, item.seller?.avatar, item.seller?.avatar_url, item.product?.seller?.avatarUrl, item.product?.seller?.avatar))
  }
  return ''
}

function getPeerStatus(item) {
  const peerId = getPeerId(item)
  const buyerId = getBuyerId(item)
  const sellerId = getSellerId(item)
  if (sameId(peerId, buyerId)) {
    return pick(item.buyerAccountStatus, item.buyerStatus, item.buyer_account_status, item.buyer?.accountStatus, item.buyer?.account_status, item.buyer?.status)
  }
  if (sameId(peerId, sellerId)) {
    return pick(item.sellerAccountStatus, item.sellerStatus, item.seller_account_status, item.seller?.accountStatus, item.seller?.account_status, item.seller?.status, item.product?.seller?.accountStatus, item.product?.seller?.account_status, item.product?.seller?.status)
  }
  return pick(
    item.targetAccountStatus,
    item.targetUserStatus,
    item.otherUserStatus,
    item.userStatus,
    item.targetUser?.accountStatus,
    item.targetUser?.account_status,
    item.targetUser?.status,
    item.otherUser?.accountStatus,
    item.otherUser?.account_status,
    item.otherUser?.status,
    item.user?.accountStatus,
    item.user?.account_status,
    item.user?.status
  )
}

function getConversationAvatar(item) {
  return getPeerAvatar(item)
}

function getConversationNickname(item) {
  return getPeerNickname(item) || pick(
    item.targetNickname,
    item.targetUserNickname,
    item.otherUserNickname,
    item.nickname,
    item.targetUser?.nickname,
    item.otherUser?.nickname,
    item.user?.nickname
  )
}

function getConversationTargetId(item) {
  return getPeerId(item)
}

const systemMessages = computed(() => messages.value.filter((item) => SYSTEM_TYPES.includes(item.type || item.messageType)))
const systemEntry = computed(() => {
  if (!systemMessages.value.length) return null
  const latest = systemMessages.value[0]
  const unread = systemMessages.value.filter((item) => !item.read).length
  return {
    unread,
    time: formatTime(latest.createdAt || latest.createTime),
    preview: unread ? `${unread} 条未读系统通知` : systemMessagePreview(latest)
  }
})
const visibleConversations = computed(() => conversations.value
  .slice()
  .filter((item) => String(item.buyerId || '') !== String(item.sellerId || '') || !item.buyerId || !item.sellerId)
  .sort((a, b) => timeValue(b.lastMessageTime || b.updateTime) - timeValue(a.lastMessageTime || a.updateTime)))
const unreadTotal = computed(() => {
  const systemUnread = systemEntry.value?.unread || 0
  const chatUnread = visibleConversations.value.reduce((sum, item) => sum + Number(item.unreadCount || 0), 0)
  return systemUnread + chatUnread
})

function timeValue(value) {
  const date = new Date(String(value || '').replace(/-/g, '/'))
  return Number.isNaN(date.getTime()) ? 0 : date.getTime()
}

function formatTime(value) {
  if (!value) return ''
  const date = new Date(String(value).replace(/-/g, '/'))
  if (Number.isNaN(date.getTime())) return String(value).slice(0, 16)
  const now = new Date()
  const pad = (n) => String(n).padStart(2, '0')
  const sameDay = date.toDateString() === now.toDateString()
  return sameDay ? `${pad(date.getHours())}:${pad(date.getMinutes())}` : `${pad(date.getMonth() + 1)}-${pad(date.getDate())}`
}

function conversationProductTitle(item) {
  return item.productTitle || item.productTitleSnapshot || item.product?.title || '商品'
}

function conversationProductId(item) {
  return item.productId || item.product_id || item.product?.id || ''
}

function getDeletedMessageIds(conversationId) {
  const value = uni.getStorageSync(`deleted-chat-messages-${conversationId}`) || []
  return Array.isArray(value) ? value.map(String) : []
}

function getDeletedMessageContents(conversationId) {
  const value = uni.getStorageSync(`deleted-chat-message-contents-${conversationId}`) || []
  return Array.isArray(value) ? value.map(String) : []
}

function conversationLastMessageId(item) {
  return item.lastMessageId || item.last_message_id || item.lastMessage?.id || item.message?.id || ''
}

function isConversationPreviewDeleted(item) {
  const lastMessageId = conversationLastMessageId(item)
  const lastContent = String(item.lastMessageContent || '').trim()
  const deletedIds = getDeletedMessageIds(item.id)
  const deletedContents = getDeletedMessageContents(item.id)
  return (lastMessageId && deletedIds.includes(String(lastMessageId)))
    || (lastContent && deletedContents.includes(lastContent))
}

function conversationPreview(item) {
  if (isConversationPreviewDeleted(item)) return '消息已删除'
  return item.lastMessageContent || `关于「${item.productTitle || '商品'}」的沟通`
}

function isStudentAuthResultMessage(item) {
  return item.title === '学生认证审核结果'
    || String(item.content || '').includes('学生认证已通过')
    || String(item.content || '').includes('学生认证未通过')
}

function systemMessagePreview(item) {
  if (!isStudentAuthResultMessage(item)) return item.content || item.title || '暂无新的系统通知'
  return String(item.content || '').includes('未通过')
    ? '学生认证审核未通过'
    : '学生认证审核通过'
}

function conversationDisplayName(item) {
  const name = item.targetNickname || 'CAU 同学'
  return item.targetBanned ? `${name}——该用户已被封禁` : name
}

function isBannedProfileError(error) {
  const message = String(error?.message || '')
  return message.includes('封禁')
    || message.includes('禁用')
    || message.includes('无法查找')
    || message.includes('没有权限')
    || message.includes('用户不存在')
    || message.includes('不存在')
    || message.includes('未找到相关数据')
    || message.toLowerCase().includes('forbidden')
    || message.toLowerCase().includes('permission')
}

async function load() {
  loading.value = true
  swipedId.value = ''
  uni.removeStorageSync('hidden-chat-conversations')
  try {
    const [systemList, chatResult] = await Promise.all([
      tradeService.getMessages().catch(() => []),
      listConversations({ page: 1, pageSize: 50 }).catch(() => ({ items: [] }))
    ])
    messages.value = systemList
    conversations.value = (chatResult.items || []).map((item) => ({
      ...item,
      id: String(item.id),
      targetNickname: getConversationNickname(item),
      targetUserId: getConversationTargetId(item),
      avatar: getConversationAvatar(item),
      targetBanned: isBannedUserStatus(getPeerStatus(item))
    }))
    conversations.value = await Promise.all(conversations.value.map(async (item) => {
      const targetUserId = item.targetUserId || getConversationTargetId(item)
      if (!targetUserId) return item
      const profile = await getPublicProfile(targetUserId).catch((error) => {
        if (isBannedProfileError(error)) return { __bannedProfile: true }
        return null
      })
      if (!profile) return item
      if (profile.__bannedProfile) {
        return {
          ...item,
          targetUserId,
          targetNickname: item.targetNickname || '该用户已封禁',
          avatar: '',
          targetBanned: true
        }
      }
      const targetBanned = isBannedUserStatus(accountStatusOf(profile?.user || profile))
      const avatar = await localizeHttpImage(normalizeImage(pick(
        profile.avatarUrl,
        profile.avatar,
        profile.avatar_url,
        profile.user?.avatarUrl,
        profile.user?.avatar,
        profile.user?.avatar_url
      )))
      return {
        ...item,
        targetUserId,
        targetNickname: targetBanned ? (profile.nickname || item.targetNickname || '该用户已封禁') : (profile.nickname || item.targetNickname),
        avatar: targetBanned ? '' : avatar,
        targetBanned
      }
    }))
    updateTabBadge()
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
}

function updateTabBadge() {
  const text = unreadTotal.value > 99 ? '99+' : String(unreadTotal.value)
  if (unreadTotal.value > 0) {
    uni.setTabBarBadge({ index: 2, text })
  } else {
    uni.removeTabBarBadge({ index: 2 })
  }
}

function openSystemMessages() {
  navigate('/pages/messages/system-messages')
}

async function openChat(item) {
  const targetUserId = item.targetUserId || getConversationTargetId(item)
  let targetNickname = item.targetNickname || ''
  let targetAvatar = item.avatar || ''
  let targetBanned = item.targetBanned
  if (targetUserId && (!targetNickname || !targetAvatar)) {
    const profile = await getPublicProfile(targetUserId).catch((error) => {
      if (isBannedProfileError(error)) return { __bannedProfile: true }
      return null
    })
    if (profile) {
      if (profile.__bannedProfile) {
        targetBanned = true
        targetNickname = targetNickname || '该用户已封禁'
        targetAvatar = ''
      } else {
      targetBanned = isBannedUserStatus(accountStatusOf(profile?.user || profile))
      targetNickname = targetBanned ? (profile.nickname || targetNickname || '该用户已封禁') : (profile.nickname || targetNickname)
      targetAvatar = targetBanned ? '' : normalizeImage(pick(
        profile.avatarUrl,
        profile.avatar,
        profile.avatar_url,
        profile.user?.avatarUrl,
        profile.user?.avatar,
        profile.user?.avatar_url,
        targetAvatar
      ))
      }
    }
  }
  navigate('/pages/chat/chat', {
    conversationId: item.id,
    title: conversationProductTitle(item) || '私信沟通',
    targetUserId,
    targetNickname: targetBanned ? '！该用户已被封禁，无法查找' : targetNickname,
    targetAvatar: !targetBanned ? targetAvatar : '',
    targetBanned: targetBanned ? 1 : 0,
    productId: conversationProductId(item)
  })
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

async function hideConversation(id) {
  try {
    await hideConversationApi(id)
    conversations.value = conversations.value.filter((item) => String(item.id) !== String(id))
    uni.showToast({ title: '会话已删除', icon: 'success' })
  } catch (error) {
    showError(error)
  } finally {
    swipedId.value = ''
    updateTabBadge()
  }
}

onShow(load)
</script>

<style scoped>
.page { min-height: 100vh; padding: 24rpx 24rpx 44rpx; background: #f5f6f7; box-sizing: border-box; }
.top { display: flex; align-items: center; justify-content: space-between; padding: 18rpx 6rpx 24rpx; }
.title, .subtitle { display: block; }
.title { color: #202124; font-size: 42rpx; font-weight: 800; }
.subtitle { margin-top: 8rpx; color: #8a8f94; font-size: 24rpx; }
.badge { display: flex; align-items: center; justify-content: center; border-radius: 999rpx; background: #f04444; color: #fff; font-size: 21rpx; }
.badge { min-width: 34rpx; height: 34rpx; padding: 0 8rpx; flex-shrink: 0; }
.list { display: flex; flex-direction: column; gap: 16rpx; margin-top: 16rpx; }
.swipe-wrap { position: relative; overflow: hidden; border-radius: 24rpx; }
.delete-action { position: absolute; top: 0; right: 0; bottom: 0; width: 136rpx; display: flex; align-items: center; justify-content: center; background: #f04444; color: #fff; font-size: 28rpx; }
.conversation { position: relative; display: flex; align-items: center; gap: 20rpx; min-height: 116rpx; padding: 24rpx; border-radius: 24rpx; background: #fff; box-shadow: 0 8rpx 28rpx rgba(23, 33, 43, .04); transition: transform .18s ease; box-sizing: border-box; }
.conversation.swiped { transform: translateX(-136rpx); }
.pinned { background: #eef0f2; box-shadow: none; }
.user-conversation { z-index: 1; }
.avatar { display: flex; width: 84rpx; height: 84rpx; flex-shrink: 0; align-items: center; justify-content: center; border-radius: 50%; color: #fff; font-size: 28rpx; font-weight: 700; }
.system-avatar { background: #87909a; }
.banned-avatar { background: #d92d20; color: #fff; font-size: 34rpx; }
.image-avatar { background: #e8ecef; }
.body { flex: 1; min-width: 0; }
.head { display: flex; align-items: center; justify-content: space-between; gap: 16rpx; }
.name { overflow: hidden; color: #222; font-size: 31rpx; font-weight: 700; text-overflow: ellipsis; white-space: nowrap; }
.time { flex-shrink: 0; color: #a0a6ad; font-size: 22rpx; }
.preview { display: block; overflow: hidden; margin-top: 12rpx; color: #7b8289; font-size: 25rpx; line-height: 34rpx; text-overflow: ellipsis; white-space: nowrap; }
.product-line { display: block; overflow: hidden; margin-top: 8rpx; color: #a0a6ad; font-size: 23rpx; line-height: 32rpx; text-overflow: ellipsis; white-space: nowrap; }
.empty-wrap { margin-top: 80rpx; border-radius: 28rpx; background: #fff; overflow: hidden; }
</style>
