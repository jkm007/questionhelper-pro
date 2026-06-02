<script setup lang="ts">
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'
import { wsClient } from '@/utils/websocket'
import { isLoggedIn } from '@/utils/auth'

onLaunch(() => {
  console.log('App Launch')
  // 已登录则连接 WebSocket
  if (isLoggedIn()) {
    wsClient.connect()
  }

  // 监听网络恢复，自动重连 WebSocket
  uni.onNetworkStatusChange((res) => {
    if (res.isConnected && isLoggedIn()) {
      wsClient.connect()
    }
  })
})

onShow(() => {
  console.log('App Show')
})

onHide(() => {
  console.log('App Hide')
})
</script>

<style>
page {
  background-color: #f5f5f5;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', 'Helvetica Neue', Helvetica, Arial, sans-serif;
}

.container {
  padding: 20rpx;
}

/* 全局按钮样式 */
.btn-primary {
  background-color: #4A90D9;
  color: #ffffff;
  border-radius: 16rpx;
  height: 88rpx;
  line-height: 88rpx;
  font-size: 32rpx;
}

.btn-primary:active {
  opacity: 0.8;
}

/* 全局卡片样式 */
.card {
  background-color: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);
}
</style>
