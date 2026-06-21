<template>
  <view class="page">
    <view class="chat-head">
      <view class="chat-title clickable-title" @click="openSellerProfile">
        <text>{{ sellerName }}</text>
        <text class="profile-link">{{ isTargetUnavailable && !isCurrentUserRestricted ? '已禁用' : '查看主页' }}</text>
      </view>
      <view class="chat-subtitle">请在平台内沟通交易细节，注意保护个人隐私</view>
    </view>

    <view v-if="productId" class="product-link" @click="openProduct">
      <view class="product-link-body">
        <text class="product-link-label">商品</text>
        <text class="product-link-title">{{ productTitle }}</text>
      </view>
      <text class="product-link-arrow">›</text>
    </view>

    <view v-if="isCurrentUserRestricted" class="restriction-banner">
      <text class="restriction-text">{{ currentUserRestrictionText }}</text>
    </view>

    <scroll-view scroll-y class="messages" :class="{ 'has-product': productId, 'has-restriction': isCurrentUserRestricted }" :scroll-into-view="lastMessageId">
      <view v-for="item in displayMessages" :id="`msg-${item.id}`" :key="item.id">
        <view v-if="item.showTime" class="time-divider">{{ item.timeText }}</view>
        <view class="message-swipe" :class="{ active: swipedMessageId === item.id }">
          <view class="message-delete" @click.stop="removeMessage(item)">删除</view>
          <view
            class="message-front"
            :class="{ swiped: swipedMessageId === item.id }"
            @touchstart="messageTouchStart($event, item.id)"
            @touchend="messageTouchEnd"
          >
            <view class="message-row" :class="{ mine: item.mine }">
              <image
                v-if="!item.mine && messageAvatar(item)"
                class="avatar image-avatar"
                :src="messageAvatar(item)"
                mode="aspectFill"
                @click.stop="openUser(item.senderId)"
              />
              <view
                v-else-if="!item.mine && isTargetUnavailable"
                class="avatar banned-avatar"
                @click.stop="openUser(item.senderId)"
              >
                禁
              </view>
              <view class="bubble">
                <view class="content">{{ item.content }}</view>
              </view>
              <image
                v-if="item.mine && messageAvatar(item)"
                class="avatar image-avatar"
                :src="messageAvatar(item)"
                mode="aspectFill"
                @click.stop="openUser(currentUserId)"
              />
            </view>
          </view>
        </view>
      </view>
      <view v-if="!messages.length && !loading" class="empty">还没有消息，先打个招呼吧</view>
    </scroll-view>

    <view class="composer" :class="{ restricted: isCurrentUserRestricted }">
      <input v-model.trim="draft" class="input" maxlength="500" confirm-type="send" :placeholder="isCurrentUserRestricted ? currentUserRestrictionText : '输入消息'" :disabled="isCurrentUserRestricted" @confirm="send" />
      <button class="send" :disabled="!draft || sending || isCurrentUserRestricted" @click="send">发送</button>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad, onPullDownRefresh } from '@dcloudio/uni-app'
import { deleteMessage, getConversationProduct, listMessages, markConversationRead, sendMessage } from '../../api/chat'
import { getPublicProfile } from '../../api/user'
import { getUser } from '../../utils/auth'
import { BASE_URL } from '../../utils/request'
import { navigate } from '../../utils/navigation'
import { accountStatusOf, isBannedUserStatus, isDisabledUserStatus } from '../../utils/user-format'

const conversationId = ref('')
const title = ref('')
const targetUserId = ref('')
const targetNicknameSnapshot = ref('')
const targetAvatarSnapshot = ref('')
const productId = ref('')
const messages = ref([])
const draft = ref('')
const loading = ref(false)
const sending = ref(false)
const targetProfile = ref(null)
const targetBannedSnapshot = ref(false)
const mineProfile = ref(null)
const profileMap = ref({})
const swipedMessageId = ref('')
let touchStartX = 0
let touchMessageId = ''

