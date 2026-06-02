<template>
  <view class="page">
    <view class="intro"><text class="title">发布你的闲置好物</text><text class="muted">真实描述物品情况，更容易遇到合适的新主人</text></view>
    <view class="card">
      <view class="field"><view class="field-head"><text>商品图片</text><text class="muted">{{ form.images.length }}/9</text></view>
        <view class="image-grid"><view v-for="(image,index) in form.images" :key="image + index" class="image-item"><image :src="image" mode="aspectFill" /><text class="remove" @click="removeImage(index)">×</text></view><view v-if="form.images.length < 9" class="upload" @click="chooseImages">+<text>添加图片</text></view></view>
        <text class="hint">最多上传 9 张图片，单张不超过 5 MB，第一张作为封面。</text>
      </view>
      <view class="field"><view class="field-head"><text>商品标题</text><text class="ai" @click="optimizeTitle">AI 优化标题</text></view><input v-model.trim="form.title" maxlength="100" placeholder="例如：九成新小米台灯，宿舍自用" /></view>
      <view class="field inline"><text>分类</text><picker :range="categories" range-key="name" @change="changeCategory"><view class="picker">{{ categoryName || '请选择分类' }} ›</view></picker></view>
      <view class="field inline"><text>成色</text><picker :range="conditions" @change="changeCondition"><view class="picker">{{ form.conditionLevel || '请选择成色' }} ›</view></picker></view>
      <view class="price-row"><view class="field"><text>原价</text><input v-model="form.originalPrice" type="digit" placeholder="选填" /></view><view class="field"><text>售价</text><input v-model="form.price" type="digit" placeholder="必填" /></view></view>
      <view class="field"><view class="field-head"><text>商品描述</text><text class="ai" @click="generateDescription">AI 生成描述</text></view><textarea v-model.trim="form.description" maxlength="1000" placeholder="说明使用情况、外观瑕疵、配件等信息" /></view>
      <view class="field"><text>建议面交地点</text><input v-model.trim="form.meetLocation" maxlength="100" placeholder="例如：东区图书馆门口" /></view>
    </view>
    <button class="submit" :loading="submitting" @click="submit">确认发布</button>
  </view>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { createProduct, generateProductDescription, listCategories, optimizeProductTitle, uploadProductImage } from '../../api/product'
