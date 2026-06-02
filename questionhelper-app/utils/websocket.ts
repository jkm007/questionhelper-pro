import { getToken } from '@/utils/auth'

class WebSocketClient {
  private socketTask: UniApp.SocketTask | null = null
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null
  private reconnectAttempts = 0
  private maxReconnectDelay = 30000
  private baseReconnectDelay = 1000
  private isManualClose = false
  private messageHandlers: Map<string, (data: any) => void> = new Map()

  connect() {
    const token = getToken()
    if (!token) return

    const wsUrl = `${this.getWsBaseUrl()}/ws?token=${token}`

    this.socketTask = uni.connectSocket({
      url: wsUrl,
      success: () => console.log('[WS] 连接中...')
    })

    this.socketTask.onOpen(() => {
      console.log('[WS] 连接成功')
      this.reconnectAttempts = 0
      this.startHeartbeat()
    })

    this.socketTask.onMessage((res) => {
      this.handleMessage(res.data)
    })

    this.socketTask.onClose(() => {
      console.log('[WS] 连接关闭')
      this.stopHeartbeat()
      if (!this.isManualClose) {
        this.scheduleReconnect()
      }
    })

    this.socketTask.onError((err) => {
      console.error('[WS] 连接错误', err)
    })
  }

  private getWsBaseUrl(): string {
    const apiUrl = import.meta.env.VITE_API_BASE_URL || '/api/v1'
    const wsProtocol = apiUrl.startsWith('https') ? 'wss' : 'ws'
    const baseUrl = apiUrl.replace(/^https?:\/\//, '').replace(/\/api\/v1$/, '')
    return `${wsProtocol}://${baseUrl}`
  }

  private startHeartbeat() {
    this.heartbeatTimer = setInterval(() => {
      this.send({ type: 'ping' })
    }, 30000)
  }

  private stopHeartbeat() {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer)
      this.heartbeatTimer = null
    }
  }

  private scheduleReconnect() {
    const delay = Math.min(
      this.baseReconnectDelay * Math.pow(2, this.reconnectAttempts),
      this.maxReconnectDelay
    )
    this.reconnectAttempts++

    console.log(`[WS] ${delay}ms 后重连 (第 ${this.reconnectAttempts} 次)`)

    this.reconnectTimer = setTimeout(() => {
      this.connect()
    }, delay)
  }

  private handleMessage(rawData: string | ArrayBuffer) {
    try {
      const data = JSON.parse(rawData as string)
      if (data.type === 'pong') return

      const handler = this.messageHandlers.get(data.type)
      if (handler) {
        handler(data)
      }

      this.onNotification(data)
    } catch (e) {
      console.error('[WS] 消息解析失败', e)
    }
  }

  private onNotification(data: any) {
    uni.$emit('notification', data)
    if (data.title) {
      uni.showToast({
        title: data.title,
        icon: 'none',
        duration: 3000
      })
    }
  }

  send(data: any) {
    this.socketTask?.send({ data: JSON.stringify(data) })
  }

  on(type: string, handler: (data: any) => void) {
    this.messageHandlers.set(type, handler)
  }

  off(type: string) {
    this.messageHandlers.delete(type)
  }

  close() {
    this.isManualClose = true
    this.stopHeartbeat()
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
    }
    this.socketTask?.close({})
  }
}

export const wsClient = new WebSocketClient()