const currentUserId = computed(() => getUser()?.id || getUser()?.userId || '')
const currentUser = computed(() => getUser() || {})
const currentUserAccountStatus = computed(() => accountStatusOf(currentUser.value))
const isCurrentUserRestricted = computed(() => isBannedUserStatus(currentUserAccountStatus.value) || isDisabledUserStatus(currentUserAccountStatus.value))
const currentUserRestrictionText = computed(() => {
  if (isBannedUserStatus(currentUserAccountStatus.value)) return '账号已被封禁，无法发送消息'
  if (isDisabledUserStatus(currentUserAccountStatus.value)) return '账号已被禁用，无法发送消息'
  return ''
})
const productTitle = computed(() => title.value || '商品详情')
const pick = (...values) => values.find((value) => value !== undefined && value !== null && value !== '') || ''
const BANNED_USER_TEXT = '！该用户已被封禁，无法查找'
const sellerIdText = computed(() => pick(targetUserId.value, targetProfile.value?.id, targetProfile.value?.userId))
const targetAccountStatus = computed(() => accountStatusOf(targetProfile.value?.user || targetProfile.value || {}))
const isTargetUnavailable = computed(() => (
  isBannedUserStatus(targetAccountStatus.value)
  || targetBannedSnapshot.value
))
const sellerName = computed(() => {
  const raw = targetProfile.value?.nickname || targetNicknameSnapshot.value || '对方'
  return raw
})
const targetAvatar = computed(() => normalizeImage(pick(targetProfile.value?.avatarUrl, targetProfile.value?.avatar, targetProfile.value?.avatar_url, targetAvatarSnapshot.value)))
const mineAvatar = computed(() => normalizeImage(pick(currentUser.value.avatarUrl, currentUser.value.avatar, mineProfile.value?.avatarUrl, mineProfile.value?.avatar, mineProfile.value?.avatar_url)))
const mineName = computed(() => currentUser.value.nickname || mineProfile.value?.nickname || '我')
const lastMessageId = computed(() => {
  const last = messages.value[messages.value.length - 1]
  return last ? `msg-${last.id}` : ''
})

const deletedMessageKey = () => `deleted-chat-messages-${conversationId.value}`
const deletedMessageContentKey = () => `deleted-chat-message-contents-${conversationId.value}`

const getDeletedMessageIds = () => {
  const value = uni.getStorageSync(deletedMessageKey()) || []
  return Array.isArray(value) ? value.map(String) : []
}

const saveDeletedMessageId = (id) => {
  const value = String(id)
  const ids = getDeletedMessageIds()
  if (!ids.includes(value)) {
    uni.setStorageSync(deletedMessageKey(), ids.concat(value))
  }
}

const getDeletedMessageContents = () => {
  const value = uni.getStorageSync(deletedMessageContentKey()) || []
  return Array.isArray(value) ? value.map(String) : []
}

const saveDeletedMessageContent = (content) => {
  const value = String(content || '').trim()
  if (!value) return
  const contents = getDeletedMessageContents()
  if (!contents.includes(value)) {
    uni.setStorageSync(deletedMessageContentKey(), contents.concat(value))
  }
}

const visibleMessageItems = (items = []) => {
  const deletedIds = getDeletedMessageIds()
  return items.filter((item) => !deletedIds.includes(String(item.id)))
}

