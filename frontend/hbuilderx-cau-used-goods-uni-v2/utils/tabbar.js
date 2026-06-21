const TAB_BAR_PAGES = [
  'pages/home/home',
  'pages/publish/publish',
  'pages/messages/messages',
  'pages/index/index'
]

export function isCurrentTabBarPage() {
  const pages = getCurrentPages()
  const current = pages[pages.length - 1]
  return TAB_BAR_PAGES.includes(current?.route)
}

export function updateMessageTabBadge(total) {
  if (!isCurrentTabBarPage()) return
  if (Number(total || 0) > 0) {
    uni.setTabBarBadge({ index: 2, text: total > 99 ? '99+' : String(total) })
  } else {
    uni.removeTabBarBadge({ index: 2 })
  }
}