import { getToken, isVerifiedUser } from '../../utils/auth'
const MAX_SIZE = 5 * 1024 * 1024
const conditions = ['全新', '九成新', '八成新', '七成新', '有明显使用痕迹']
const categories = ref([])
const submitting = ref(false)
const form = reactive({ images: [], title: '', categoryId: '', conditionLevel: '', originalPrice: '', price: '', description: '', meetLocation: '' })
const categoryName = computed(() => categories.value.find((item) => item.id === Number(form.categoryId))?.name || '')
const toast = (title) => uni.showToast({ title, icon: 'none' })
const changeCategory = ({ detail }) => { form.categoryId = categories.value[detail.value].id }
const changeCondition = ({ detail }) => { form.conditionLevel = conditions[detail.value] }
const removeImage = (index) => form.images.splice(index, 1)
const chooseImages = async () => {
  try {
    const result = await uni.chooseImage({ count: 9 - form.images.length, sizeType: ['compressed'] })
    const files = result.tempFiles.filter((file) => { if (file.size > MAX_SIZE) { toast('单张图片不能超过 5 MB'); return false } return true })
    if (!files.length) return
    uni.showLoading({ title: '上传图片中' })
    for (const file of files) { const uploaded = await uploadProductImage(file.path); form.images.push(uploaded.imageUrl) }
  } catch (error) { if (!error?.errMsg?.includes('cancel')) toast(error.message || '图片上传失败') } finally { uni.hideLoading() }
}
const optimizeTitle = async () => {
  if (!form.title) return toast('请先填写一个基础标题')
  try { uni.showLoading({ title: 'AI 正在优化' }); const { titles } = await optimizeProductTitle({ title: form.title, categoryName: categoryName.value, conditionLevel: form.conditionLevel }); uni.showActionSheet({ itemList: titles, success: ({ tapIndex }) => { form.title = titles[tapIndex] } }) }
  catch (error) { toast(error.message || 'AI 服务暂时不可用，你仍可手动填写') } finally { uni.hideLoading() }
}
const generateDescription = async () => {
  if (!form.title) return toast('请先填写标题')
  try { uni.showLoading({ title: 'AI 正在生成' }); const result = await generateProductDescription({ title: form.title, categoryName: categoryName.value, conditionLevel: form.conditionLevel, meetLocation: form.meetLocation }); form.description = result.description }
  catch (error) { toast(error.message || 'AI 服务暂时不可用，你仍可手动填写') } finally { uni.hideLoading() }
}
const validate = () => {
  if (!form.images.length) return '请至少上传一张商品图片'
  if (!form.title) return '请填写商品标题'
  if (!form.categoryId) return '请选择商品分类'
  if (!form.conditionLevel) return '请选择商品成色'
  if (!form.price || Number(form.price) <= 0) return '请填写正确的商品售价'
  if (!form.description) return '请填写商品描述'
  return ''
}
const submit = async () => {
  const message = validate(); if (message) return toast(message)
  submitting.value = true
  try { await createProduct({ ...form, price: Number(form.price), originalPrice: Number(form.originalPrice || 0) }); uni.setStorageSync('PRODUCT_LIST_DIRTY', true); uni.showToast({ title: '发布成功', icon: 'success' }); setTimeout(() => uni.switchTab({ url: '/pages/home/home' }), 600) }
  catch (error) { toast(error.message || '发布失败') } finally { submitting.value = false }
}
onLoad(async () => {
  if (!getToken()) return uni.navigateTo({ url: '/pages/login/login' })
  if (!isVerifiedUser()) return uni.navigateTo({ url: '/pages/student-auth/student-auth' })
  try { categories.value = (await listCategories()).filter((item) => Number(item.id) !== 0) } catch (error) { toast(error.message) }
})
</script>

<style scoped>
.page { min-height: 100vh; padding: 28rpx; box-sizing: border-box; }.intro { margin-bottom: 22rpx; }.title,.muted,.hint { display: block; }.title { font-size: 38rpx; font-weight: 700; }.muted,.hint { margin-top: 8rpx; color: #929c98; font-size: 22rpx; }
.card { padding: 4rpx 22rpx; border-radius: 18rpx; background: #fff; }.field { padding: 22rpx 0; border-bottom: 2rpx solid #f0f2f1; }.field-head,.inline { display: flex; align-items: center; justify-content: space-between; }.ai,.picker { color: #23734f; font-size: 25rpx; }
input { height: 66rpx; margin-top: 8rpx; }textarea { width: 100%; height: 190rpx; margin-top: 14rpx; padding: 14rpx; border-radius: 14rpx; background: #f6f8f5; box-sizing: border-box; }
.image-grid { display: grid; grid-template-columns: repeat(3,1fr); gap: 12rpx; margin-top: 16rpx; }.image-item,.upload { position: relative; height: 170rpx; overflow: hidden; border-radius: 14rpx; }.image-item image { width: 100%; height: 100%; }.remove { position: absolute; top: 6rpx; right: 6rpx; width: 34rpx; height: 34rpx; border-radius: 50%; background: rgba(0,0,0,.6); color: #fff; text-align: center; line-height: 32rpx; }.upload { display: flex; flex-direction: column; align-items: center; justify-content: center; border: 2rpx dashed #aac4b8; color: #6b8b7e; font-size: 42rpx; }.upload text { font-size: 22rpx; }
.price-row { display: grid; grid-template-columns: repeat(2,1fr); gap: 24rpx; }.submit { margin-top: 24rpx; border-radius: 999rpx; background: #23734f; color: #fff; }
</style>