function normalizeImage(url) {
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

function setProfile(id, profile) {
  if (!id || !profile) return
  profileMap.value = {
    ...profileMap.value,
    [String(id)]: profile
  }
}

function seedInitialProfiles() {
  const mineAvatarUrl = normalizeImage(pick(currentUser.value.avatarUrl, currentUser.value.avatar))
  if (currentUserId.value && (mineAvatarUrl || currentUser.value.nickname)) {
    setProfile(currentUserId.value, { ...currentUser.value, avatarUrl: mineAvatarUrl })
  }
  if (targetUserId.value && (targetAvatarSnapshot.value || targetNicknameSnapshot.value)) {
    const profile = {
      id: targetUserId.value,
      nickname: targetNicknameSnapshot.value,
      avatarUrl: normalizeImage(targetAvatarSnapshot.value)
    }
    targetProfile.value = profile
    setProfile(targetUserId.value, profile)
  }
}

function profileAvatar(id) {
  const profile = profileMap.value[String(id || '')]
  return normalizeImage(pick(profile?.avatarUrl, profile?.avatar, profile?.avatar_url))
}

function messageAvatar(item) {
  if (item.mine) return profileAvatar(currentUserId.value) || mineAvatar.value
  return profileAvatar(item.senderId) || (String(item.senderId) === String(targetUserId.value) ? targetAvatar.value : '')
}

const messageTime = (item) => new Date(String(item.createTime || '').replace(/-/g, '/')).getTime()
const formatTime = (value) => {
  const date = new Date(String(value || '').replace(/-/g, '/'))
  if (Number.isNaN(date.getTime())) return value || ''
  const pad = (n) => String(n).padStart(2, '0')
  return `${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`
}

const displayMessages = computed(() => {
  let previousTime = 0
  return messages.value.map((item, index) => {
    const currentTime = messageTime(item)
    const showTime = index === 0 || !previousTime || currentTime - previousTime > 5 * 60 * 1000
    if (currentTime) previousTime = currentTime
    return {
      ...item,
      mine: Number(item.senderId) === Number(currentUserId.value),
      showTime,
      timeText: formatTime(item.createTime)
    }
  })
})

const load = async () => {
  if (!conversationId.value) return
  loading.value = true
  swipedMessageId.value = ''
  try {
    const result = await listMessages(conversationId.value, { page: 1, pageSize: 50 })
    messages.value = visibleMessageItems(result.items || [])
    if (!targetUserId.value) {
      const otherMessage = messages.value.find((item) => Number(item.senderId) !== Number(currentUserId.value))
      targetUserId.value = otherMessage?.senderId || ''
    }
    await markConversationRead(conversationId.value)
  } catch (error) {
    uni.showToast({ title: error.message || '消息加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

const loadProfiles = async () => {
  const ids = new Set([currentUserId.value, targetUserId.value])
  messages.value.forEach((item) => {
    if (item.senderId) ids.add(item.senderId)
  })
  const jobs = Array.from(ids).filter(Boolean).map((id) => (
    getPublicProfile(id).then(async (data) => {
      const avatarUrl = await localizeHttpImage(normalizeImage(pick(data?.avatarUrl, data?.avatar, data?.avatar_url)))
      const profile = { ...data, avatarUrl }
      setProfile(id, profile)
      if (String(id) === String(targetUserId.value)) {
        targetProfile.value = profile
        targetBannedSnapshot.value = isBannedUserStatus(accountStatusOf(profile?.user || profile))
      }
      if (String(id) === String(currentUserId.value)) mineProfile.value = profile
    }).catch(() => {})
  ))
  await Promise.all(jobs)
}

async function loadInitialTargetProfile() {
  if (!targetUserId.value || (targetNicknameSnapshot.value && targetAvatarSnapshot.value)) return
  const data = await getPublicProfile(targetUserId.value).catch(() => null)
  if (!data) return
  const avatarUrl = await localizeHttpImage(normalizeImage(pick(data?.avatarUrl, data?.avatar, data?.avatar_url)))
  const profile = { ...data, avatarUrl }
  const targetBanned = isBannedUserStatus(accountStatusOf(profile?.user || profile))
  targetBannedSnapshot.value = targetBanned
  targetNicknameSnapshot.value = targetBanned
    ? BANNED_USER_TEXT
    : profile.nickname || targetNicknameSnapshot.value
  targetAvatarSnapshot.value = avatarUrl || targetAvatarSnapshot.value
  targetProfile.value = profile
  setProfile(targetUserId.value, profile)
}

const send = async () => {
  if (!draft.value || sending.value) return
  if (isCurrentUserRestricted.value) {
    uni.showToast({ title: currentUserRestrictionText.value, icon: 'none' })
    return
  }
  const original = draft.value
  sending.value = true
  try {
    const message = await sendMessage(conversationId.value, original)
    messages.value = messages.value.concat(message)
    await loadProfiles()
    draft.value = ''
    if (message.content !== original) {
      uni.showToast({ title: '禁止使用敏感词汇，已替换为 *', icon: 'none' })
    }
  } catch (error) {
    uni.showToast({ title: error.message || '发送失败', icon: 'none' })
  } finally {
    sending.value = false
  }
}

const messageTouchStart = (event, id) => {
  touchStartX = event.changedTouches?.[0]?.clientX || 0
  touchMessageId = id
  swipedMessageId.value = swipedMessageId.value === id ? '' : swipedMessageId.value
}

const messageTouchEnd = (event) => {
  const endX = event.changedTouches?.[0]?.clientX || 0
  const distance = endX - touchStartX
  if (distance < -42) {
    swipedMessageId.value = touchMessageId
    return
  }
  if (distance > 20) swipedMessageId.value = ''
}

const removeMessage = async (item) => {
  const id = item?.id
  if (!id) return
  saveDeletedMessageId(id)
  saveDeletedMessageContent(item.content)
  messages.value = messages.value.filter((item) => String(item.id) !== String(id))
  swipedMessageId.value = ''
  try {
    await deleteMessage(id)
  } catch (error) {
    // Older backend builds do not expose per-message deletion; local deletion is kept.
  }
}

function openUser(id) {
  if (isCurrentUserRestricted.value) {
    uni.showModal({
      title: '提示',
      content: '你已被封禁，无法查看用户主页',
      showCancel: false
    })
    return
  }
  if (isTargetUnavailable.value && String(id || targetUserId.value) === String(targetUserId.value)) {
    uni.showModal({
      title: '提示',
      content: '该用户已被永久禁用',
      showCancel: false
    })
    return
  }
  const profileId = id || targetUserId.value
  if (!profileId) return
  navigate('/pages/user-profile/user-profile', { id: profileId })
}

function openSellerProfile() {
  if (isCurrentUserRestricted.value) {
    uni.showModal({
      title: '提示',
      content: '你已被封禁，无法查看用户主页',
      showCancel: false
    })
    return
  }
  if (isTargetUnavailable.value) {
    uni.showModal({
      title: '提示',
      content: '该用户已被永久禁用',
      showCancel: false
    })
    return
  }
  const profileId = sellerIdText.value || targetUserId.value
  if (!profileId) return
  navigate('/pages/user-profile/user-profile', { id: profileId })
}

async function openProduct() {
  if (!conversationId.value && !productId.value) return
  try {
    const product = conversationId.value ? await getConversationProduct(conversationId.value) : null
    const id = product?.id || productId.value
    if (!id) return
    const image = Array.isArray(product?.images) ? product.images[0] : (product?.coverImage || product?.image || '')
    navigate('/pages/detail/detail', {
      id,
      readonly: product?.status === 'ON_SALE' ? '' : '1',
      snapshotTitle: product?.title ? encodeURIComponent(product.title) : '',
      snapshotPrice: product?.price || '',
      snapshotImage: image ? encodeURIComponent(image) : '',
      snapshotStatus: product?.status || '',
      snapshotSellerId: product?.sellerId || product?.seller?.id || '',
      snapshotSellerName: product?.sellerName || product?.seller?.nickname || '',
      sellerUnavailable: isTargetUnavailable.value ? '1' : ''
    })
  } catch (error) {
    if (productId.value) {
      navigate('/pages/detail/detail', { id: productId.value })
      return
    }
    uni.showToast({ title: '商品暂不可查看', icon: 'none' })
  }
}

onLoad(async (options) => {
  conversationId.value = options.conversationId || ''
  title.value = options.title ? decodeURIComponent(options.title) : ''
  targetUserId.value = options.targetUserId || ''
  targetNicknameSnapshot.value = options.targetNickname ? decodeURIComponent(options.targetNickname) : ''
  targetAvatarSnapshot.value = options.targetAvatar ? decodeURIComponent(options.targetAvatar) : ''
  targetBannedSnapshot.value = options.targetBanned === '1' || options.targetBanned === 1
  productId.value = options.productId || ''
  seedInitialProfiles()
  await loadInitialTargetProfile()
  await load()
  await loadProfiles()
})

onPullDownRefresh(async () => {
  try {
    await load()
    await loadProfiles()
  } finally {
    uni.stopPullDownRefresh()
  }
})
</script>

<style scoped>
.page { min-height: 100vh; padding-bottom: 120rpx; background: #f5f8f6; box-sizing: border-box; }
.chat-head { padding: 28rpx; background: #fff; }
.chat-title { color: #26342f; font-size: 32rpx; font-weight: 700; }
.clickable-title { display: flex; align-items: center; gap: 16rpx; }
.profile-link { color: #23734f; font-size: 23rpx; font-weight: 500; }
.chat-subtitle { margin-top: 10rpx; color: #89938f; font-size: 23rpx; }
.product-link { display: flex; align-items: center; justify-content: space-between; gap: 16rpx; margin: 18rpx 28rpx 0; padding: 18rpx 22rpx; border-radius: 16rpx; background: #fff; box-shadow: 0 8rpx 24rpx rgba(23, 33, 43, .04); }
.product-link-body { min-width: 0; flex: 1; }
.product-link-label { display: block; color: #89938f; font-size: 22rpx; }
.product-link-title { display: block; overflow: hidden; margin-top: 6rpx; color: #26342f; font-size: 27rpx; font-weight: 700; text-overflow: ellipsis; white-space: nowrap; }
.product-link-arrow { color: #9aa5a1; font-size: 42rpx; line-height: 1; }
.messages { height: calc(100vh - 210rpx); padding: 20rpx 28rpx 24rpx; box-sizing: border-box; }
.messages.has-product { height: calc(100vh - 318rpx); }
.messages.has-restriction { height: calc(100vh - 268rpx); }
.messages.has-product.has-restriction { height: calc(100vh - 376rpx); }
.restriction-banner { padding: 14rpx 28rpx; background: #fef3f2; border-bottom: 1rpx solid #fee4e2; text-align: center; }
.restriction-text { color: #d92d20; font-size: 26rpx; font-weight: 500; }
.time-divider { width: fit-content; max-width: 420rpx; margin: 18rpx auto; padding: 6rpx 16rpx; border-radius: 999rpx; background: #dfe7e3; color: #7b8782; font-size: 21rpx; text-align: center; }
.message-swipe { position: relative; overflow: hidden; margin: 0 -28rpx 18rpx; padding: 0 28rpx; }
.message-delete { position: absolute; top: 6rpx; right: 28rpx; bottom: 6rpx; display: flex; width: 120rpx; align-items: center; justify-content: center; border-radius: 16rpx; background: #f04444; color: #fff; font-size: 26rpx; opacity: 0; transition: opacity .12s ease; }
.message-swipe.active .message-delete { opacity: 1; }
.message-front { position: relative; z-index: 1; width: 100%; min-height: 58rpx; background: #f5f8f6; transition: transform .18s ease; }
.message-front.swiped { transform: translateX(-140rpx); }
.message-row { display: flex; align-items: flex-end; gap: 12rpx; }
.message-row.mine { justify-content: flex-end; }
.bubble { max-width: 520rpx; padding: 18rpx 22rpx; border-radius: 18rpx; background: #fff; color: #26342f; box-sizing: border-box; }
.avatar { display: flex; width: 58rpx; height: 58rpx; flex-shrink: 0; align-items: center; justify-content: center; border-radius: 50%; color: #fff; font-size: 22rpx; font-weight: 700; }
.seller-avatar { background: #6b8b7e; }
.buyer-avatar { background: #23734f; }
.banned-avatar { background: #d92d20; color: #fff; font-size: 25rpx; }
.image-avatar { background: #e8ecef; }
.mine .bubble { background: #23734f; color: #fff; }
.content { font-size: 27rpx; line-height: 1.55; word-break: break-word; }
.empty { margin-top: 180rpx; color: #929c98; text-align: center; }
.composer { position: fixed; right: 0; bottom: 0; left: 0; display: flex; gap: 14rpx; padding: 18rpx 22rpx calc(18rpx + env(safe-area-inset-bottom)); background: #fff; box-sizing: border-box; }
.composer.restricted { background: #f5f5f5; }
.composer.restricted .input { background: #e8e8e8; color: #999; }
.composer.restricted .send { background: #b8c5c0; }
.input { flex: 1; height: 76rpx; padding: 0 24rpx; border-radius: 999rpx; background: #f2f5f3; font-size: 27rpx; box-sizing: border-box; }
.send { width: 132rpx; height: 76rpx; border-radius: 999rpx; background: #23734f; color: #fff; font-size: 27rpx; line-height: 76rpx; }
.send[disabled] { background: #b8c5c0; }
</style>
