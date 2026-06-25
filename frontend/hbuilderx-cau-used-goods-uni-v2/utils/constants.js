export const ORDER_STATUS = {
  PENDING_CONFIRM: { label: '\u5f85\u5356\u5bb6\u786e\u8ba4', tone: 'warning' },
  WAIT_MEET: { label: '\u5f85\u7ebf\u4e0b\u4ea4\u6613', tone: 'primary' },
  COMPLETED: { label: '\u5df2\u5b8c\u6210', tone: 'success' },
  CANCELED: { label: '\u5df2\u53d6\u6d88', tone: 'muted' },
  CANCELLED: { label: '\u5df2\u53d6\u6d88', tone: 'muted' },
  EXCEPTION_CLOSED: { label: '\u5f02\u5e38\u5173\u95ed', tone: 'danger' }
}

export const REPORT_STATUS = {
  PENDING: { label: '\u5f85\u5904\u7406', tone: 'warning' },
  PROCESSING: { label: '\u5904\u7406\u4e2d', tone: 'primary' },
  APPROVED: { label: '\u5df2\u901a\u8fc7', tone: 'success' },
  RESOLVED: { label: '\u5df2\u5904\u7406', tone: 'success' },
  REJECTED: { label: '\u5df2\u9a73\u56de', tone: 'muted' },
  CLOSED: { label: '\u5df2\u5173\u95ed', tone: 'muted' }
}

export const APPEAL_STATUS = {
  PENDING: { label: '\u5f85\u5904\u7406', tone: 'warning' },
  PROCESSING: { label: '\u5904\u7406\u4e2d', tone: 'primary' },
  APPROVED: { label: '\u5df2\u901a\u8fc7', tone: 'success' },
  REJECTED: { label: '\u5df2\u9a73\u56de', tone: 'muted' },
  CLOSED: { label: '\u5df2\u64a4\u56de', tone: 'muted' }
}

export const REPORT_REASON = {
  FAKE_PRODUCT: '\u865a\u5047\u6216\u8fdd\u89c4\u5546\u54c1',
  INAPPROPRIATE_CONTENT: '\u4e0d\u5f53\u5185\u5bb9',
  SCAM: '\u6b3a\u8bc8\u98ce\u9669',
  TRADE_DISPUTE: '\u4ea4\u6613\u7ea0\u7eb7',
  OTHER: '\u5176\u4ed6\u95ee\u9898'
}

export const TARGET_TYPE = {
  PRODUCT: '\u5546\u54c1',
  ORDER: '\u4ea4\u6613\u8ba2\u5355',
  USER: '\u7528\u6237',
  REPORT: '\u4e3e\u62a5'
}

export const MESSAGE_TYPE = {
  ORDER: '\u8ba2\u5355\u6d88\u606f',
  REPORT: '\u4e3e\u62a5\u6d88\u606f',
  SYSTEM: '\u7cfb\u7edf\u6d88\u606f'
}
