const TOKEN_KEY = 'CAU_USED_GOODS_TOKEN'
const USER_KEY = 'CAU_USED_GOODS_USER'

export const getToken = () => {
  return uni.getStorageSync(TOKEN_KEY) || ''
}

export const setToken = (token) => {
  uni.setStorageSync(TOKEN_KEY, token)
}

export const getUser = () => {
  return uni.getStorageSync(USER_KEY) || null
}

export const setUser = (user) => {
  uni.setStorageSync(USER_KEY, user)
}

export const clearAuth = () => {
  uni.removeStorageSync(TOKEN_KEY)
  uni.removeStorageSync(USER_KEY)
}

export const saveLoginResult = (result) => {
  if (result?.token) {
    setToken(result.token)
  }
  if (result?.user) {
    setUser(result.user)
  }
}

export const isVerifiedUser = (user = getUser()) => {
  return user?.authStatus === 'VERIFIED' && user?.accountStatus === 'NORMAL'
}
