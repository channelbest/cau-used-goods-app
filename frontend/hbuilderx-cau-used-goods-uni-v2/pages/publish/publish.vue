<template>
  <view class="page">
    <view v-if="!canPublish" class="locked-card">
      <view class="locked-title">{{ lockedTitle }}</view>
      <view class="locked-desc">{{ lockedDesc }}</view>
      <button v-if="showStudentAuthButton" class="locked-button" @click="goStudentAuth">{{ lockedButtonText }}</button>
    </view>

    <view v-else>
      <view class="hero">
        <text class="eyebrow">CAU CAMPUS MARKET</text>
        <text class="title">{{ editMode ? '编辑商品信息' : '发布你的闲置好物' }}</text>
        <text class="hero-copy">{{ editMode ? '修改商品信息和图片后保存，买家会看到最新内容' : '真实描述物品情况，更容易遇到合适的新主人' }}</text>
      </view>

      <view class="card">
        <view class="field">
          <view class="field-head">
            <text>商品图片</text>
            <text class="muted">{{ form.images.length }}/{{ MAX_IMAGES }}</text>
          </view>

          <view class="upload-row">
            <view v-if="form.images.length < MAX_IMAGES" class="upload" @click="chooseImages">
              +
              <text>添加图片</text>
            </view>
            <view class="upload-copy">
              <text>选择图片后可拖动、缩放，自主裁剪 4:3 商品图。</text>
              <text>确认裁剪后在下方预览。</text>
            </view>
          </view>

          <view v-if="previewImages.length" class="preview-block">
            <view class="preview-head">
              <text>裁剪预览</text>
              <text class="muted">商品卡片展示效果</text>
            </view>
            <view class="image-grid">
              <view v-for="(image, index) in previewImages" :key="image + index" class="image-item" @click="replaceImage(index)">
                <image :src="image" mode="aspectFill" @click.stop="replaceImage(index)" @tap.stop="replaceImage(index)" />
                <text class="ratio-tag">4:3</text>
                <text class="edit-tag">点击修改</text>
                <text class="remove" @click.stop="removeImage(index)">x</text>
              </view>
            </view>
          </view>
        </view>

        <view class="field">
          <view class="field-head">
            <text>商品标题</text>
            <text class="ai" @click="optimizeTitle">AI 优化标题</text>
          </view>
          <input v-model.trim="form.title" maxlength="100" placeholder="例如：九成新小米台灯，宿舍自用" />
        </view>

        <view class="field inline">
          <text>一级分类</text>
          <picker :range="primaryCategories" range-key="name" @change="changePrimaryCategory">
            <view class="picker">{{ primaryName || '请选择一级分类' }} ›</view>
          </picker>
        </view>

        <view v-if="childCategories.length" class="field inline">
          <text>二级分类</text>
          <picker :range="childCategories" range-key="name" @change="changeChildCategory">
            <view class="picker">{{ categoryName || '请选择二级分类' }} ›</view>
          </picker>
        </view>

        <view class="field inline">
          <text>成色</text>
          <picker :range="conditions" @change="changeCondition">
            <view class="picker">{{ form.conditionLevel || '请选择成色' }} ›</view>
          </picker>
        </view>

        <view class="price-row">
          <view class="field">
            <text>原价</text>
            <input v-model="form.originalPrice" type="digit" placeholder="选填" />
          </view>
          <view class="field">
            <text>售价</text>
            <input v-model="form.price" type="digit" placeholder="必填" />
          </view>
        </view>

        <view class="field">
          <view class="field-head">
            <text>商品描述</text>
            <text class="muted">{{ form.description.length }}/{{ DESCRIPTION_LIMIT }}</text>
          </view>
          <textarea v-model.trim="form.description" :maxlength="DESCRIPTION_LIMIT" placeholder="说明使用情况、外观瑕疵、配件等信息" />
          <view class="field-foot">
            <text class="ai" @click="generateDescription">AI 生成描述</text>
          </view>
        </view>

        <view class="field last-field">
          <text>建议面交地点</text>
          <input v-model.trim="form.meetLocation" maxlength="100" placeholder="例如：东区图书馆门口" />
        </view>
      </view>

      <button class="submit" :loading="submitting" @click="submit">{{ editMode ? '保存修改' : '确认发布' }}</button>
    </view>

    <view v-if="cropperVisible" class="crop-mask">
      <view class="crop-panel">
        <view class="crop-head">
          <text>裁剪商品图</text>
          <text class="crop-close" @click="cancelCrop">取消</text>
        </view>
        <view
          class="crop-frame"
          :style="{ width: cropBox.width + 'px', height: cropBox.height + 'px' }"
          @touchstart.stop="startDrag"
          @touchmove.stop.prevent="moveDrag"
          @touchend.stop="endDrag"
        >
          <image
            class="crop-image"
            :src="cropSource"
            mode="aspectFit"
            :style="cropImageStyle"
            draggable="false"
          />
          <view class="crop-border"></view>
        </view>
        <view class="zoom-row">
          <text>缩小/放大</text>
          <slider
            class="zoom-slider"
            :value="cropState.zoomPercent"
            min="50"
            max="300"
            block-size="18"
            activeColor="#23734f"
            backgroundColor="#dbe5df"
            @changing="changeZoom"
            @change="changeZoom"
          />
        </view>
        <view class="crop-actions">
          <button class="crop-secondary" @click="resetCropPosition">居中</button>
          <button class="crop-primary" @click="confirmCrop">确认裁剪</button>
        </view>
      </view>
    </view>

    <canvas
      canvas-id="cropCanvas"
      class="crop-canvas"
      :style="{ width: canvasSize.width + 'px', height: canvasSize.height + 'px' }"
    ></canvas>
  </view>
</template>

<script setup>
import { computed, getCurrentInstance, reactive, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import {
  createProduct,
  generateProductDescription,
  listCategories,
  optimizeProductTitle,
  replaceProductImages,
  updateProduct,
  uploadProductImage
} from '../../api/product'
import { getCurrentUser } from '../../api/auth'
import { getToken, getUser, setUser } from '../../utils/auth'
import { normalizeImage } from '../../utils/product-format'
import { accountStatusOf, userTradeRestrictionMessage } from '../../utils/user-format'

const MAX_SIZE = 5 * 1024 * 1024
const MAX_IMAGES = 9
const CANVAS_WIDTH = 800
const CANVAS_HEIGHT = 600
const DESCRIPTION_LIMIT = 500
const conditions = ['全新', '九成新', '八成新', '七成新', '有明显使用痕迹']
const categories = ref([])
const selectedPrimaryId = ref('')
const selectedChildName = ref('')
const submitting = ref(false)
const canPublish = ref(false)
const publishLockedMessage = ref('\u672a\u5b8c\u6210\u5b66\u751f\u8ba4\u8bc1')
const publishLockedType = ref('auth')
const previewImages = ref([])
const editProductId = ref('')
const originalImages = ref([])
const newImages = ref([])
const imageFlowActive = ref(false)
const cropperVisible = ref(false)
const cropSource = ref('')
const cropResolve = ref(null)
const instance = getCurrentInstance()
const form = reactive({
  images: [],
  title: '',
  categoryId: '',
  conditionLevel: '',
  originalPrice: '',
  price: '',
  description: '',
  meetLocation: ''
})
const canvasSize = reactive({ width: CANVAS_WIDTH, height: CANVAS_HEIGHT })
const cropBox = reactive({ width: 320, height: 240 })
const cropState = reactive({
  imageWidth: 0,
  imageHeight: 0,
  baseScale: 1,
  zoomPercent: 100,
  offsetX: 0,
  offsetY: 0,
  startX: 0,
  startY: 0,
  startOffsetX: 0,
  startOffsetY: 0,
  dragging: false
})

const parentIdOf = (item) => Number(item?.parentId || item?.parent_id || item?.pid || 0)
const normalizeCategories = (list = [], parentId = 0) => {
  const source = Array.isArray(list) ? list : (list.items || list.list || list.data || [])
  const result = []
  ;(source || []).forEach((item) => {
    const children = item.children || item.childList || item.subCategories || item.sub_categories || []
    const normalized = {
      ...item,
      id: item.id || item.categoryId || item.category_id,
      name: item.name || item.categoryName || item.category_name,
      parentId: item.parentId || item.parent_id || parentId || 0
    }
    if (Number(normalized.id) !== 0) result.push(normalized)
    if (Array.isArray(children) && children.length) {
      result.push(...normalizeCategories(children, normalized.id))
    }
  })
  return result.filter((item) => item.id && item.name)
}
const primaryCategories = computed(() => categories.value.filter((item) => parentIdOf(item) === 0))
const childCategories = computed(() => {
  if (!selectedPrimaryId.value) return []
  return categories.value.filter((item) => parentIdOf(item) === Number(selectedPrimaryId.value) && Number(item.id) !== Number(selectedPrimaryId.value))
})
const primaryName = computed(() => primaryCategories.value.find((item) => Number(item.id) === Number(selectedPrimaryId.value))?.name || '')
const categoryName = computed(() => selectedChildName.value || categories.value.find((item) => Number(item.id) === Number(form.categoryId))?.name || '')
const toast = (title) => uni.showToast({ title, icon: 'none' })
const zoom = computed(() => cropState.zoomPercent / 100)
const displayWidth = computed(() => cropState.imageWidth * cropState.baseScale * zoom.value)
const displayHeight = computed(() => cropState.imageHeight * cropState.baseScale * zoom.value)
const cropImageStyle = computed(() => ({
  width: `${displayWidth.value}px`,
  height: `${displayHeight.value}px`,
  transform: `translate(${cropState.offsetX}px, ${cropState.offsetY}px)`
}))
const editMode = computed(() => !!editProductId.value)
const lockedTitle = computed(() => {
  if (publishLockedType.value === 'account') return '\u8d26\u53f7\u72b6\u6001\u53d7\u9650'
  if (publishLockedType.value === 'login') return '\u8bf7\u5148\u767b\u5f55'
  return '\u5148\u5b8c\u6210\u8ba4\u8bc1\uff0c\u518d\u53d1\u5e03\u95f2\u7f6e'
})
const lockedDesc = computed(() => {
  if (publishLockedType.value === 'account') return `${publishLockedMessage.value}\u3002\u5982\u9700\u7533\u8bc9\uff0c\u8bf7\u8054\u7cfb\u5e73\u53f0\u7ba1\u7406\u5458\u5904\u7406\u3002`
  if (publishLockedType.value === 'login') return '\u767b\u5f55\u540e\u624d\u80fd\u53d1\u5e03\u548c\u7ba1\u7406\u95f2\u7f6e\u5546\u54c1\u3002'
  return '\u4f60\u7684\u5b66\u751f\u8ba4\u8bc1\u5c1a\u672a\u901a\u8fc7\uff0c\u53d1\u5e03\u8868\u5355\u5df2\u6682\u65f6\u5173\u95ed\u3002\u8ba4\u8bc1\u901a\u8fc7\u540e\u5373\u53ef\u53d1\u5e03\u5546\u54c1\u3002'
})
const showStudentAuthButton = computed(() => publishLockedType.value === 'auth')
const lockedButtonText = computed(() => '\u53bb\u5b66\u751f\u8ba4\u8bc1')

const resetForm = () => {
  editProductId.value = ''
  originalImages.value = []
  newImages.value = []
  form.images = []
  previewImages.value = []
  form.title = ''
  selectedPrimaryId.value = ''
  selectedChildName.value = ''
  form.categoryId = ''
  form.conditionLevel = ''
  form.originalPrice = ''
  form.price = ''
  form.description = ''
  form.meetLocation = ''
}

const clearEditStorage = () => {
  uni.removeStorageSync('PUBLISH_EDIT_PRODUCT_ID')
  uni.removeStorageSync('PUBLISH_EDIT_PRODUCT_DATA')
  uni.removeStorageSync('PUBLISH_EDIT_INTENT')
}

const fillEditForm = (product = {}) => {
  editProductId.value = String(product.id || '')
  form.images = Array.isArray(product.images) ? product.images.slice() : []
  originalImages.value = form.images.slice()
  newImages.value = []
  previewImages.value = form.images.map((image) => normalizeImage(image)).filter(Boolean)
  form.title = product.title || ''
  form.categoryId = product.categoryId || product.category_id || ''
  form.conditionLevel = product.conditionLevel || product.condition_level || ''
  form.originalPrice = product.originalPrice ?? product.original_price ?? ''
  form.price = product.price ?? ''
  form.description = product.description || ''
  form.meetLocation = product.meetLocation || product.meet_location || ''

  const category = categories.value.find((item) => Number(item.id) === Number(form.categoryId))
  if (category && parentIdOf(category)) {
    selectedPrimaryId.value = category.parentId
    selectedChildName.value = category.name
  } else if (category) {
    selectedPrimaryId.value = category.id
    selectedChildName.value = ''
  }
}

const applyEditState = () => {
  const hasEditIntent = uni.getStorageSync('PUBLISH_EDIT_INTENT') === '1'
  uni.removeStorageSync('PUBLISH_EDIT_INTENT')
  if (!hasEditIntent) {
    clearEditStorage()
    return
  }
  const id = uni.getStorageSync('PUBLISH_EDIT_PRODUCT_ID')
  const product = uni.getStorageSync('PUBLISH_EDIT_PRODUCT_DATA')
  if (!id || !product?.id) {
    clearEditStorage()
    if (editProductId.value) resetForm()
    return
  }
  if (String(id) === String(editProductId.value)) return
  fillEditForm(product)
}

const isVerified = (user) => {
  const status = user?.authStatus || user?.auth_status || ''
  return status === 'VERIFIED'
}

const normalizeProductError = (error, action = 'publish') => {
  const message = String(error?.message || '')
  if (/BANNED|PERM_BANNED|PERMANENT_BANNED|\u5c01\u7981/i.test(message)) return userTradeRestrictionMessage('BANNED', action)
  if (/DISABLED|\u7981\u7528/i.test(message)) return userTradeRestrictionMessage('DISABLED', action)
  if (/CANCELED|CANCELLED|\u6ce8\u9500/i.test(message)) return userTradeRestrictionMessage('CANCELED', action)
  if (/403|FORBIDDEN|PERMISSION|VERIFY|VERIFIED|\u6743\u9650/i.test(message)) return action === 'sale' ? '\u8d26\u53f7\u72b6\u6001\u4e0d\u6ee1\u8db3\u4e0a\u67b6\u6761\u4ef6' : '\u8d26\u53f7\u72b6\u6001\u4e0d\u6ee1\u8db3\u53d1\u5e03\u6761\u4ef6'
  return message || (action === 'sale' ? '\u4e0a\u67b6\u5931\u8d25\uff0c\u8bf7\u7a0d\u540e\u91cd\u8bd5' : '\u53d1\u5e03\u5931\u8d25')
}

const loadPublishState = async () => {
  if (!getToken()) {
    canPublish.value = false
    publishLockedType.value = 'login'
    publishLockedMessage.value = '\u8bf7\u5148\u767b\u5f55'
    uni.showToast({ title: publishLockedMessage.value, icon: 'none' })
    return
  }

  try {
    const cached = getUser() || {}
    const current = { ...cached, ...(await getCurrentUser()) }
    setUser(current)
    const restrictionMessage = userTradeRestrictionMessage(accountStatusOf(current), 'publish')
    if (restrictionMessage) {
      canPublish.value = false
      publishLockedType.value = 'account'
      publishLockedMessage.value = restrictionMessage
      toast(restrictionMessage)
      return
    }
    canPublish.value = isVerified(current)
    if (!canPublish.value) {
      publishLockedType.value = 'auth'
      publishLockedMessage.value = '\u672a\u5b8c\u6210\u5b66\u751f\u8ba4\u8bc1'
      uni.showToast({ title: publishLockedMessage.value, icon: 'none' })
      return
    }
    publishLockedType.value = ''
    publishLockedMessage.value = ''
    categories.value = normalizeCategories(await listCategories())
    if (!imageFlowActive.value) applyEditState()
  } catch (error) {
    canPublish.value = false
    const cachedRestriction = userTradeRestrictionMessage(accountStatusOf(getUser() || {}), 'publish')
    publishLockedType.value = 'account'
    publishLockedMessage.value = cachedRestriction || normalizeProductError(error, 'publish')
    toast(publishLockedMessage.value || '\u52a0\u8f7d\u5931\u8d25')
  }
}

onShow(loadPublishState)

const goStudentAuth = () => uni.navigateTo({ url: '/pages/student-auth/student-auth' })
const changePrimaryCategory = ({ detail }) => {
  const primary = primaryCategories.value[detail.value]
  selectedPrimaryId.value = primary?.id || ''
  selectedChildName.value = ''
  form.categoryId = childCategories.value.length ? '' : selectedPrimaryId.value
}
const changeChildCategory = ({ detail }) => {
  const child = childCategories.value[detail.value]
  form.categoryId = child?.id || ''
  selectedChildName.value = child?.name || ''
}
const changeCondition = ({ detail }) => { form.conditionLevel = conditions[detail.value] }
const removeImage = (index) => {
  const image = form.images[index]
  if (newImages.value.includes(image)) {
    newImages.value = newImages.value.filter((item) => item !== image)
  }
  form.images.splice(index, 1)
  previewImages.value.splice(index, 1)
}
const appendImage = (imageUrl, previewUrl) => {
  if (form.images.length >= MAX_IMAGES) return
  form.images.push(imageUrl)
  previewImages.value.push(normalizeImage(previewUrl || imageUrl))
}
const replaceImageAt = (index, imageUrl, previewUrl) => {
  if (index < 0 || index >= form.images.length) return
  const oldImage = form.images[index]
  form.images.splice(index, 1, imageUrl)
  previewImages.value.splice(index, 1, normalizeImage(previewUrl || imageUrl))
  if (editMode.value) {
    newImages.value = newImages.value.filter((item) => item !== oldImage)
    newImages.value.push(imageUrl)
  }
}

const clampOffset = () => {
  const width = displayWidth.value
  const height = displayHeight.value
  cropState.offsetX = width <= cropBox.width
    ? (cropBox.width - width) / 2
    : Math.min(0, Math.max(cropBox.width - width, cropState.offsetX))
  cropState.offsetY = height <= cropBox.height
    ? (cropBox.height - height) / 2
    : Math.min(0, Math.max(cropBox.height - height, cropState.offsetY))
}

const resetCropPosition = () => {
  cropState.zoomPercent = 100
  cropState.offsetX = (cropBox.width - displayWidth.value) / 2
  cropState.offsetY = (cropBox.height - displayHeight.value) / 2
  clampOffset()
}

const openCropper = (src) => {
  return new Promise((resolve) => {
    uni.getImageInfo({
      src,
      success: (info) => {
        const system = uni.getSystemInfoSync()
        cropBox.width = Math.min(system.windowWidth - 56, 360)
        cropBox.height = Math.round(cropBox.width * 3 / 4)
        cropState.imageWidth = info.width
        cropState.imageHeight = info.height
        cropState.baseScale = Math.min(cropBox.width / info.width, cropBox.height / info.height)
        cropState.zoomPercent = 100
        cropSource.value = src
        cropResolve.value = resolve
        cropperVisible.value = true
        setTimeout(resetCropPosition, 20)
      },
      fail: () => {
        toast('图片读取失败')
        resolve('')
      }
    })
  })
}

const startDrag = (event) => {
  const touch = event.touches?.[0]
  if (!touch) return
  cropState.dragging = true
  cropState.startX = touch.clientX
  cropState.startY = touch.clientY
  cropState.startOffsetX = cropState.offsetX
  cropState.startOffsetY = cropState.offsetY
}

const moveDrag = (event) => {
  if (!cropState.dragging) return
  const touch = event.touches?.[0]
  if (!touch) return
  cropState.offsetX = cropState.startOffsetX + touch.clientX - cropState.startX
  cropState.offsetY = cropState.startOffsetY + touch.clientY - cropState.startY
  clampOffset()
}

const endDrag = () => {
  cropState.dragging = false
}

const changeZoom = ({ detail }) => {
  const oldWidth = displayWidth.value
  const oldHeight = displayHeight.value
  const centerX = cropBox.width / 2
  const centerY = cropBox.height / 2
  const ratioX = oldWidth ? (centerX - cropState.offsetX) / oldWidth : 0.5
  const ratioY = oldHeight ? (centerY - cropState.offsetY) / oldHeight : 0.5
  cropState.zoomPercent = detail.value
  cropState.offsetX = centerX - displayWidth.value * ratioX
  cropState.offsetY = centerY - displayHeight.value * ratioY
  clampOffset()
}

const finishCrop = (path) => {
  const resolve = cropResolve.value
  cropperVisible.value = false
  cropResolve.value = null
  cropSource.value = ''
  if (resolve) resolve(path)
}

const cancelCrop = () => {
  finishCrop('')
}

const confirmCrop = () => {
  const ctx = uni.createCanvasContext('cropCanvas', instance?.proxy)
  const scaleX = CANVAS_WIDTH / cropBox.width
  const scaleY = CANVAS_HEIGHT / cropBox.height
  ctx.clearRect(0, 0, CANVAS_WIDTH, CANVAS_HEIGHT)
  ctx.setFillStyle('#fff')
  ctx.fillRect(0, 0, CANVAS_WIDTH, CANVAS_HEIGHT)
  ctx.drawImage(
    cropSource.value,
    cropState.offsetX * scaleX,
    cropState.offsetY * scaleY,
    displayWidth.value * scaleX,
    displayHeight.value * scaleY
  )
  ctx.draw(false, () => {
    uni.canvasToTempFilePath({
      canvasId: 'cropCanvas',
      x: 0,
      y: 0,
      width: CANVAS_WIDTH,
      height: CANVAS_HEIGHT,
      destWidth: CANVAS_WIDTH,
      destHeight: CANVAS_HEIGHT,
      fileType: 'jpg',
      quality: 0.9,
      success: (res) => finishCrop(res.tempFilePath),
      fail: () => {
        toast('裁剪失败，请重新选择图片')
        finishCrop('')
      }
    }, instance?.proxy)
  })
}

const chooseImageFiles = async (count) => {
  const result = await uni.chooseImage({
    count,
    sizeType: ['compressed']
  })
  const selectedFiles = Array.isArray(result.tempFiles) && result.tempFiles.length
    ? result.tempFiles
    : (result.tempFilePaths || []).map((path) => ({ path, size: 0 }))
  return selectedFiles.filter((file) => {
    if (file.size > MAX_SIZE) {
      toast('单张图片不能超过 5 MB')
      return false
    }
    return true
  })
}

const cropAndUploadImage = async (file, loadingTitle = '上传裁剪图中') => {
  const filePath = file.path || file.tempFilePath
  if (!filePath) return ''
  const croppedPath = await openCropper(filePath)
  if (!croppedPath) return ''
  uni.showLoading({ title: loadingTitle })
  try {
    const uploaded = await uploadProductImage(croppedPath)
    return uploaded.imageUrl
  } finally {
    uni.hideLoading()
  }
}

const replaceImage = async (index) => {
  if (!canPublish.value) return toast(publishLockedMessage.value || '\u672a\u5b8c\u6210\u5b66\u751f\u8ba4\u8bc1')
  imageFlowActive.value = true
  try {
    const files = await chooseImageFiles(1)
    const file = files[0]
    if (!file) return
    const imageUrl = await cropAndUploadImage(file, '上传新图片中')
    if (!imageUrl) return
    replaceImageAt(index, imageUrl)
  } catch (error) {
    if (!error?.errMsg?.includes('cancel')) toast(error.message || '图片修改失败')
  } finally {
    imageFlowActive.value = false
    uni.hideLoading()
  }
}

const chooseImages = async () => {
  if (!canPublish.value) return toast(publishLockedMessage.value || '\u672a\u5b8c\u6210\u5b66\u751f\u8ba4\u8bc1')
  const remaining = MAX_IMAGES - form.images.length
  if (remaining <= 0) return toast(`最多上传 ${MAX_IMAGES} 张商品图片`)

  imageFlowActive.value = true
  try {
    const files = (await chooseImageFiles(remaining)).slice(0, remaining)
    if (!files.length) return

    for (const file of files) {
      if (form.images.length >= MAX_IMAGES) break
      const imageUrl = await cropAndUploadImage(file)
      if (!imageUrl) continue
      appendImage(imageUrl)
      if (editMode.value) newImages.value.push(imageUrl)
    }
  } catch (error) {
    if (!error?.errMsg?.includes('cancel')) toast(error.message || '图片选择或上传失败')
  } finally {
    imageFlowActive.value = false
    uni.hideLoading()
  }
}

const optimizeTitle = async () => {
  if (!canPublish.value) return toast(publishLockedMessage.value || '\u672a\u5b8c\u6210\u5b66\u751f\u8ba4\u8bc1')
  if (!form.title) return toast('请先填写一个基础标题')
  try {
    uni.showLoading({ title: 'AI 正在优化' })
    const { titles } = await optimizeProductTitle({
      title: form.title,
      categoryName: categoryName.value,
      conditionLevel: form.conditionLevel
    })
    uni.showActionSheet({
      itemList: titles,
      success: ({ tapIndex }) => { form.title = titles[tapIndex] }
    })
  } catch (error) {
    toast('AI 服务暂时不可用')
  } finally {
    uni.hideLoading()
  }
}

const generateDescription = async () => {
  if (!canPublish.value) return toast(publishLockedMessage.value || '\u672a\u5b8c\u6210\u5b66\u751f\u8ba4\u8bc1')
  if (!form.title) return toast('请先填写标题')
  try {
    uni.showLoading({ title: 'AI 正在生成' })
    const result = await generateProductDescription({
      title: form.title,
      categoryName: categoryName.value,
      conditionLevel: form.conditionLevel,
      meetLocation: form.meetLocation
    })
    form.description = String(result.description || '').slice(0, DESCRIPTION_LIMIT)
  } catch (error) {
    toast('AI 服务暂时不可用')
  } finally {
    uni.hideLoading()
  }
}

const validate = () => {
  if (!canPublish.value) return publishLockedMessage.value || '\u672a\u5b8c\u6210\u5b66\u751f\u8ba4\u8bc1'
  if (!form.images.length) return '请至少上传一张商品图片'
  if (!form.title) return '请填写商品标题'
  if (!selectedPrimaryId.value) return '请选择一级分类'
  if (childCategories.value.length && !form.categoryId) return '请选择二级分类'
  if (!form.categoryId) form.categoryId = selectedPrimaryId.value
  if (!form.conditionLevel) return '请选择商品成色'
  if (!form.price || Number(form.price) <= 0) return '请填写正确的商品售价'
  if (!form.description) return '请填写商品描述'
  if (form.description.length > DESCRIPTION_LIMIT) return `商品描述不能超过 ${DESCRIPTION_LIMIT} 字`
  return ''
}
const submit = async () => {
  const message = validate()
  if (message) return toast(message)
  submitting.value = true
  try {
    const payload = {
      ...form,
      price: Number(form.price),
      originalPrice: Number(form.originalPrice || 0)
    }
    if (editMode.value) {
      const { images, ...productPayload } = payload
      await updateProduct(editProductId.value, productPayload)
      await replaceProductImages(editProductId.value, images)
      clearEditStorage()
      resetForm()
      uni.setStorageSync('PRODUCT_LIST_DIRTY', true)
      uni.showToast({ title: '修改成功', icon: 'success' })
      setTimeout(() => uni.navigateTo({ url: '/pages/my-products/my-products' }), 600)
      return
    }
    await createProduct(payload)
    resetForm()
    uni.setStorageSync('PRODUCT_LIST_DIRTY', true)
    uni.showToast({ title: '发布成功', icon: 'success' })
    setTimeout(() => uni.switchTab({ url: '/pages/home/home' }), 600)
  } catch (error) {
    toast(normalizeProductError(error, 'publish'))
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.page { min-height: 100vh; padding-bottom: 36rpx; background: #f5f8f6; box-sizing: border-box; }
.hero { margin: 24rpx 28rpx 0; padding: 34rpx 30rpx; border-radius: 24rpx; background: #2d855f; color: #fff; }
.eyebrow, .title, .hero-copy, .muted { display: block; }
.eyebrow { color: rgba(255,255,255,.72); font-size: 20rpx; letter-spacing: 3rpx; }
.title { margin-top: 12rpx; font-size: 38rpx; font-weight: 700; }
.hero-copy { margin-top: 10rpx; color: rgba(255,255,255,.78); font-size: 24rpx; }
.locked-card, .card { margin: 28rpx; padding: 28rpx; border-radius: 18rpx; background: #fff; }
.locked-title { color: #26342f; font-size: 34rpx; font-weight: 700; }
.locked-desc { margin-top: 16rpx; color: #7d8984; font-size: 26rpx; line-height: 1.6; }
.locked-button, .submit { margin-top: 30rpx; border-radius: 999rpx; background: #23734f; color: #fff; }
.field { margin-bottom: 26rpx; }
.last-field { margin-bottom: 0; }
.field-head, .inline, .price-row, .preview-head, .field-foot { display: flex; align-items: center; justify-content: space-between; gap: 18rpx; }
.field text, .field-head { color: #26342f; font-size: 27rpx; font-weight: 600; }
.muted { color: #89938f; font-size: 22rpx; font-weight: 400; }
.ai { color: #23734f !important; font-size: 24rpx !important; font-weight: 600 !important; }
input, textarea, .picker { width: 100%; margin-top: 14rpx; padding: 18rpx 22rpx; border-radius: 16rpx; background: #f3f6f4; color: #26342f; font-size: 27rpx; box-sizing: border-box; }
input, .picker { min-height: 72rpx; }
textarea { height: 180rpx; }
.field-foot { margin-top: 10rpx; justify-content: flex-end; }
.inline .picker { min-width: 360rpx; margin-top: 0; text-align: right; }
.price-row .field { flex: 1; margin-bottom: 0; }
.upload-row { display: flex; align-items: center; gap: 18rpx; margin-top: 16rpx; }
.upload { display: flex; width: 176rpx; height: 132rpx; flex-shrink: 0; flex-direction: column; align-items: center; justify-content: center; border-radius: 14rpx; background: #edf4f1; color: #23734f; font-size: 42rpx; }
.upload text { margin-top: 4rpx; color: #23734f; font-size: 22rpx; }
.upload-copy { display: flex; min-width: 0; flex: 1; flex-direction: column; gap: 8rpx; }
.upload-copy text { color: #7d8984; font-size: 23rpx; font-weight: 400; line-height: 1.35; }
.preview-block { margin-top: 20rpx; }
.preview-head { margin-bottom: 12rpx; }
.image-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16rpx; }
.image-item { position: relative; aspect-ratio: 4 / 3; overflow: hidden; border-radius: 14rpx; background: #edf4f1; }
.image-item image { width: 100%; height: 100%; }
.ratio-tag, .remove, .edit-tag { position: absolute; border-radius: 999rpx; color: #fff !important; font-size: 20rpx !important; font-weight: 700 !important; line-height: 34rpx; text-align: center; }
.ratio-tag, .remove { top: 8rpx; }
.ratio-tag { left: 8rpx; padding: 0 10rpx; background: rgba(35,115,79,.82); }
.remove { right: 8rpx; width: 34rpx; height: 34rpx; background: rgba(0,0,0,.55); }
.edit-tag { right: 8rpx; bottom: 8rpx; padding: 0 10rpx; background: rgba(0,0,0,.48); }
.submit { margin: 0 28rpx; height: 88rpx; line-height: 88rpx; font-size: 30rpx; }
.crop-mask { position: fixed; z-index: 99; top: 0; right: 0; bottom: 0; left: 0; display: flex; align-items: center; justify-content: center; padding: 28rpx; background: rgba(0,0,0,.58); box-sizing: border-box; }
.crop-panel { width: 100%; padding: 28rpx; border-radius: 20rpx; background: #fff; box-sizing: border-box; }
.crop-head, .zoom-row, .crop-actions { display: flex; align-items: center; justify-content: space-between; gap: 18rpx; }
.crop-head { margin-bottom: 20rpx; color: #26342f; font-size: 30rpx; font-weight: 700; }
.crop-close { color: #89938f; font-size: 25rpx; font-weight: 500; }
.crop-frame { position: relative; overflow: hidden; margin: 0 auto; background: #111; touch-action: none; }
.crop-image { position: absolute; top: 0; left: 0; will-change: transform; }
.crop-border { position: absolute; top: 0; right: 0; bottom: 0; left: 0; border: 4rpx solid rgba(255,255,255,.96); box-shadow: inset 0 0 0 2rpx rgba(35,115,79,.7); pointer-events: none; }
.zoom-row { margin-top: 22rpx; color: #26342f; font-size: 25rpx; }
.zoom-slider { flex: 1; }
.crop-actions { margin-top: 18rpx; }
.crop-secondary, .crop-primary { flex: 1; height: 78rpx; border-radius: 999rpx; font-size: 27rpx; line-height: 78rpx; }
.crop-secondary { background: #edf4f1; color: #23734f; }
.crop-primary { background: #23734f; color: #fff; }
.crop-canvas { position: fixed; left: -9999px; top: -9999px; opacity: 0; pointer-events: none; }
</style>





